package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"strings"

	kyaml "github.com/ghodss/yaml" // we intentionally use a different yaml package here for Kubernetes objects because gopkg.in/yaml.v2 is not meant to serialize k8s objects because of UnmarshalJSON/UnmarshalYAML and `json:""`/`yaml:""` dichotomy resulting in panic when used
	cmdutil "github.com/redhat-developer/opencompose/pkg/cmd/util"
	"github.com/redhat-developer/opencompose/pkg/encoding"
	"github.com/redhat-developer/opencompose/pkg/object"
	"github.com/redhat-developer/opencompose/pkg/transform"
	"github.com/redhat-developer/opencompose/pkg/transform/kubernetes"
	"github.com/redhat-developer/opencompose/pkg/transform/openshift"
	pkgutil "github.com/redhat-developer/opencompose/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/meta"
	"k8s.io/client-go/pkg/runtime"
	"os"
)

var (
	convertExample = `opencompose convert -f opencompose.yaml`
)

const (
	retryAttempts = 3
)

func NewCmdConvert(v *viper.Viper, out, outerr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "convert",
		Short:   "Convert OpenCompose files into Kubernetes (and OpenShift) artifacts",
		Long:    "Convert OpenCompose files into Kubernetes (and OpenShift) artifacts",
		Example: convertExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunConvert(v, cmd, out, outerr)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Parent().PersistentPreRunE != nil {
				if err := cmd.Parent().PersistentPreRunE(cmd, args); err != nil {
					return err
				}
			}

			// We have to bind Viper in Run because there is only one instance to avoid collisions between subcommands
			cmdutil.AddIOFlagsViper(v, cmd)

			return nil
		},
	}

	cmdutil.AddIOFlags(cmd)

	return cmd
}

func GetValidatedObject(v *viper.Viper, cmd *cobra.Command, out, outerr io.Writer) (*object.OpenCompose, error) {
	files := v.GetStringSlice(cmdutil.Flag_File_Key)
	if len(files) < 1 {
		return nil, cmdutil.UsageError(cmd, "there has to be at least one file")
	}
	var ocObjects []*object.OpenCompose

	for _, file := range files {
		var err error
		var in object.Input

		// Check if the passed resource points to STDIN or not
		if file == "-" {
			in.STDIN = true
			in.Data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("unable to read from stdin: %s", err)
			}
		} else if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
			// Check if the passed resource is a URL or not
			// TODO: add test for validating this against an actual resource on the web
			in.URL = &file
			in.Data, err = pkgutil.GetURLData(file, retryAttempts)
			if err != nil {
				return nil, fmt.Errorf("an error occurred while fetching data from the url %v: %v", file, err)
			}
		} else {
			in.FilePath = &file
			in.Data, err = ioutil.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("unable to read file '%s': %s", file, err)
			}
		}
		decoder, err := encoding.GetDecoderFor(in.Data)
		if err != nil {
			return nil, fmt.Errorf("could not find decoder for resource '%s': %s", file, err)
		}

		o, err := decoder.Decode(&in)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal data for file '%s': %s", file, err)
		}

		ocObjects = append(ocObjects, o)
	}

	// FIXME: implement merging OpenCompose obejcts
	openCompose := ocObjects[0]

	if err := openCompose.Validate(); err != nil {
		return nil, err
	}

	return openCompose, nil
}

func RunConvert(v *viper.Viper, cmd *cobra.Command, out, outerr io.Writer) error {
	o, err := GetValidatedObject(v, cmd, out, outerr)
	if err != nil {
		return err
	}

	var transformer transform.Transformer
	distro := v.GetString("distro")
	switch d := strings.ToLower(distro); d {
	case "kubernetes":
		transformer = &kubernetes.Transformer{}
	case "openshift":
		transformer = &openshift.Transformer{}
	default:
		return fmt.Errorf("unknown distro '%s'", distro)
	}

	runtimeObjects, err := transformer.Transform(o)
	if err != nil {
		return fmt.Errorf("transformation failed: %s", err)
	}

	var writeObject func(o runtime.Object, data []byte) error
	outputDir := v.GetString(cmdutil.Flag_OutputDir_Key)
	if outputDir == "" || outputDir == "-" {
		// don't use dir but write it to out (stdout)
		writeObject = func(o runtime.Object, data []byte) error {
			_, err := fmt.Fprintln(out, "---")
			if err != nil {
				return err
			}

			_, err = out.Write(data)
			return err
		}
	} else {
		// write files
		writeObject = func(o runtime.Object, data []byte) error {
			kind := o.GetObjectKind().GroupVersionKind().Kind
			m, ok := o.(meta.Object)
			if !ok {
				return fmt.Errorf("failed to cast runtime.object to meta.object (type is %s): %s", reflect.TypeOf(o).String(), err)
			}

			filename := fmt.Sprintf("%s-%s.yaml", m.GetName(), strings.ToLower(kind))

			err := ioutil.WriteFile(path.Join(outputDir, filename), data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %q: %s", filename, err)
			}

			fmt.Fprintf(out, "created file %q\n", filename)

			return nil
		}
	}

	for _, runtimeObject := range runtimeObjects {
		gvk, isUnversioned, err := api.Scheme.ObjectKind(runtimeObject)
		if err != nil {
			return fmt.Errorf("ConvertToVersion failed: %s", err)
		}
		if isUnversioned {
			return fmt.Errorf("ConvertToVersion failed: can't output unversioned type: %T", runtimeObject)
		}

		runtimeObject.GetObjectKind().SetGroupVersionKind(gvk)

		data, err := kyaml.Marshal(runtimeObject)
		if err != nil {
			return fmt.Errorf("failed to marshal object: %s", err)
		}

		err = writeObject(runtimeObject, data)
		if err != nil {
			return fmt.Errorf("failed to write object: %s", err)
		}
	}

	return nil
}

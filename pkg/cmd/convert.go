package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdutil "github.com/tnozicka/opencompose/pkg/cmd/util"
	"github.com/tnozicka/opencompose/pkg/encoding"
	"github.com/tnozicka/opencompose/pkg/object"
	"github.com/tnozicka/opencompose/pkg/transform"
	"github.com/tnozicka/opencompose/pkg/transform/kubernetes"
	"github.com/tnozicka/opencompose/pkg/transform/openshift"
	//"k8s.io/client-go/pkg/api"
	//"k8s.io/client-go/pkg/runtime/schema"
)

var (
	convertExample = `  # Converts file
  opencompose convert -f opencompose.yaml`
)

func NewCmdConvert(v *viper.Viper, out, outerr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "convert",
		Short:   "Converts OpenCompose files into Kubernetes (and OpenShift) artifacts",
		Long:    "Converts OpenCompose files into Kubernetes (and OpenShift) artifacts",
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
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("unable to read file '%s': %s", file, err)
		}

		decoder, err := encoding.GetDecoderFor(data)
		if err != nil {
			return nil, fmt.Errorf("could not find decoder for file '%s': %s", file, err)
		}

		o, err := decoder.Unmarshal(data)
		if err != nil {
			return nil, fmt.Errorf("could not unmarsha data for file '%s': %s", file, err)
		}

		ocObjects = append(ocObjects, o)
	}

	// FIXME: implement merging OpenCompose obejcts
	openCompose := ocObjects[0]

	openCompose.Validate()

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

	outputDir := v.GetString(cmdutil.Flag_OutputDir_Key)
	if outputDir == "-" {
		// don't use dir but write it to out (stdout)
		fmt.Fprintf(out, "runtimeObjects: %#v\n", runtimeObjects)
		for i, runtimeObject := range runtimeObjects {
			if i > 0 {
				fmt.Fprintln(out, "---")
			}

			// FIXME: (bellow)
			versionedObject := runtimeObject
			//versionedObject, err := api.Scheme.ConvertToVersion(runtimeObject, schema.GroupVersion{})
			//if err != nil {
			//	return fmt.Errorf("ConvertToVersion failed: %s", err)
			//}

			data, err := yaml.Marshal(versionedObject)
			if err != nil {
				return fmt.Errorf("failed to marshal object: %s", err)
			}
			fmt.Fprintln(out, string(data))
		}
	} else {
		// write files
	}

	return nil
}

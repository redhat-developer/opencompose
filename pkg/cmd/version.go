package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tnozicka/opencompose/pkg/version"
)

var (
	versionExample = `
		# Print version information
		kubectl version`
)

func NewCmdVersion(v *viper.Viper, out, outerr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print version information",
		Long:    "Print version information",
		Example: versionExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersion(v, cmd, out, outerr)
		},
	}
	return cmd
}

func RunVersion(v *viper.Viper, cmd *cobra.Command, out, outerr io.Writer) error {
	info := version.Get()
	fmt.Fprintln(out, info)
	return nil
}

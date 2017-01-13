package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdutil "github.com/tnozicka/opencompose/pkg/cmd/util"
)

var (
	validateExample = `
		# Print validate information
		opencompose validate`
)

func NewCmdValidate(v *viper.Viper, out, outerr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Print validate information",
		Long:    "Print validate information",
		Example: validateExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunValidate(v, cmd, out, outerr)
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

func RunValidate(v *viper.Viper, cmd *cobra.Command, out, outerr io.Writer) error {
	_, err := GetValidatedObject(v, cmd, out, outerr)
	return err
}

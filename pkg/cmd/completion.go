package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdutil "github.com/tnozicka/opencompose/pkg/cmd/util"
)

var (
	completionExample = `  source <(opencompose completion)`

	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		//"zsh":  runCompletionZsh,
	}
)

func NewCmdCompletion(v *viper.Viper, out, outerr io.Writer) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:     "completion",
		Short:   "Generate the opencompose completion code for bash",
		Long:    "Generate the opencompose completion code for bash",
		Example: completionExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCompletion(v, cmd, args, out, outerr)
		},
		ValidArgs: shells,
	}

	return cmd
}

func RunCompletion(v *viper.Viper, cmd *cobra.Command, args []string, out, outerr io.Writer) error {
	if len(args) == 0 {
		return cmdutil.UsageError(cmd, "Shell not specified.")
	}

	if len(args) > 1 {
		return cmdutil.UsageError(cmd, "Too many arguments. Expected only the shell type.")
	}

	runFunc, found := completionShells[args[0]]
	if !found {
		return cmdutil.UsageError(cmd, "Unsupported shell type %q.", args[0])
	}

	return runFunc(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, rootCmd *cobra.Command) error {
	return rootCmd.GenBashCompletion(out)
}

//func runCompletionZsh(out io.Writer, rootCmd *cobra.Command) error {
//	// TODO: when https://github.com/spf13/cobra/issues/107 is fixed
//}

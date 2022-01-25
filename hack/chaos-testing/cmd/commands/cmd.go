package commands

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
	"github.com/argoproj/argo-cd/v2/hack/gen-resources/util"
	"github.com/argoproj/argo-cd/v2/util/cli"
)

const (
	cliName = "argocd-chaos-testing"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cli.SetLogFormat(cmdutil.LogFormat)
	cli.SetLogLevel(cmdutil.LogLevel)
}

// NewCommand returns a new instance of an argocd command
func NewCommand() *cobra.Command {

	var generateOpts util.GenerateOpts

	var command = &cobra.Command{
		Use:   cliName,
		Short: "Chaos testing",
		Run: func(c *cobra.Command, args []string) {

		},
		DisableAutoGenTag: true,
	}

	command.PersistentFlags().StringVar(&generateOpts.Namespace, "kube-namespace", "argocd", "Name of the namespace on which Argo agent should be installed [$KUBE_NAMESPACE]")
	return command
}

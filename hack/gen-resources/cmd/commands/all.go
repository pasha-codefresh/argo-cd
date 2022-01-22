package commands

import (
	"log"
	"os"

	"context"

	generator "github.com/argoproj/argo-cd/v2/hack/gen-resources/generators"
	"github.com/argoproj/argo-cd/v2/hack/gen-resources/tools"

	"github.com/argoproj/argo-cd/v2/util/db"
	"github.com/argoproj/argo-cd/v2/util/settings"

	"github.com/argoproj/argo-cd/v2/hack/gen-resources/util"

	"github.com/spf13/cobra"
)

func NewAllResourcesCommand(opts *util.GenerateOpts) *cobra.Command {
	var command = &cobra.Command{
		Use:   "all",
		Short: "Manage all resources",
		Long:  "Manage all resources",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
			os.Exit(1)
		},
	}
	command.AddCommand(NewAllResourcesGenerationCommand(opts))
	command.AddCommand(NewAllResourcesCleanCommand(opts))
	return command
}

func NewAllResourcesGenerationCommand(opts *util.GenerateOpts) *cobra.Command {
	var command = &cobra.Command{
		Use:   "generate",
		Short: "Generate all resources",
		Long:  "Generate all resources",
		Run: func(c *cobra.Command, args []string) {
			argoClientSet := util.ConnectToK8sArgoClientSet()
			clientSet := tools.ConnectToK8sClientSet()

			settingsMgr := settings.NewSettingsManager(context.TODO(), clientSet, opts.Namespace)
			argoDB := db.NewDB(opts.Namespace, settingsMgr, clientSet)

			pg := generator.NewProjectGenerator(argoClientSet)
			ag := generator.NewApplicationGenerator(argoClientSet, clientSet, argoDB)
			rg := generator.NewRepoGenerator(tools.ConnectToK8sClientSet())

			err := pg.Generate(opts)
			if err != nil {
				log.Fatalf("Something went wrong, %v", err.Error())
			}
			err = ag.Generate(opts)
			if err != nil {
				log.Fatalf("Something went wrong, %v", err.Error())
			}
			err = rg.Generate(opts)
			if err != nil {
				log.Fatalf("Something went wrong, %v", err.Error())
			}
		},
	}
	return command
}

func NewAllResourcesCleanCommand(opts *util.GenerateOpts) *cobra.Command {
	var command = &cobra.Command{
		Use:   "clean",
		Short: "Clean all resources",
		Long:  "Clean all resources",
		Run: func(c *cobra.Command, args []string) {
			argoClientSet := util.ConnectToK8sArgoClientSet()
			clientSet := tools.ConnectToK8sClientSet()
			settingsMgr := settings.NewSettingsManager(context.TODO(), clientSet, opts.Namespace)
			argoDB := db.NewDB(opts.Namespace, settingsMgr, clientSet)
			pg := generator.NewProjectGenerator(argoClientSet)
			ag := generator.NewApplicationGenerator(argoClientSet, clientSet, argoDB)
			err := pg.Clean(opts)
			if err != nil {
				log.Fatalf("Something went wrong, %v", err.Error())
			}
			err = ag.Clean(opts)
			if err != nil {
				log.Fatalf("Something went wrong, %v", err.Error())
			}
		},
	}
	return command
}

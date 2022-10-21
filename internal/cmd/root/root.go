package root

import (
	"github.com/spf13/cobra"
	"kubernetes-controller/internal/manager"
)

var cfg manager.Config

func init() {
	rootCmd.Flags().AddFlagSet(cfg.FlagSet())
}

var rootCmd = &cobra.Command{
	PersistentPreRunE: bindEnvVars,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run(&cfg)
	},
	SilenceUsage: true,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

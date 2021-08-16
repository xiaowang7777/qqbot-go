package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qqbot-go/config"
	"qqbot-go/pkg/client"
)

var ROOT_CMD = &cobra.Command{
	Use:   "qqbot",
	Short: "qq bot",
	Long:  `qq bot commend`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		c := &config.Config{}
		c.Load()
		client.Run(c)
	},
}

func Execute() {
	if err := ROOT_CMD.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {

}

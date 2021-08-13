package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qqbot-go/config"
)

var CONFIG_CMD = cobra.Command{
	Use:   "config",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.New()
		if account, err := cmd.Flags().GetString("account"); err != nil {
			logrus.Fatal(err)
		} else {
			c.Account = account
		}

		if password, err := cmd.Flags().GetString("password"); err != nil {
			logrus.Fatal(err)
		} else {
			c.Password = password
		}

		c.Write()
	},
}

func init() {
	CONFIG_CMD.Flags().StringP("account", "a", "", "qq authentication account")
	CONFIG_CMD.Flags().StringP("password", "p", "", "qq authentication password,use QR if null")
}

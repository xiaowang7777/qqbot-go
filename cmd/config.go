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
		if account, err := cmd.Flags().GetInt64("account"); err != nil {
			logrus.Fatal(err)
		} else {
			c.Account.Uin = account
		}

		if password, err := cmd.Flags().GetString("password"); err != nil {
			logrus.Fatal(err)
		} else {
			c.Account.Password = password
		}

		c.Write()
	},
}

func init() {
	CONFIG_CMD.Flags().Int64P("account", "u", 0, "qq authentication account")
	CONFIG_CMD.Flags().StringP("password", "p", "", "qq authentication password,use QR if null")
}

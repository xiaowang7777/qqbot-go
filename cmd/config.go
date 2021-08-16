package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qqbot-go/config"
)

var configCmd = &cobra.Command{
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
	configCmd.Flags().Int64P("account", "u", 0, "QQ账号.")
	configCmd.Flags().StringP("password", "p", "", "QQ密码.")
	configCmd.Flags().BoolP("encrypt", "e", false, "是否需要开启加密.")
	configCmd.Flags().IntP("en-type", "", 1, `加密方式.
1.RSA.
2.DES.`)
}

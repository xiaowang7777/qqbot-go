package client

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
	"qqbot-go/utils"
)

const (
	passwordEncryptFileName = "password.encrypt"
)

func handlePasswordEncrypt() {
	conf := goal.GetConfig()
	if !conf.Account.Encrypt.Enable {
		return
	}

	if len(conf.Account.Password) == 0 && !utils.CheckFileExists(filepath.Join(config.GetFileDir(), passwordEncryptFileName)) {
		logrus.Error("已开启密码加密功能，但未找到加密后的文件，请在bot配置文件修改后重新启动bot.")
		logrus.Warnf("bot配置文件路径-> %s", config.GetConfigFilePath())
		logrus.Warn("bot即将关闭，请按Enter键继续.")
		utils.Readline()
		os.Exit(1)
	}

}

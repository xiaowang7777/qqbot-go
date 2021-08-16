package client

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
)

func Run(conf *config.Config) {
	goal.SetConfig(conf)
	_, err := login()
	if err != nil {
		logrus.Errorf("登陆失败，错误信息：%v", err)
	}
}

func newQQClient() *client.QQClient {
	c := client.NewClientEmpty()
	c.OnServerUpdated(func(qqClient *client.QQClient, event *client.ServerUpdatedEvent) bool {
		return true
	})
	c.OnLog(func(qqClient *client.QQClient, event *client.LogEvent) {
		switch event.Type {
		case "INFO":
			logrus.Info(event.Message)
		case "WARN":
			logrus.Warn(event.Message)
		case "ERROR":
			logrus.Error(event.Message)
		case "DEBUG":
			logrus.Debug(event.Message)
		}
	})
	return c
}

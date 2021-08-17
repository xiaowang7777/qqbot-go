package client

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"qqbot-go/config"
	"qqbot-go/pkg/controller"
	"qqbot-go/pkg/goal"
)

func Run(conf *config.Config) {
	goal.SetConfig(conf)
	q := &controller.QQ{}
	q.OnStarted(func(q *controller.QQ) {
		q.Cli.OnPrivateMessage(func(cli *client.QQClient, msg *message.PrivateMessage) {
			logrus.Infof("get private message from %d", msg.Sender.Uin)
			sendingMessage := message.NewSendingMessage()
			sendingMessage.Elements = append(sendingMessage.Elements, &message.TextElement{Content: "hello,word"})
			cli.SendPrivateMessage(msg.Sender.Uin, sendingMessage)
		})
	})

	err := q.Login()
	if err != nil {
		logrus.Errorf("登陆失败，错误信息：%v", err)
	}
	q.WatchTermSignal()
}

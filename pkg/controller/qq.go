package controller

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path/filepath"
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
	"qqbot-go/utils"
	"syscall"
)

type OnQuit func(q *QQ)

type BeforeStart func() bool

type OnStarted func(q *QQ)

type QQ struct {
	stopChannel chan os.Signal
	Cli         *client.QQClient
	beforeStart []BeforeStart
	onStarted   []OnStarted
	onQuit      []OnQuit
}

func (q *QQ) WatchTermSignal() {
	q.stopChannel = make(chan os.Signal)
	signal.Notify(q.stopChannel, syscall.SIGINT)
	signal.Notify(q.stopChannel, syscall.SIGTERM)
	signal.Notify(q.stopChannel, syscall.SIGKILL)
	<-q.stopChannel
	logrus.Warn("程序即将关闭...")
	q.onQuitExec()
	q.Cli.Disconnect()
	q.Cli.Release()
}

func (q *QQ) BeforeStart(start BeforeStart) {
	q.beforeStart = append(q.beforeStart, start)
}

func (q *QQ) OnStarted(started OnStarted) {
	q.onStarted = append(q.onStarted, started)
}

func (q *QQ) OnQuit(quit OnQuit) {
	q.onQuit = append(q.onQuit, quit)
}

func (q *QQ) onQuitExec() {
	for i := 0; i < len(q.onQuit); i++ {
		q.onQuit[i](q)
	}
}

func (q *QQ) beforeStartExec() {
	for i := 0; i < len(q.beforeStart); i++ {
		if !q.beforeStart[i]() {
			break
		}
	}
}

func (q *QQ) onStartedExec() {
	for i := 0; i < len(q.onStarted); i++ {
		q.onStarted[i](q)
	}
}

//Login qq登录
func (q *QQ) Login() error {
	q.beforeStartExec()
	q.Cli = newQQClient()
	passwordFile := filepath.Join(config.GetFileDir(), passwordEncryptFileName)
	//配置文件中未设置密码并且未找到密码文件时使用二维码登录
	if len(goal.GetConfig().Account.Password) == 0 && !utils.CheckFileExists(passwordFile) {
		if cliN, err := qrcodeLogin(q.Cli); err != nil {
			return err
		} else {
			q.Cli = cliN
		}
	} else {
		if cliN, err := commonLogin(q.Cli); err != nil {
			return err
		} else {
			q.Cli = cliN
		}
	}
	logrus.Info("qq login complete.")
	q.onStartedExec()
	return nil
}

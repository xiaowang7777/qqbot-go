package client

import (
	"bytes"
	"errors"
	qrTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/sirupsen/logrus"
	"github.com/tuotoo/qrcode"
	"os"
	"path/filepath"
	"qqbot-go/config"
	"qqbot-go/utils"
	"strings"
	"time"
)

var (
	qrTimeoutChannel = make(chan interface{})
)

func commonLogin(cli *client.QQClient) error {
	response, err := cli.Login()
	if err != nil {
		return err
	}
	return handleQQClientLogin(cli, response)
}

//qrcodeLogin 处理二维码登录的扫码操作
func qrcodeLogin(cli *client.QQClient) error {
	qrcodeFile := filepath.Join(config.GetFileDir(), "qrcode.png")
	response, err := cli.FetchQRCode()
	if err != nil {
		return err
	}
	f, err := qrcode.Decode(bytes.NewBuffer(response.ImageData))
	if err != nil {
		return err
	}
	_ = os.WriteFile(qrcodeFile, response.ImageData, 0o664)
	defer func() {
		_ = os.Remove(qrcodeFile)
	}()

	if cli.Uin != 0 {
		logrus.Infof("请使用账号：%v 登陆手机QQ扫描二维码", cli.Uin)
	} else {
		logrus.Infof("请登录手机QQ扫描登录")
	}

	qrTerminal.New().Get(f.Content).Print()

	return handleQRCodeWithTimeout(cli, response.Sig, 60)
}

//handleQRCodeWithTimeout 处理二维码登录扫码后的操作
func handleQRCodeWithTimeout(cli *client.QQClient, sig []byte, sec int32) error {
	startTimeout(sec)
	defer close(qrTimeoutChannel)

	reference := utils.NewAtomicReference(false)

	end := false
	var errCon error = nil

	for {
		select {
		case <-qrTimeoutChannel:
			return errors.New("qrcode login timeout")
		case <-time.After(time.Duration(1)):
			if reference.CompareAndSet(false, true) {
				func() {
					defer reference.Set(false)
					status, err := cli.QueryQRCodeStatus(sig)
					if err != nil {
						errCon = err
						end = true
						return
					}
					switch status.State {
					case client.QRCodeCanceled:
						logrus.Fatal("扫码被取消")
					case client.QRCodeTimeout:
						logrus.Fatal("二维码过期")
					case client.QRCodeWaitingForConfirm:
						logrus.Info("扫码成功，请在手机端确认")
					case client.QRCodeConfirmed:
						response, err := cli.QRCodeLogin(status.LoginInfo)
						if err != nil {
							errCon = err
						} else {
							errCon = handleQQClientLogin(cli, response)
						}
						end = true
						return
					default:
						//忽略其他状态
					}
				}()
				if end {
					return errCon
				}
			}
		}
	}
}

func handleQQClientLogin(cli *client.QQClient, resp *client.LoginResponse) error {
	for {
		switch resp.Error {
		case client.SliderNeededError:
			logrus.Warn("登录需要滑条验证码。")
			logrus.Warn("请参考文档 -> https://docs.go-cqhttp.org/faq/slider.html <- 进行处理")
			logrus.Warn("1. 自行抓包并获取 Ticket 输入.")
			logrus.Warn("2. 使用手机QQ扫描二维码登入. (推荐)")
			logrus.Warn("请输入(1 - 2) (将在10秒后自动选择2)：")
			t := utils.ReadlineWithTimeout(time.Second*10, "2")
			if strings.Contains(t, "1") {
				logrus.Infoln("--------------------")
				logrus.Infof("请用浏览器打开 -> %v <- 并获取Ticket.", resp.VerifyUrl)
				logrus.Infoln("--------------------")
				logrus.Infoln("请输入Ticket： (Enter 提交)")
				tick := utils.Readline()
				response, err := cli.SubmitTicket(tick)
				if err != nil {
					return err
				}
				return handleQQClientLogin(cli, response)
			}
			cli.Disconnect()
			cli.Release()
			qqClient := newQQClient()
			return qrcodeLogin(qqClient)
		case client.SMSNeededError:

		}
	}
}

func startTimeout(sec int32) {
	time.Sleep(time.Duration(sec))
	qrTimeoutChannel <- &config.Config{}
}

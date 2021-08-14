package client

import (
	"bytes"
	"errors"
	qrTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Mrs4s/MiraiGo/client"
	errorsPkg "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tuotoo/qrcode"
	"os"
	"path/filepath"
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
	newError "qqbot-go/pkg/nerror"
	"qqbot-go/utils"
	"strings"
	"time"
)

var (
	qrTimeoutChannel = make(chan interface{})
)

//login qq登录
func login() (*client.QQClient, error) {
	cli := newQQClient()
	if len(goal.GetConfig().Account.Password) != 0 {
		if cliN, err := qrcodeLogin(cli); err != nil {
			return nil, err
		} else {
			cli = cliN
		}
	} else {
		if cliN, err := commonLogin(cli); err != nil {
			return nil, err
		} else {
			cli = cliN
		}
	}
	logrus.Info("qq login complete.")
	return cli, nil
}

//用户名密码登录
func commonLogin(cli *client.QQClient) (*client.QQClient, error) {
	cli.Uin = goal.GetConfig().Account.Uin
	cli.PasswordMd5 = handlePassword(goal.GetConfig().Account.Password)
	response, err := cli.Login()
	if err != nil {
		return nil, err
	}
	return handleQQClientLogin(cli, response)
}

//qrcodeLogin 处理二维码登录的扫码操作
func qrcodeLogin(cli *client.QQClient) (*client.QQClient, error) {
	qrcodeFile := filepath.Join(config.GetFileDir(), "qrcode.png")
	response, err := cli.FetchQRCode()
	if err != nil {
		return nil, err
	}
	f, err := qrcode.Decode(bytes.NewBuffer(response.ImageData))
	if err != nil {
		return nil, err
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

//handleQRCodeWithTimeout 处理二维码登录扫码获取返回信息的操作
func handleQRCodeWithTimeout(cli *client.QQClient, sig []byte, sec int32) (*client.QQClient, error) {
	startTimeout(sec)
	defer close(qrTimeoutChannel)

	reference := utils.NewAtomicReference(false)
	end := false
	var errCon error = nil

	for {
		select {
		case <-qrTimeoutChannel:
			return nil, errors.New("qrcode login timeout")
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
							cli, errCon = handleQQClientLogin(cli, response)
						}
						end = true
						return
					default:
						//忽略其他状态
					}
				}()
				if end {
					return cli, errCon
				}
			}
		}
	}
}

//handleQQClientLogin 处理QQ登陆后的返回信息
func handleQQClientLogin(cli *client.QQClient, resp *client.LoginResponse) (*client.QQClient, error) {
	switch resp.Error {
	//滑动条
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
				return nil, err
			}
			return handleQQClientLogin(cli, response)
		}
		cli.Disconnect()
		cli.Release()
		cli = newQQClient()
		return qrcodeLogin(cli)
	//需要发送验证码
	case client.SMSNeededError:
		logrus.Warnf("需要发送手机验证码->%s，点击Enter发送.", resp.SMSPhone)
		utils.Readline()
		return sendSms(cli)
	case client.SMSOrVerifyNeededError:
		logrus.Info("账号已开启设备锁，请选择验证方式：")
		logrus.Infof("1.向手机：-> %s <-发送验证码.")
		logrus.Info("2.手机扫码登录.")
		logrus.Warn("请输入（1-2），10秒后自动选2")
		text := utils.ReadlineWithTimeout(time.Second*10, "2")
		if strings.Contains(text, "1") {
			return sendSms(cli)
		} else if strings.Contains(text, "2") {
			cli.Disconnect()
			cli.Release()
			cli = newQQClient()
			return qrcodeLogin(cli)
		}
		fallthrough
	case client.UnsafeDeviceError:
		logrus.Warnf("账号已开启设备锁，请前往-> %s <-验证后重启bot.", resp.VerifyUrl)
		logrus.Info("bot即将关闭，请按Enter键继续.")
		utils.Readline()
		os.Exit(1)
		return nil, nil
	case client.NeedCaptcha, client.OtherLoginError, client.TooManySMSRequestError, client.UnknownLoginError:
		logrus.Warn("发生不可恢复的错误！QQ返回错误信息：")
		logrus.Errorf("%s", resp.ErrorMessage)
		//os.Exit(-1)
		return nil, errorsPkg.WithStack(newError.TypeNotFoundError)
	}
	if resp.Success {
		logrus.Info("登陆成功！")
		return cli, nil
	} else {
		logrus.Errorf("未知异常！QQ返回错误消息-> %s <-", resp.ErrorMessage)
		return nil, errorsPkg.WithStack(newError.UnknownError)
	}
}

//sendSms 发送手机验证码
func sendSms(cli *client.QQClient) (*client.QQClient, error) {
	if !cli.RequestSMS() {
		logrus.Error("验证码发送失败！可能是发送过于频繁.")
		return nil, errorsPkg.WithStack(newError.SendSMSError)
	}
	text := utils.Readline()
	if response, err := cli.SubmitSMS(text); err != nil {
		return nil, err
	} else {
		return handleQQClientLogin(cli, response)
	}
}

func startTimeout(sec int32) {
	time.Sleep(time.Duration(sec))
	qrTimeoutChannel <- &config.Config{}
}

func handlePassword(pass string) [16]byte {

	return [16]byte{}
}

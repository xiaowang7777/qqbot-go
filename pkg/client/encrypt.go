package client

import (
	"bytes"
	"crypto/md5"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"qqbot-go/config"
	"qqbot-go/pkg/goal"
	"qqbot-go/utils"
	"strings"
)

const (
	passwordEncryptFileName = "password.encrypt"

	//DES加密文件
	desEncryptFileName = "des.encrypt.pem"
	//RSA私钥文件
	rsaPrivateFileName = "rsa.private.pem"
	//RSA公钥文件
	rsaPublicFileName = "rsa.public.pem"
)

var (
	nilByte [16]byte
)

//handlePasswordEncrypt 密码加密处理，加密后写入文件供后续使用
func handlePasswordEncrypt() {
	conf := goal.GetConfig()
	password := conf.Account.Password
	passwordFile := filepath.Join(config.GetFileDir(), passwordEncryptFileName)
	//判断密码文件是否已存在，已存在则不做处理
	if utils.CheckFileExists(passwordFile) || !conf.Account.Encrypt.Enable {
		return
	}

	//密码文件不存在时，创建它
	switch conf.Account.Encrypt.Type {
	case config.RSA:
		handleRSAEncrypt(password, passwordFile)
	case config.DES:

	}
}

//handleRSAEncrypt 当选在RSA作为加密方式时，未找到密码信息的加密处理
func handleRSAEncrypt(password, passwordFile string) {
	rsaPublicFilePath := filepath.Join(config.GetFileDir(), rsaPublicFileName)
	rsaPrivateFilePath := filepath.Join(config.GetFileDir(), rsaPrivateFileName)
	logrus.Warn("即将使用RSA算法对密码进行加密.")
	logrus.Warn("注意：程序只有在每次启动bot，并且需要进行密码解密时才会需要私钥信息，且程序不会保存用户的私钥文件，请用户妥善保存！")
	handlePasswordWriteToFile := func(rsaPublicFilePath, password, passwordFile string) {
		rsaPublicFile, err := os.Open(rsaPublicFilePath)
		if err != nil {
			logrus.Error("打开公钥文件失败！程序即将退出，错误信息：%v", err)
			os.Exit(1)
		}
		publicFileStat, err := rsaPublicFile.Stat()

		sum := md5.New().Sum(utils.StringToByte(password))

		if encryptRes, err := goal.RSAEncrypt(rsaPublicFile, publicFileStat.Size(), sum); err != nil {
			logrus.Errorf("进行公钥加密时失败！程序即将退出，错误信息：%v", err)
			os.Exit(1)
		} else {
			if err := os.WriteFile(passwordFile, encryptRes, 0644); err != nil {
				logrus.Errorf("将加密后的密码写入文件时失败！程序即将退出，错误信息:%v", err)
				os.Exit(1)
			}
			logrus.Warnf("将加密密码写入文件成功！文件路劲-> %s", passwordFile)
		}
	}
	//不存在公钥文件时
	if !utils.CheckFileExists(rsaPublicFilePath) {
		logrus.Warnf("未在文件路径：—> %s <- 找到RSA公钥文件！", rsaPublicFilePath)
		logrus.Warn("RSA公钥配置选择（5秒后自动选3）：")
		logrus.Warn("1.已有RSA公钥文件，直接输入公钥信息.")
		logrus.Warn("2.已有RSA公钥文件，输入公钥文件路径.")
		logrus.Warn("3.重新生成RSA加、解密密钥，并保存到文件.")
		text := utils.Readline()
		if strings.Contains("3", text) {
			logrus.Warn("RSA加、解密密钥生成中，请稍后...")
			if err := goal.GenerateRSAKey(config.GetFileDir(), rsaPrivateFileName, rsaPublicFileName); err != nil {
				logrus.Error("RSA加、解密密钥生成失败，bot即将退出.")
				logrus.Errorf("失败信息：%v", err)
				os.Exit(1)
			}
			logrus.Warn("RSA加、解密密钥生成成功，请妥善保存该信息.")
			logrus.Warnf("RSA公钥文件路径-> %s", rsaPublicFilePath)
			logrus.Warnf("RSA私钥文件路径-> %s", rsaPrivateFilePath)
			logrus.Warnf("使用生成公私钥文件进行加、解密！")
			handlePasswordWriteToFile(rsaPublicFilePath, password, passwordFile)
		} else if strings.Contains("2", text) {
			logrus.Warn("请输入RSA公钥文件路径，并按Enter键继续：")
			text := utils.Readline()
			handlePasswordWriteToFile(text, password, passwordFile)
		} else {
			logrus.Warn("请输入RSA公钥信息，并按Enter键继续：")
			text := utils.Readline()
			sum := md5.New().Sum(utils.StringToByte(password))
			resPublicContext := utils.StringToByte(text)
			if encryptRes, err := goal.RSAEncrypt(bytes.NewReader(resPublicContext), int64(len(resPublicContext)), sum); err != nil {
				logrus.Errorf("进行公钥加密时失败！程序即将退出，错误信息：%v", err)
				os.Exit(1)
			} else {
				if err := os.WriteFile(passwordFile, encryptRes, 0644); err != nil {
					logrus.Errorf("将加密后的密码写入文件时失败！程序即将退出，错误信息:%v", err)
					os.Exit(1)
				}
				logrus.Warnf("将加密密码写入文件成功！文件路劲-> %s", passwordFile)
			}
		}
	} else {
		//存在公钥公钥文件时
		logrus.Warnf("将使用位于-> %s <-的RSA公钥文件进行加密.", rsaPublicFilePath)
		handlePasswordWriteToFile(rsaPublicFilePath, password, passwordFile)
	}
}

////getPassword 获取明文的登录密码
//func getPassword(conf *config.Config) string {
//	password := conf.Account.Password
//	if len(password) > 0 {
//		return password
//	}
//	if conf.Account.Encrypt.Enable {
//		logrus.Warn("已开启密码加密功能，但未找到密码信息.")
//	} else {
//		logrus.Warn("未开启密码加密功能，但未找到密码信息.")
//	}
//
//	logrus.Warn("请选择以下方法输入密码（5秒后自动选2）：")
//	logrus.Warn("1.直接输入.")
//	logrus.Warn("2.在配置文件中写入密码后重启bot.")
//	text := utils.ReadlineWithTimeout(time.Second*5, "2")
//	if strings.Contains("1", text) {
//		logrus.Warn("请输入密码，并按Enter键继续（10秒后退出）：")
//
//		text := utils.ReadlineWithTimeout(time.Second*10, "")
//		if len(text) <= 0 {
//			logrus.Error("未读取到密码输入！程序即将退出.")
//			os.Exit(1)
//		}
//		return text
//	}
//	logrus.Warnf("bot即将关闭，请在配置文件添加密码信息后启动bot！配置文件路劲-> %s", config.GetConfigFilePath())
//	os.Exit(1)
//	return password
//}

func handlePasswordDecrypt() [16]byte {
	conf := goal.GetConfig()

	//不需要加、解密时
	if !conf.Account.Encrypt.Enable {
		return md5.Sum(utils.StringToByte(conf.Account.Password))
	}

	//需要加、解密
	passwordFilePath := filepath.Join(config.GetFileDir(), passwordEncryptFileName)
	passwordFile, err := os.Open(passwordFilePath)
	if err != nil {
		logrus.Errorf("打开密码文件失败！程序即将关闭，错误信息：%v", err)
		os.Exit(1)
	}

	defer passwordFile.Close()

	passwordFileStat, err := passwordFile.Stat()
	if err != nil {
		logrus.Errorf("出现了未知错误！程序即将关闭，错误信息：%v", err)
		os.Exit(1)
	}

	passwordEncrypt := make([]byte, passwordFileStat.Size())

	if _, err = passwordFile.Read(passwordEncrypt); err != nil {
		logrus.Errorf("读取文件失败！程序即将关闭，错误信息：%v", err)
		os.Exit(1)
	}

	switch conf.Account.Encrypt.Type {
	case config.RSA:
		return handleRSADecode(passwordEncrypt)
	case config.DES:

	}
	return nilByte
}

func handleRSADecode(passwordEncrypt []byte) [16]byte {
	rsaPrivateFilePath := filepath.Join(config.GetFileDir(), rsaPrivateFileName)

	f := func(rsaPrivateFilePath string, passwordEncrypt []byte) ([16]byte, error) {
		logrus.Warnf("尝试使用私钥-> %s <-进行解密", rsaPrivateFilePath)
		rsaPrivateFile, err := os.Open(rsaPrivateFilePath)
		if err != nil {
			logrus.Errorf("打开私钥文件-> %s <-失败！程序即将退出，错误信息：%v", rsaPrivateFilePath, err)
			os.Exit(1)
		}
		rsaPrivateFileStat, err := rsaPrivateFile.Stat()
		if rsaDecrypt, err := goal.RSADecrypt(rsaPrivateFile, rsaPrivateFileStat.Size(), passwordEncrypt); err != nil {
			return nilByte, err
		} else {
			return md5.Sum(rsaDecrypt), nil
		}
	}

	logrus.Warn("即将使用RSA算法对密码进行解密.")
	if utils.CheckFileExists(rsaPrivateFilePath) {
		logrus.Infof("找到默认私钥文件-> %s")
		if decodeRes, err := f(rsaPrivateFilePath, passwordEncrypt); err == nil {
			return decodeRes
		}
	}
	logrus.Warn("选择密钥输入格式（5秒后自动选2）：")
	logrus.Warn("1.直接输入密钥字符串.")
	logrus.Warn("2.指定密钥文件路径.")
	text := utils.Readline()
	if strings.Contains("2", text) {
		logrus.Warn("请输入私钥文件路径，并Enter键以继续.")
		text := utils.Readline()
		if decodeRes, err := f(text, passwordEncrypt); err == nil {
			return decodeRes
		} else {
			logrus.Errorf("使用私钥解码失败，请确认私钥是否合法，并且密码文件是由对应公钥加密，错误信息：%v", err)
			os.Exit(1)
		}
	}
	logrus.Warn("请输入私钥字符串，并按Enter键继续：")
	text = utils.Readline()
	reader := strings.NewReader(text)
	if rsaDecrypt, err := goal.RSADecrypt(reader, int64(reader.Len()), passwordEncrypt); err != nil {
		logrus.Errorf("使用私钥解码失败，请确认私钥是否合法，并且密码文件是由对应公钥加密，错误信息：%v", err)
		os.Exit(1)
	} else {
		return md5.Sum(rsaDecrypt)
	}
	return nilByte
}

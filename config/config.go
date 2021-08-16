package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"qqbot-go/utils"
	"runtime"
)

const (
	configFileName = ".qqbot.yaml"
	configPathName = "/qq_bot/"

	RSA EncryptType = 1
	DES EncryptType = 2
)

type EncryptType uint8
type GroupUin []int64
type FriendUin []int64

type Config struct {
	Account struct {
		Uin      int64  `yaml:"uin"`
		Password string `yaml:"password"`
		Encrypt  struct {
			Enable bool        `yaml:"enable"`
			Type   EncryptType `yaml:"type"`
			//DecodeFilePath string      `yaml:"decode_file_path"`
			//EncodeFilePath string      `yaml:"encode_file_path"`
		} `yaml:"encrypt"`
		ReLogin struct {
			Enable   bool `yaml:"enable"`
			Delay    uint `yaml:"delay"`
			MaxTimes uint `yaml:"max_times"`
		} `yaml:"re_login"`
	} `yaml:"account"`
}

//New 新建Config结构体，并从配置文件中读取信息
func New() *Config {
	c := &Config{}

	c.Load()

	return c
}

//Load 从配置文件中加载信息
func (c *Config) Load() {
	createFileIfNotExists()

	filePath := GetConfigFilePath()

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal(bs)
	}

	if err := yaml.Unmarshal(bs, c); err != nil {
		logrus.Fatal(err)
	}

}

//Write 将Config信息写入到配置文件中
func (c *Config) Write() {
	//配置文件路径
	filePath := GetConfigFilePath()

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	//将配置写入文件
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(c); err != nil {
		logrus.Fatal(err)
	}
}

func GetFileDir() string {
	return filepath.Join(getConfigDirPath(), configPathName)
}

func GetConfigFilePath() string {
	return filepath.Join(GetFileDir(), configFileName)
}

//createFileIfNotExists 当配置文件不存在时，创建它
func createFileIfNotExists() {
	if !utils.CheckFileExists(GetFileDir()) {
		if err := os.MkdirAll(GetFileDir(), 0666); err != nil {
			panic(err)
		}
	}
	filePath := filepath.Join(GetFileDir(), configFileName)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {

			file, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()
		} else {
			panic(err)
		}
	}
}

//getConfigDirPath 获取配置文件所在文件夹
func getConfigDirPath() string {
	if home := os.Getenv("QQ_BOT_HOME"); home != "" {
		return home
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

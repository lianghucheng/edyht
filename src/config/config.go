package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/szxby/tools/log"
)

// JSONConfig config in json
type JSONConfig struct {
	RWC        chan bool
	IPList     []string `json:"IpList"`
	Port       string   `json:"Port"`
	RedisAddr  string   `json:"RedisAddr"`  // redis地址
	RedisPass  string   `json:"RedisPass"`  // redis密码
	RedisDB    int      `json:"RedisDB"`    // redis库
	MongoAddr  string   `json:"MongoAddr"`  // mongo地址
	GameMongo  string   `json:"GameMongo"`  // mongo地址
	GameServer string   `json:"GameServer"` // 游戏服地址
	LocalIP    string   `json:"LocalIP"`    // 本地ip
	PassURL    []string `json:"PassURL"`    // 跳过验证
	ExportURL  []string `json:"ExportURL"`  // 批量打款操作
	BackDB     string   `json:"BackDB"`     // 后台数据库
	GameDB     string   `json:"GameDB"`     // 游戏数据库
}

// Enviroment develop----0 release----1
const (
	Develop = iota
	Release
)

var serverConfig = JSONConfig{}

func init() {
	log.Debug("init config")
	file := "config.json"
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("init config fail:%v", err)
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("read file fail %v", err)
	}
	envDta := map[string]interface{}{}
	err = json.Unmarshal(content, &envDta)
	if err != nil {
		log.Fatal("unmarshal config fail %v", err)
	}
	env, ok := envDta["env"].(float64)
	// log.Debug("check,%v", reflect.TypeOf(data["env"]))
	if !ok {
		log.Fatal("unmarshal config fail %v", envDta)
	}
	data := map[string]JSONConfig{}
	json.Unmarshal(content, &data)
	// log.Debug("check,%v", reflect.TypeOf(data["dev"]))
	if env == Develop {
		serverConfig = data["dev"]
	} else if env == Release {
		serverConfig = data["release"]
	}
	serverConfig.RWC = make(chan bool, 1)
	log.Debug("get config:%+v", serverConfig)
}

// GetConfig return config for server
func GetConfig() JSONConfig {
	return serverConfig
}

// SetIP from db
func SetIP(list []string) error {
	serverConfig.RWC <- true
	defer func() {
		<-serverConfig.RWC
	}()
	log.Debug("SetIP:%v", list)
	if len(list) == 0 {
		return errors.New("invalid list")
	}
	serverConfig.IPList = list
	return nil
}


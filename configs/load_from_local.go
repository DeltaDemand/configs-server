package appConfigs

import (
	"encoding/json"
	"github.com/DeltaDemand/configs-server/internal/client"
	"github.com/DeltaDemand/configs-server/server"
	"io/ioutil"
	"log"
)

type Config struct {
	ConfigServer server.Server `json:"config_server"`
	EtcdConfig   client.Etcd   `json:"Etcd"`
}

var config = Config{}
var configFile = "configs/config.json"

func LoadingConfigs() Config {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("配置文件读取失败\n", err)
		//读取失败，返回初始化零值的Config
		return config
	}
	configErr := json.Unmarshal(configBytes, &config)
	if configErr != nil {
		if err != nil {
			//无解析失败
			log.Println("解析json配置文件读取失败")
		}
	}
	return config
}

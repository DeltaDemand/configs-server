package inputArgs

import (
	"flag"
	appConfigs "github.com/DeltaDemand/configs-server/configs"

	"strings"
)

// Parse :configServer命令行参数解析
func Parse(confs *appConfigs.Config) {
	flag.StringVar(&confs.ConfigServer.ServerName, "name", confs.ConfigServer.ServerName, "本服务的名称")
	flag.StringVar(&confs.ConfigServer.ListenPort, "p", confs.ConfigServer.ListenPort, "开放的端口")

	var ends string
	flag.StringVar(&ends, "endPoints", "112.74.60.132:2379", "etcd节点的地址可以多个，用逗号(,)隔开")

	flag.Parse()
	//初始化全局变量值
	confs.EtcdConfig.EndPoints = strings.Split(ends, ",")
}

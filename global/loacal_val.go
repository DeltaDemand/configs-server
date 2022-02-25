package global

import (
	"github.com/DeltaDemand/configs-server/internal/model"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Groups map[string]model.Group //该配置有多少组要管理
	Split  = "/"
	Cli    *clientv3.Client //etcd客户端
)

// swagger:operation GET
// ---
// @summary: 忽略具体数据获取Groups的结构
// @Produce  json
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /struct [get]
func DisplayGroups(c *gin.Context) {
	for groupName, group := range Groups {
		c.String(200, "<%s>组成员:\n", groupName)
		for agentName, _ := range group {
			c.String(200, "<%s>\t", agentName)
		}
		c.String(200, "\n")
	}
}

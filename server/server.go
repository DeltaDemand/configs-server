package server

import (
	"context"
	"github.com/DeltaDemand/configs-server/global"
	"github.com/DeltaDemand/configs-server/internal/model"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

var opTimeout = time.Second * 5

type Server struct {
	ServerName string `json:"server_name"`
	ListenPort string `json:"listen_port"`
}

// swagger:operation GET
// ---
// @summary: 获取全部数据
// @Produce  json
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router / [get]
func (e *Server) EtcdGetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	response, err := kvc.Get(ctx, global.Split+e.ServerName, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}
	if len(response.Kvs) == 0 {
		log.Printf("etcd not found")
	}
	global.Groups = make(map[string]model.Group, 3)
	//把所有的配置解析，存入global.Groups中
	for _, kv := range response.Kvs {
		parseKV(kv.Key, kv.Value)
	}
	c.JSON(200, global.Groups)
}

// swagger:operation GET
// ---
// @summary: 获取该组内容
// @description: 获取该组所有配置内容
// @Produce  json
// @Param group_name query string true "组名"
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /group [get]
func (e *Server) EtcdGetGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	groupName := c.Query("group_name")

	response, err := kvc.Get(ctx, global.Split+e.ServerName+global.Split+groupName, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}
	if len(response.Kvs) == 0 {
		log.Printf("etcd not found")
	}
	for _, kv := range response.Kvs {
		parseKV(kv.Key, kv.Value)
	}
	c.JSON(200, global.Groups[groupName])
}

// swagger:operation PUT
// ---
// @summary: 修改该group内所有Agent某项配置
// @description: 获取该组所有配置内容
// @Produce  json
// @Param group_name query string true "组名"
// @Param config_name query string true "配置名"
// @Param value query string true "配置值"
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /group [put]
func (e *Server) EtcdPutGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	groupName := c.Query("group_name")
	configName := c.Query("config_name")
	value := c.Query("value")
	//获取该组
	group := global.Groups[groupName]
	for agentName, agent := range group {
		//对组内每个Agent的configName配置项更新
		_, err := kvc.Put(ctx, global.Split+e.ServerName+global.Split+groupName+global.Split+agentName+global.Split+configName, value)
		if err == nil {
			//更新etcd成功，也更新global.Groups的数据
			agent[configName] = value
			c.String(200, "%s修改成功，%s改为:\n%v", agentName, configName, agent[configName])
		}
	}
}

// swagger:operation GET
// ---
// summary: 获取Agent配置
// @Produce  json
// @Param group_name query string true "组名"
// @Param agent_name query string true "agent名"
// @Param config_name query string false "配置名"
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /agent [get]
func (e *Server) EtcdGetAgent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	groupName := c.Query("group_name")
	agentName := c.Query("agent_name")
	configName := c.Query("config_name")
	key := global.Split + e.ServerName + global.Split + groupName + global.Split + agentName
	response, err := kvc.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}
	if len(response.Kvs) == 0 {
		log.Printf("etcd key not found")
	}
	for _, kv := range response.Kvs {
		parseKV(kv.Key, kv.Value)
	}
	//用户不输入具体配置项，返回全部配置项
	if configName != "" {
		c.String(200, "%s组的%s[Agent]配置项%s为：\n%v", groupName, agentName, configName, global.Groups[groupName][agentName][configName])
	} else {
		c.JSON(200, global.Groups[groupName][agentName])
	}
}

// swagger:operation PUT
// ---
// summary: 修改Agent某项配置
// summary: 获取Agent配置
// @Produce  json
// @Param group_name query string true "组名"
// @Param agent_name query string true "agent名"
// @Param config_name query string true "配置名"
// @Param value query string true "配置值"
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /agent [put]
func (e *Server) EtcdPutAgent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	groupName := c.Query("group_name")
	agentName := c.Query("agent_name")
	configName := c.Query("config_name")
	value := c.Query("value")

	//对gent的configName配置项更新
	_, err := kvc.Put(ctx, global.Split+e.ServerName+global.Split+groupName+global.Split+agentName+global.Split+configName, value)
	if err == nil {
		//更新成功 agent[configName]=value
		global.Groups[groupName][agentName][configName] = value
		c.String(200, "%s修改成功，%s改为:\n%v", agentName, configName, global.Groups[groupName][agentName][configName])
	}
}

// swagger:operation Del
// ---
// summary: 删除Agent配置
// summary: agent名为为空则删除整组
// @Produce  json
// @Param group_name query string true "组名"
// @Param agent_name query string false "agent名"
// @Success 200 {string} string "请求成功"
// @Failure 400 {string} string "请求错误"
// @Router /agent [delete]
func (e *Server) EtcdDel(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	kvc := clientv3.NewKV(global.Cli)
	groupName := c.Query("group_name")
	agentName := c.Query("agent_name")

	//用户不输入具体配置项，返回全部配置项
	if agentName == "" {
		if _, ok := global.Groups[groupName]; ok {
			_, err := kvc.Delete(ctx, global.Split+e.ServerName+global.Split+groupName+global.Split, clientv3.WithPrefix())
			if err == nil {
				//删除成功，也删除groups内的元素
				delete(global.Groups, groupName)
				c.String(200, "<%s>组删除成功\n", groupName)
			}
		} else {
			c.String(200, "<%s>组不存在本系统，删除失败\n", groupName)
		}
	} else {
		if _, ok := global.Groups[groupName][agentName]; ok {
			_, err := kvc.Delete(ctx, global.Split+e.ServerName+global.Split+groupName+global.Split+agentName+global.Split, clientv3.WithPrefix())
			if err == nil {
				//删除成功，也删除groups内的元素
				delete(global.Groups[groupName], agentName)
				c.String(200, "<%s>Agent配置删除成功\n", agentName)
			}
		} else {
			c.String(200, "<%s>Agent不存在本系统，删除失败\n", groupName)
		}
	}
}

//根据kv存入数据global.Groups中
func parseKV(k []byte, v []byte) {
	keySli := strings.Split(string(k), global.Split)
	//大于三层的才要
	if len(keySli) > 3 {
		//keySli[0]是服务名，忽略
		groupName := keySli[2]
		agentName := keySli[3]
		configName := keySli[4]
		//该组不存在，创建
		if _, ok := global.Groups[groupName]; !ok {
			newGroup := make(map[string]model.ConfigItems, 3)
			global.Groups[groupName] = newGroup
		}
		//该agent配置不存在，创建
		if _, ok := global.Groups[groupName][agentName]; !ok {
			newConfigs := make(map[string]string, 7)
			global.Groups[groupName][agentName] = newConfigs
		}
		//填入值
		global.Groups[groupName][agentName][configName] = string(v)
	}
}

package client

import (
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type Etcd struct {
	EndPoints   []string `json:"endPoints"`
	DialTimeout int      `json:"dialTimeout"`
}

var (
	cli *clientv3.Client
)

func (e *Etcd) GetEtcdClient() *clientv3.Client {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   e.EndPoints,
		DialTimeout: time.Duration(e.DialTimeout) * time.Second,
	})
	if err != nil {
		//如果连接Etcd失败，直接退出
		log.Fatal("connect to etcd failed, err:%v\n", err)
		return nil
	}
	return cli
}

func (e *Etcd) CloseConn() {
	cli.Close()
}

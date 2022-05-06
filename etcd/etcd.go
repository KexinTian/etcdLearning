package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	Cli *clientv3.Client
)

//日志信息
type LogEntry struct {
	Path  string `json:"path"`  //日志存放的路径
	Topic string `json:"topic"` //日志发往kafka的to'pi'c
}

//初始化etcd的方法
func Init(addr string, timeout time.Duration) (err error) {
	Cli, err = clientv3.New(clientv3.Config{ //注意，这里不要用 := ，一旦用了这个，cli就是新建的局部变量了
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("Connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("Connect to etcd success")
	// defer cli.Close()
	return
}

//从etcd中根据 key 获取配置项信息
func GetConf(key string) (LogEntrys []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := Cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		// fmt.Printf("配置JSON:%v\n", ev.Value)//输出的是二进制数组
		fmt.Printf("配置JSON:%s\n", ev.Value)
		// err = json.Unmarshal(ev.Value, LogEntrys)
		err = json.Unmarshal(ev.Value, &LogEntrys)
		if err != nil {
			fmt.Printf("Unmarshal etcd value failed, err:%v\n", err)
			return
		}
	}
	return
}

package main

import (
	"context"
	"etcdLearning/conf"
	"etcdLearning/etcd"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-ini/ini"
	"go.etcd.io/etcd/clientv3"
)

// put  get del
//mian——1

func main_1() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success") //就算本地的etcd服务没有启动，这里也是连接成功的，会在下一步报错
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) //接受一个 Context 和超时时间作为参数，返回其子Context和取消函数cancel
	_, err = cli.Put(ctx, "hh", "helloworld")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "hh")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
	// del
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	delResp, err := cli.Delete(ctx, "hh")
	cancel()
	if err != nil {
		fmt.Printf("delete from etcd failed, err:%v\n", err)
		return
	}
	fmt.Printf("成功删除了 %d 条数据\n", delResp.Deleted)

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err = cli.Get(ctx, "hh")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	fmt.Printf("键值 %s 的value数量为：%d\n", "hh", resp.Count)
	// test
}

//watch

func main_2() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch key:q1mi change
	/*
		rch := cli.Watch(context.Background(), "hh")
		也可以写成:

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		rch := cli.Watch(ctx, "hh")
	*/
	rch := cli.Watch(context.Background(), "hh")
	//rch 是一个通道
	// WatchChan <-chan WatchResponse，表示这是一个接受WatchResponse（数据类型）的通道
	//chan表示双向通道，<-chan 表示接受通道，chan<- 表示发送通道
	//详细请搜索，golang chan
	fmt.Printf("打印：%#v\n", rch) // 打印：(clientv3.WatchChan)(0xc000095ce0)
	for wresp := range rch {    //rch是一个监听hh key的通道，每次hh 的value改变，这里都会监听到进入到循环体内
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value) //分别输出事件类型、key值、value值
		}
	}
}

//加载ini配置信息
func main_3() {
	var cfg = new(conf.AppConf)
	//加载配置文件到结构体
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed by map, error:%v\n", err)
		return
	}
	fmt.Printf("load ini by map success\n")
	fmt.Printf("cgf struct: %#v\n", cfg)

	//加载配置文件, 读取命名分区数据
	cfgs, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed, error:%v\n", err)
		return
	}
	fmt.Printf("load ini success\n")
	val := cfgs.Section("kafka").Key("address").Value()
	fmt.Printf("命名分区数据[kafka]address: %#v\n", val)

	//加载配置文件, 读取默认分区数据
	val = cfgs.Section("").Key("name").Value() //Value()则直接不做换回直接返回，只有一个返回值
	fmt.Printf("默认分区数据name: %#v\n", val)
	res, err := cfgs.Section("").Key("name").Int() //键值转为Int
	fmt.Printf("默认分区数据name: %#v\n", res)

}

//用配置文件初始化etcd ,并收集日志
func main_4() {
	var cfg = new(conf.AppConf)
	//1、 加载配置文件到结构体，初始化etcd
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Printf("load ini failed by map, error:%v\n", err)
		return
	}
	fmt.Printf("Load ini by map success\n")
	//初始化etcd
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("Init etcd error:%v\n", err)
		return
	}
	fmt.Printf("Init etcd success\n")

	//2从etcd中获取日志收集的配置项信息
	logEntitys, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("Get conf from etcd, error:%v\n", err)
		return
	}
	fmt.Printf("Get conf from etcd success :%v\n", logEntitys)
	for index, logEntity := range logEntitys {
		fmt.Printf("index: %v, logEntity: %v\n", index, logEntity)
	}
	// 循环每一个日志收集想，创建tailObj,发往kafka
	//未完待续

}

type user struct {
	key         string
	account     int
	modRevision int64
}

//事务
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("Init etcd error:%v\n", err)
		return
	}
	fmt.Printf("Init etcd success\n")

	userA := user{key: "userA"}
	userB := user{key: "userB"}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	respA, err := cli.Get(ctx, userA.key)
	cancel()
	if err != nil {
		fmt.Printf("获取数据失败 err:%v\n", err)
		return
	}
	if len(respA.Kvs) == 1 {
		fmt.Printf("respA.Kvs == 1\n")
		fmt.Printf("BEFOR: userA =%v\n", userA)
		userA.modRevision = respA.Kvs[0].ModRevision
		//实现strconv Package与基本数据类型的字符串表示之间的转换.Atoi()函数，
		//该函数等效于ParseInt(str string，base int，bitSize int)用于将字符串类型转换为int类型。
		userA.account, _ = strconv.Atoi(string(respA.Kvs[0].Value))
		fmt.Printf("AFTER: userA =%v\n", userA)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	respB, err := cli.Get(ctx, userB.key)
	cancel()
	if err != nil {
		fmt.Printf("获取数据失败 err:%v\n", err)
		return
	}
	if len(respB.Kvs) == 1 {
		fmt.Printf("respB.Kvs == 1\n")
		fmt.Printf("BEFOR: userB =%v\n", userB)
		userB.modRevision = respB.Kvs[0].ModRevision
		userB.account, _ = strconv.Atoi(string(respB.Kvs[0].Value))
		fmt.Printf("AFTER: userB =%v\n", userB)
	}

	userA.account -= 100
	userB.account += 100

	//ModRevision 是当前key最后一次修改的版本(针对单个key)
	//https://www.modb.pro/db/79479
	//Itoa() 函数，它相当于 FormatInt(int64(x), 10)。或者换句话说，Itoa() 函数在基数为 10 时返回 x 的字符串表示。
	txnReps, err := cli.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.ModRevision(userA.key), "=", userA.modRevision),
		clientv3.Compare(clientv3.ModRevision(userB.key), "=", userB.modRevision),
	).Then(
		clientv3.OpPut(userA.key, strconv.Itoa(userA.account)),
		clientv3.OpPut(userB.key, strconv.Itoa(userB.account)),
	).Else(
	//do something
	).Commit()
	if err != nil {
		log.Printf("事务执行失败 err:%v\n", err)
		return
	}
	if !txnReps.Succeeded {
		log.Println("事务失败")
		return
	}
	log.Println("事务成功")
}

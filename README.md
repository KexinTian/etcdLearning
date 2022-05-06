# etcdLearning
	1）module declares its path as: go.etcd.io/bbolt but was required as: github.com/coreos/bbolt
	解决办法：
	go mod edit -replace github.com/coreos/bbolt@v1.3.6=go.etcd.io/bbolt@v1.3.6
	

	2） google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.46.0), but does not contain package google.golang.org/grpc/naming
	解决办法：
	go mod edit -replace google.golang.org/grpc@v1.46.0=google.golang.org/grpc@v1.26.0
	
	3）使用 ini.MapTo(cfg, "./conf/config.ini") 将 .ini 文件中的配置信息，转化到cfg结构体中时报错。
	解决办法：
	命令行执行：go get github.com/go-ini/ini

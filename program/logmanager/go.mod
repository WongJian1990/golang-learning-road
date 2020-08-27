module logmanager

go 1.14

require (
	github.com/Shopify/sarama v1.27.0
	github.com/astaxie/beego v1.12.2
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.1.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/grpc v1.31.1 // indirect
	gopkg.in/olivere/elastic.v2 v2.0.61
)

replace github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

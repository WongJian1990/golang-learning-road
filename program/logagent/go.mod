module logagent

go 1.14

require (
	github.com/Shopify/sarama v1.27.0
	github.com/astaxie/beego v1.12.2
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.7 // indirect
	github.com/hpcloud/tail v1.0.0
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/smallnest/rpcx v0.0.0-20200822084631-5599102026f5 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	go.etcd.io/etcd v3.3.24+incompatible
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/net v0.0.0-20200528225125-3c3fba18258b
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/grpc v1.31.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

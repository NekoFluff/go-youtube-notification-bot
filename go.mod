module github.com/NekoFluff/hololive-livestream-notifier-go

// +heroku goVersion go1.23
go 1.23

// replace github.com/NekoFluff/discord => ../discord

require (
	github.com/NekoFluff/discord v1.10.1
	github.com/bwmarrin/discordgo v0.28.1
	github.com/dpup/gohubbub v0.0.0-20140517235056-2dc6969d22d8
	github.com/joho/godotenv v1.5.1
	github.com/robfig/cron v1.2.0
	github.com/urfave/cli/v2 v2.27.5
	go.mongodb.org/mongo-driver v1.17.1
	golang.org/x/text v0.19.0
	google.golang.org/api v0.204.0
)

require (
	cloud.google.com/go/auth v0.10.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.5 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.56.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
)

module github.com/NekoFluff/go-hololive-notification-bot

// +heroku goVersion go1.23
go 1.23

// replace github.com/NekoFluff/discord => ../discord

require (
	github.com/NekoFluff/discord v1.1.1
	github.com/bwmarrin/discordgo v0.28.1
	github.com/dpup/gohubbub v0.0.0-20140517235056-2dc6969d22d8
	github.com/joho/godotenv v1.5.1
	github.com/robfig/cron v1.2.0
	go.mongodb.org/mongo-driver v1.11.7
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
)

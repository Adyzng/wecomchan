set -ex

GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main
zip main.zip main
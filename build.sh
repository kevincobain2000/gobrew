#! /bin/sh

echo "building linux 64"
GOOS=linux GOARCH=amd64 go build cmd/gobrew/main.go && mv main bin/gobrew-linux-64
echo "building linux done"

echo "building darwin 64"
GOOS=darwin GOARCH=amd64 go build cmd/gobrew/main.go && mv main bin/gobrew-darwin-64
echo "building darwin done"

echo "building windows 64"
GOOS=windows GOARCH=amd64 go build cmd/gobrew/main.go && mv main.exe bin/gobrew-windows-64.exe
echo "building windows done"
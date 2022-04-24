#! /bin/sh

echo "building linux amd 64"
GOOS=linux GOARCH=amd64 go build cmd/gobrew/main.go && mv main bin/gobrew-linux-amd-64
echo "building linux amd 64 done"

echo "building linux arm 64"
GOOS=linux GOARCH=arm64 go build cmd/gobrew/main.go && mv main bin/gobrew-linux-arm-64
echo "building linux arm64  done"

echo "building darwin 64"
GOOS=darwin GOARCH=amd64 go build cmd/gobrew/main.go && mv main bin/gobrew-darwin-64
echo "building darwin done"

echo "building darwin arm-64 (m1)"
GOOS=darwin GOARCH=arm64 go build cmd/gobrew/main.go && mv main bin/gobrew-darwin-arm-64
echo "building darwin arm-64 (m1) done"

echo "building windows 64"
GOOS=windows GOARCH=amd64 go build cmd/gobrew/main.go && mv main.exe bin/gobrew-windows-64.exe
echo "building windows done"

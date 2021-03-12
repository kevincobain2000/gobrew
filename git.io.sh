#! /bin/sh

GOBREW_BIN_DIR=$HOME/.gobrew/bin
mkdir -p $GOBREW_BIN_DIR

GOBREW_ARCH_BIN=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        GOBREW_ARCH_BIN="gobrew-linux-64"
elif [[ "$OSTYPE" == "linux"*  ]]; then
        GOBREW_ARCH_BIN="gobrew-linux-64"        
elif [[ "$OSTYPE" == "darwin"* ]]; then
        GOBREW_ARCH_BIN="gobrew-darwin-64"
elif [[ "$OSTYPE" == "win64" ]]; then
        GOBREW_ARCH_BIN="gobrew-windows-64.exe"
# elif [[ "$OSTYPE" == "freebsd"* ]]; then
#         # ...
# else
#         # Unknown.
fi


curl -kLs https://raw.githubusercontent.com/kevincobain2000/gobrew/master/bin/$GOBREW_ARCH_BIN -o $GOBREW_BIN_DIR/gobrew

chmod +x $GOBREW_BIN_DIR/gobrew

echo "Installed successfully"


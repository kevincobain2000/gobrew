#! /bin/sh

GOBREW_BIN_DIR=$HOME/.gobrew/bin
mkdir -p $GOBREW_BIN_DIR

GOBREW_ARCH_BIN=''
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
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



curl -ks https://raw.githubusercontent.com/kevincobain2000/gobrew/bin/$GOBREW_ARCH_BIN -o $GOBREW_BIN_DIR/gobrew

chmod 755 $GOBREW_BIN_DIR/gobrew
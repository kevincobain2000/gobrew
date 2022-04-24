#!/bin/sh

GOBREW_BIN_DIR=$HOME/.gobrew/bin
mkdir -p $GOBREW_BIN_DIR

GOBREW_ARCH_BIN=''

THISOS=$(uname -s)
ARCH=$(uname -m)

case $THISOS in
   Linux*)
      case $ARCH in
        arm64)
          GOBREW_ARCH_BIN="gobrew-linux-arm-64"
          ;;
        aarch64)
          GOBREW_ARCH_BIN="gobrew-linux-arm-64"
          ;;
        *)
          GOBREW_ARCH_BIN="gobrew-linux-amd-64"
          ;;
      esac
      ;;
   Darwin*)
      case $ARCH in
        arm64)
          GOBREW_ARCH_BIN="gobrew-darwin-arm-64"
          ;;
        *)
          GOBREW_ARCH_BIN="gobrew-darwin-64"
          ;;
      esac
      ;;
   Windows*)
      GOBREW_ARCH_BIN="gobrew-windows-64.exe"
      ;;
esac

if [ -z "$GOBREW_VERSION" ]
then
      GOBREW_VERSION=master
      echo "Using gobrew version latest\n"
else
      echo "Using gobrew version $GOBREW_VERSION\n"
fi

curl -kLs https://raw.githubusercontent.com/kevincobain2000/gobrew/master/bin/$GOBREW_ARCH_BIN -o $GOBREW_BIN_DIR/gobrew

chmod +x $GOBREW_BIN_DIR/gobrew

echo "Installed successfully to: $GOBREW_BIN_DIR/gobrew"

echo "============================"
$GOBREW_BIN_DIR/gobrew help
echo "============================"

echo
echo "***Please add PATH below to your ~/.bashrc manually***"
echo
echo 'export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"'
echo
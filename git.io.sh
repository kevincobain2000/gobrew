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
          GOBREW_ARCH_BIN="gobrew-linux-arm64"
          ;;
        aarch64)
          GOBREW_ARCH_BIN="gobrew-linux-arm64"
          ;;
        armv6l)
          GOBREW_ARCH_BIN="gobrew-linux-arm_6"
          ;;
        armv7l)
          GOBREW_ARCH_BIN="gobrew-linux-arm_7"
          ;;
        *)
          GOBREW_ARCH_BIN="gobrew-linux-amd64"
          ;;
      esac
      ;;
   Darwin*)
      case $ARCH in
        arm64)
          GOBREW_ARCH_BIN="gobrew-darwin-arm64"
          ;;
        *)
          GOBREW_ARCH_BIN="gobrew-darwin-amd64"
          ;;
      esac
      ;;
   Windows*)
      GOBREW_ARCH_BIN="gobrew-windows-amd64.exe"
      ;;
esac

echo "Installing gobrew...\n"

curl -kLs https://github.com/kevincobain2000/gobrew/releases/latest/download/$GOBREW_ARCH_BIN -o $GOBREW_BIN_DIR/gobrew

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
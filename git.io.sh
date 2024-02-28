#!/bin/sh

GOBREW_BIN_DIR=$HOME/.gobrew/bin

# check if env GOBREW_ROOT is set
if [ -n "$GOBREW_ROOT" ]; then
  GOBREW_BIN_DIR=$GOBREW_ROOT/.gobrew/bin
fi

mkdir -p $GOBREW_BIN_DIR

GOBREW_ARCH_BIN=''
GOBREW_BIN='gobrew'
GOBREW_VERSION=${1:-latest}

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
          GOBREW_ARCH_BIN="gobrew-linux-arm"
          ;;
        armv7l)
          GOBREW_ARCH_BIN="gobrew-linux-arm"
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
   Windows|MINGW64_NT*)
      GOBREW_ARCH_BIN="gobrew-windows-amd64.exe"
      GOBREW_BIN="gobrew.exe"
      ;;
esac

if [ -z "$GOBREW_ARCH_BIN" ]; then
   echo "This script is not supported on $THISOS and $ARCH"
   exit 1
fi

DOWNLOAD_URL="https://github.com/kevincobain2000/gobrew/releases/download/$GOBREW_VERSION/$GOBREW_ARCH_BIN"
if [ "$GOBREW_VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/kevincobain2000/gobrew/releases/$GOBREW_VERSION/download/$GOBREW_ARCH_BIN"
fi

echo "Installing gobrew from: $DOWNLOAD_URL"
echo ""

curl -L --progress-bar $DOWNLOAD_URL -o $GOBREW_BIN_DIR/$GOBREW_BIN

chmod +x $GOBREW_BIN_DIR/$GOBREW_BIN

echo "Installed successfully to: $GOBREW_BIN_DIR/$GOBREW_BIN"

echo "============================"
$GOBREW_BIN_DIR/$GOBREW_BIN help
echo "============================"

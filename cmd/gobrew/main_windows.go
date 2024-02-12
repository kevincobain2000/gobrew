//go:build windows

package main

const usageMsg = `
    # Add gobrew to your environment variables
    PATH="%USERPROFILE%\.gobrew\current\bin;%USERPROFILE%\.gobrew\bin;%PATH%"
    GOROOT="%USERPROFILE%\.gobrew\current\go"

`

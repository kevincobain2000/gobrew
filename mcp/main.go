package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Gobrew MCP ☕",
		"1.0.0",
	)

	// Register tools
	s.AddTool(mcp.NewTool("install_gobrew",
		mcp.WithDescription("Install gobrew CLI tool based on your OS"),
	), installGobrewHandler)

	s.AddTool(mcp.NewTool("install_go_version",
		mcp.WithDescription("Install a specific Go version using gobrew"),
		mcp.WithString("version", mcp.Required(), mcp.Description("Go version to install")),
	), installGoHandler)

	s.AddTool(mcp.NewTool("use_go_version",
		mcp.WithDescription("Use a specific Go version"),
		mcp.WithString("version", mcp.Required(), mcp.Description("Go version to use")),
	), useGoHandler)

	s.AddTool(mcp.NewTool("uninstall_go_version",
		mcp.WithDescription("Uninstall a specific Go version"),
		mcp.WithString("version", mcp.Required(), mcp.Description("Go version to uninstall")),
	), uninstallGoHandler)

	s.AddTool(mcp.NewTool("list_versions",
		mcp.WithDescription("List installed Go versions"),
	), listHandler)

	s.AddTool(mcp.NewTool("list_remote_versions",
		mcp.WithDescription("List remote Go versions available"),
	), listRemoteHandler)

	s.AddTool(mcp.NewTool("gobrew_version",
		mcp.WithDescription("Show gobrew version"),
	), versionHandler)

	s.AddTool(mcp.NewTool("self_update_gobrew",
		mcp.WithDescription("Self-update gobrew"),
	), selfUpdateHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
func installGobrewHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("bash", "-c", `curl -sL https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | bash`)
	case "windows":
		cmd = exec.Command("powershell", "-Command", `Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))`)
	default:
		return nil, errors.New("unsupported OS")
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, errors.New("failed to install gobrew: " + stderr.String())
	}

	return mcp.NewToolResultText("✅ gobrew installed successfully! Add it to your PATH if needed."), nil
}

func installGoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrewWithVersion("install", request)
}

func useGoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrewWithVersion("use", request)
}

func uninstallGoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrewWithVersion("uninstall", request)
}

func listHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrew("list")
}

func listRemoteHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrew("ls-remote")
}

func versionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrew("version")
}

func selfUpdateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return runGobrew("self-update")
}
func runGobrew(command string) (*mcp.CallToolResult, error) {
	if _, err := exec.LookPath("gobrew"); err != nil {
		return nil, errors.New("gobrew is not installed or not in PATH")
	}

	cmd := exec.Command("gobrew", command)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, errors.New(strings.TrimSpace(errBuf.String()))
	}

	return mcp.NewToolResultText(strings.TrimSpace(out.String())), nil
}

func runGobrewWithVersion(command string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := exec.LookPath("gobrew"); err != nil {
		return nil, errors.New("gobrew is not installed or not in PATH")
	}

	version, ok := request.Params.Arguments["version"].(string)
	if !ok || version == "" {
		return nil, errors.New("version must be a non-empty string")
	}

	cmd := exec.Command("gobrew", command, version)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return nil, errors.New(strings.TrimSpace(errBuf.String()))
	}

	return mcp.NewToolResultText(strings.TrimSpace(out.String())), nil
}

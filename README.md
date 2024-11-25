<p align="center">
  <a href="https://github.com/kevincobain2000/gobrew">
    <img alt="gobrew" src="https://imgur.com/WkKYQPI.png" width="360">
  </a>
</p>
<p align="center">
  Go version manager, written in Go<br>
  Update and switch Go versions easily<br>
  Install Go on Linux or Mac (intel) or Mac with Apple chip (M1 to M4) or Windows
</p>

**Quick Setup:** One command to install Go and manage versions.

**Hassle Free:** Doesn't require root or sudo, or shell re-hash.

**Platform:** Supports (arm64, arch64, Mac, Mac M1, Ubuntu and Windows).

**Flexible:** Manage multiple Go versions including beta and rc.

**Colorful:** Colorful output.

![CI](https://github.com/kevincobain2000/gobrew/actions/workflows/build.yml/badge.svg)

## Install or update

### Step 1

**Using** curl (mac, linux) - recommended

```sh
curl -sL https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | bash
```

**Using** powershell (windows)

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))

```

**Using** go

```sh
go install github.com/kevincobain2000/gobrew/cmd/gobrew@latest
```

### Step 2

Now add `PATH` setting your shell config file (`.bashrc` or `.zshrc`).

 ```sh
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
```

**All DONE!**

Execute `gobrew` command from any dir.

```sh
gobrew
```

### Quick Usage

Simply use command `gobrew` from any dir. It will auto detect if Go version is set, or not latest, or not same as `go.mod` file.

<p align="center">
  <a href="https://github.com/kevincobain2000/gobrew">
    <img alt="gobrew command" src="https://imgur.com/vaqbS5o.png">
  </a>
</p>


### Full Usage


**Smart command**

```sh
gobrew
```

**Specific commands**

Will install and set Go

```sh
gobrew use 1.16
```

Will automatically install and set Go

```sh
gobrew use mod #from go.mod
gobrew use latest #latest stable
gobrew use dev-latest #latest of latest including rc|beta
```

Will only install it

```sh
gobrew install 1.16
gobrew use 1.16
```

Uninstall a version

```sh
gobrew uninstall 1.16
```

List installed versions `gobrew ls`

<p align="center">
  <a href="https://github.com/kevincobain2000/gobrew">
    <img alt="gobrew command" src="https://imgur.com/qCuaUMD.png">
  </a>
</p>



List available versions `gobrew ls-remote`

<p align="center">
  <a href="https://github.com/kevincobain2000/gobrew">
    <img alt="gobrew command" src="https://imgur.com/5FcBGUA.png">
  </a>
</p>


# All commands

```sh
╰─$ gobrew help

gobrew 1.10.7

Usage:

    gobrew use <version>           Install and set <version>
    gobrew ls                      Alias for list
    gobrew ls-remote               List remote versions (including rc|beta versions)

    gobrew install <version>       Only install <version> (binary from official or GOBREW_REGISTRY env)
    gobrew uninstall <version>     Uninstall <version>
    gobrew list                    List installed versions
    gobrew self-update             Self update this tool
    gobrew prune                   Uninstall all go versions except current version
    gobrew version                 Show gobrew version
    gobrew help                    Show this message

Examples:
    gobrew use 1.16                # use go version 1.16
    gobrew use 1.16.1              # use go version 1.16.1
    gobrew use 1.16rc1             # use go version 1.16rc1

    gobrew use 1.16@latest         # use go version latest of 1.16

    gobrew use 1.16@dev-latest     # use go version latest of 1.16, including rc and beta
                                   # Note: rc and beta become no longer latest upon major release

    gobrew use mod                 # use go version listed in the go.mod file
    gobrew use latest              # use go version latest available
    gobrew use dev-latest          # use go version latest avalable, including rc and beta

Installation Path:

# Add gobrew to your ~/.bashrc or ~/.zshrc
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
export GOROOT="$HOME/.gobrew/current/go"
```

# Uninstall gobrew

```sh
rm -rf $HOME/.gobrew
```

# Use it in Github Actions

For more details: https://github.com/kevincobain2000/action-gobrew

```yaml
on: [push]
name: CI
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.13, 1.18, 1.18@latest, 1.19beta1, 1.19@dev-latest, latest, dev-latest, mod]
        os: [ubuntu-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v2
      - uses: kevincobain2000/action-gobrew@v2
        with:
          version: ${{ matrix.go-version }}

      - name: Go
        run: go version
```

# Using Bash completions

```sh
curl https://raw.githubusercontent.com/kevincobain2000/gobrew/master/completions/bash/gobrew-completion > /usr/local/etc/bash_completion.d/gobrew

# or
curl https://raw.githubusercontent.com/kevincobain2000/gobrew/master/completions/bash/gobrew-completion >> ~/.zshrc

# or
curl https://raw.githubusercontent.com/kevincobain2000/gobrew/master/completions/bash/gobrew-completion >> ~/.bashrc
```

# Customization

By default, gobrew is installed in `$HOME` as `$HOME/.gobrew`.

You can change this by setting the `GOBREW_ROOT` environment variable.

```sh
# optionally set
echo "export GOBREW_ROOT=/usr/local/share" >> ~/.bashrc
# optionally set
echo "export GOBREW_ROOT=/usr/local/share" >> ~/.zshrc


#then
curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | sh
```

Set `GOROOT` and `GOPATH` in your shell config file (`.bashrc` or `.zshrc`).

```sh
# optionally set
export GOROOT="$HOME/.gobrew/current/go"
# optionally set
export GOPATH="$HOME/.gobrew/current/go"
```

# CHANGELOG

- **v1.2.0** - Added rc|beta versions, appended at the end of list
- **v1.5.0** - Mac M1 support
- **v1.5.5** - arm|M1|darwin support added
- **v1.5.8** - Show download progress and use Go's compression instead of tar command
- **v1.6.0** - Added support for @latest and @dev-latest and progress bar for download
- **v1.6.2** - Using goreleaser #35 by @juev
- **v1.6.3** - Added latest and dev-latest
- **v1.6.4** - Github action publish
- **v1.6.7** - Fixes rate limit issue
- **v1.7.4** - Added 2 new options `gobrew version` and `gobrew prune`
- **v1.7.5** - Fixes strange output on `gobrew use latest`
- **v1.7.8** - Windows support, self-update fixes
- **v1.7.9** - Windows fix ups and bash-completions
- **v1.8.0** - Windows support, including actions
- **v1.8.4** - Light background terminal support
- **v1.8.6** - Fixes where 1.20.0 was detected as 1.20
- **v1.9.0** - v1.8.6 ~ v1.9.0, updates colors packages, fixes UT issues for Github status codes
- **v1.9.4** - `gobrew` interactive
- **v1.9.8** - bug fix where 1.21 is not detected as 1.21.0
- **v1.10.10** - `ls-remote` is blazing fast, cached.
- **v1.10.11** - Optional options for cache and ttl.
- **v1.10.12** - Icons on `gobrew` and install command.


# DEVELOPMENT NOTES

```sh
go run ./cmd/gobrew -h
golangci-lint run ./...
go test -v ./...
```
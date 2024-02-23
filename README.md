<p align="center">
  <a href="https://github.com/kevincobain2000/gobrew">
    <img alt="gobrew" src="https://imgur.com/WkKYQPI.png" width="360">
  </a>
</p>
<p align="center">
  Go version manager, written in Go<br>
  Update and switch Go versions easily<br>
  Install Go on Linux or Mac (intel) or Mac with Apple chip (M1, M2 etc) or Windows
</p>

**Quick Setup:** One command to install Go and manage versions.

**Hassle Free:** Doesn't require root or sudo, or shell re-hash.

**Platform:** Supports (arm64, arch64, Mac, Mac M1, Ubuntu and Windows).

**Flexible:** Manage multiple Go versions including beta and rc.

**Colorful:** Colorful output.


# Build Status

![CI](https://github.com/kevincobain2000/gobrew/actions/workflows/build.yml/badge.svg)
![coverage](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=coverage&branch=master)
![build-time](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=build-time&branch=master)
![go-binary-size](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=go-binary-size&branch=master)
![unit-test-run-time](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=unit-test-run-time&branch=master)

![go-mod-dependencies](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=go-mod-dependencies&branch=master)
![go-sec-issues](https://coveritup.app/badge?org=kevincobain2000&repo=gobrew&type=go-sec-issues&branch=master)

![coverage](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=coverage&branch=master&output=svg&width=160&height=160&line=fill)
![build-time](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=build-time&branch=master&output=svg&width=160&height=160)
![go-binary-size](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=go-binary-size&branch=master&output=svg&width=160&height=160&line=fill)
![go-mod-dependencies](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=go-mod-dependencies&branch=master&output=svg&width=160&height=160)
![go-sec-issues](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=go-sec-issues&branch=master&output=svg&width=160&height=160&line=fill)
![unit-test-run-time](https://coveritup.app/chart?org=kevincobain2000&repo=gobrew&type=unit-test-run-time&branch=master&output=svg&width=160&height=160)





## Install or update

### Step 1

**Using** curl (mac, linux) - recommended

```sh
curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | sh
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
    <img alt="gobrew command" src="https://imgur.com/X6P3r8L.png">
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

List installed versions

```sh
╰─$ gobrew ls
1.15.1
1.16
1.16.3
1.17
1.17.1
1.17.2
1.17.3
1.17.4
1.17.5
1.17.6
1.17.7
1.17.8
1.18
1.18.1
1.18beta1
1.18rc1
1.22.0*
1.22rc1
1.22rc2

current: 1.20.0
```

List available versions

```sh
╰─$ gobrew ls-remote
==> [Info] Fetching remote versions
1.5	1.5.0  1.5.1  1.5.2  1.5.3  1.5.4
	1.5beta1  1.5beta2  1.5beta3  1.5rc1

1.6	1.6.0  1.6.1  1.6.2  1.6.3  1.6.4
	1.6beta1  1.6beta2  1.6rc1  1.6rc2

1.7	1.7.0  1.7.1  1.7.2  1.7.3  1.7.4
	1.7.5  1.7.6
	1.7beta1  1.7beta2  1.7rc1  1.7rc2  1.7rc3  1.7rc4
	1.7rc5  1.7rc6

1.8	1.8.0  1.8.1  1.8.2  1.8.3  1.8.4
	1.8.5  1.8.6  1.8.7
	1.8beta1  1.8beta2  1.8rc1  1.8rc2  1.8rc3  1.8.5rc4
	1.8.5rc5

1.9	1.9.0  1.9.1  1.9.2  1.9.3  1.9.4
	1.9.5  1.9.6  1.9.7
	1.9beta1  1.9beta2  1.9rc1  1.9rc2

1.10	1.10.0  1.10.1  1.10.2  1.10.3  1.10.4
	1.10.5  1.10.6  1.10.7  1.10.8
	1.10beta1  1.10beta2  1.10rc1  1.10rc2

1.11	1.11.0  1.11.1  1.11.2  1.11.3  1.11.4
	1.11.5  1.11.6  1.11.7  1.11.8  1.11.9  1.11.10
	1.11.11  1.11.12  1.11.13
	1.11beta1  1.11beta2  1.11beta3  1.11rc1  1.11rc2

1.12	1.12.0  1.12.1  1.12.2  1.12.3  1.12.4
	1.12.5  1.12.6  1.12.7  1.12.8  1.12.9  1.12.10
	1.12.11  1.12.12  1.12.13  1.12.14  1.12.15  1.12.16
	1.12.17
	1.12beta1  1.12beta2  1.12rc1

1.13	1.13.0  1.13.1  1.13.2  1.13.3  1.13.4
	1.13.5  1.13.6  1.13.7  1.13.8  1.13.9  1.13.10
	1.13.11  1.13.12  1.13.13  1.13.14  1.13.15
	1.13beta1  1.13rc1  1.13rc2

1.14	1.14.0  1.14.1  1.14.2  1.14.3  1.14.4
	1.14.5  1.14.6  1.14.7  1.14.8  1.14.9  1.14.10
	1.14.11  1.14.12  1.14.13  1.14.14  1.14.15
	1.14beta1  1.14rc1

1.15	1.15.0  1.15.1  1.15.2  1.15.3  1.15.4
	1.15.5  1.15.6  1.15.7  1.15.8  1.15.9  1.15.10
	1.15.11  1.15.12  1.15.13  1.15.14  1.15.15
	1.15beta1  1.15rc1  1.15rc2

1.16	1.16.0  1.16.1  1.16.2  1.16.3  1.16.4
	1.16.5  1.16.6  1.16.7  1.16.8  1.16.9  1.16.10
	1.16.11  1.16.12  1.16.13  1.16.14  1.16.15
	1.16beta1  1.16rc1

1.17	1.17.0  1.17.1  1.17.2  1.17.3  1.17.4
	1.17.5  1.17.6  1.17.7  1.17.8  1.17.9  1.17.10
	1.17.11  1.17.12  1.17.13
	1.17beta1  1.17rc1  1.17rc2

1.18	1.18.0  1.18.1  1.18.2  1.18.3  1.18.4
	1.18.5  1.18.6  1.18.7  1.18.8  1.18.9  1.18.10
	1.18beta1  1.18beta2  1.18rc1

1.19	1.19.0  1.19.1  1.19.2  1.19.3  1.19.4
	1.19.5  1.19.6  1.19.7  1.19.8  1.19.9  1.19.10
	1.19.11  1.19.12  1.19.13
	1.19beta1  1.19rc1  1.19rc2

1.20	1.20.0  1.20.1  1.20.2  1.20.3  1.20.4
	1.20.5  1.20.6  1.20.7  1.20.8  1.20.9  1.20.10
	1.20.11  1.20.12  1.20.13  1.20.14
	1.20rc1  1.20rc2  1.20rc3

1.21	1.21.0  1.21.1  1.21.2  1.21.3  1.21.4
	1.21.5  1.21.6  1.21.7
	1.21rc1  1.21rc2  1.21rc3  1.21rc4

1.22	1.22.0
	1.22rc1  1.22rc2
```

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

# Change Log

- v1.2.0 - Added rc|beta versions, appended at the end of list
- v1.4.1 - Go mod updated to 1.17
- v1.5.0 - Mac M1 support
- v1.5.1 - Oops
- v1.5.5 - arm|M1|darwin support added
- v1.5.6 - README updated
- v1.5.7 - Exit code fixed
- v1.5.8 - Show download progress and use Go's compression instead of tar command
- v1.5.9 - oops
- v1.6.0 - Added support for @latest and @dev-latest and progress bar for download
- v1.6.1 - Bug on use
- v1.6.2 - Using goreleaser #35 by @juev
- v1.6.3 - Added latest and dev-latest
- v1.6.4 - Github action publish
- v1.6.7 - Fixes rate limit issue
- v1.6.9 - Fixes #52, download error on status != 200
- v1.7.4 - Added 2 new options `gobrew version` and `gobrew prune`
- v1.7.5 - Fixes strange output on `gobrew use latest`
- v1.7.7 - Windows support?
- v1.7.8 - Windows support, self-update fixes
- v1.7.9 - Windows fix ups and bash-completions
- v1.8.0 - Windows support, including actions
- v1.8.2 - Prune old go versions, bug fix and `go.mod` auto detection
- v1.8.4 - Light background terminal support
- v1.8.6 - Fixes where 1.20.0 was detected as 1.20
- v1.9.0 - v1.8.6 ~ v1.9.0, updates colors packages, fixes UT issues for Github status codes
- v1.9.1 - Minor logging fixes
- v1.9.2 - Minor log message updated
- v1.9.4 - `gobrew` interactive
- v1.9.5 - bug
- v1.9.8 - bug fix where 1.21 is not detected as 1.21.0
- v1.10.0 - bug fixes for wrong gobrew version on windows
- v1.10.6 - fix for powershell
- v1.10.7 - fixes where `dev-latest` was not from `rc` or `beta`. `go get -u` deps updated and use go `1.22` in `go.mod`

# Build Status

| Branch  | Status                                                                                     |
| :------ | :----------------------------------------------------------------------------------------- |
| master  | ![Test](https://github.com/kevincobain2000/gobrew/workflows/Test/badge.svg?branch=master)  |
| develop | ![Test](https://github.com/kevincobain2000/gobrew/workflows/Test/badge.svg?branch=develop) |

# gobrew

Go version manager

## Install or update

With curl

```sh
$ curl -sLk https://git.io/gobrew | sh -
```

or install using go

```sh
$ go get -u github.com/kevincobain2000/gobrew/cmd/gobrew
```

Add `GOPATH` & `PATH` setting your shell config file (`.bashrc` or `.zshrc`).

 ```sh
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"

```

Reload config.

**All DONE!**

(optional)

```sh
export GOPATH="$HOME/go"
## or set for version specific
export GOPATH="$HOME/.gobrew/current/go"
```

### Confirm

```sh
$ gobrew help
```

### Usage

**Will install and set Go**

```sh
$ gobrew use 1.16
```

Will only install it and then use it

```sh
$ gobrew install 1.16
$ gobrew use 1.16
```

Uninstall a version

```sh
$ gobrew uninstall 1.16
```

List installed versions

```sh
$ gobrew ls

1.15.8
1.16*

current: 1.16
```

List available versions

```sh
$ gobrew ls-remote

...
1.15.1
1.15.2
1.15.3
1.15.4
1.15.5
1.15.6
1.15.7
1.15.8
...
1.16
1.16beta1
1.16rc1
```

# All commands

```sh
$ gobrew help                         Show this message
$ gobrew use <version>                Use <version>
$ gobrew install <version>            Download and install <version> (from binary))
$ gobrew uninstall <version>          Uninstall <version>
$ gobrew list                         List installed versions
$ gobrew ls                           Alias for list
$ gobrew ls-remote                    List remote versions
```

# Screenshots

![colors-ls-remote](https://i.imgur.com/gTBCfZL.png)
![colors-ls](https://i.imgur.com/KQbiuyH.png)


# Setup Go in Github actions

```yml
on: [push, pull_request]
name: Test
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.13, 1.14, 1.15, 1.16.7, 1.17rc1, 1.17rc2]
        os: [ubuntu-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
    - name: Set Env
      run: |
         echo "GOPATH=$HOME/.gobrew/current/go" >> $GITHUB_ENV
         echo "$HOME/.gobrew/bin" >> $GITHUB_PATH
         echo "$HOME/.gobrew/current/bin/" >> $GITHUB_PATH
    - name: Install Gobrew
      run: |
         curl -sLk https://git.io/gobrew | sh -
         gobrew use ${{ matrix.go-version }}
         go version

    - name: Checkout code
      uses: actions/checkout@v2

    - name: Go version
      run: go version
```

# Change Log

- v1.2.0 - Added rc|beta versions, appended at the end of list

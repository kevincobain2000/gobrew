# gobrew

Go version manager

## Install

Install with curl

```sh
$ curl https://raw.githubusercontent.com/kevincobain2000/gobrew/master/gobrew.sh | sh -
```

Add `PATH` setting your shell config file (`.bashrc` or `.zshrc`).

 ```sh
export GOPATH="$HOME/.gobrew/current/go:$PATH"
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
```

Reload config.
<span style="color:#03c03c">DONE!</span>.

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
$ gobrew ls-remote                    (not implemented yet) List remote versions
```

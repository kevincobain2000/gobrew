# Installation

### Step 1

```sh
curl https://raw.githubusercontent.com/kevincobain2000/gobrew/gobrew.sh | sh -
```
### Step 2

Add to your ``~/.bashrc``

 ```sh
export GOPATH="$HOME/.gobrew/current/go:$PATH"
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
 ```

### All done

```
gobrew use 1.16
gobrew ls
```

# Usage

 ```sh
 gobrew help
 ```

 ```sh
gobrew 1.0.0

Usage:
	gobrew help                         Show this message
	gobrew use <version>                Use <version>
	gobrew install <version>            Download and install <version> (from binary))
	gobrew uninstall <version>          Uninstall <version>
	gobrew list                         List installed versions
	gobrew ls                           Alias for list
	gobrew ls-remote                    (not implemented yet) List remote versions

Example:
	# install
	gobrew install 1.16
	gobrew install 1.15.8
	gobrew use 1.16

Reference:
	# Go versions
	https://golang.org/dl/
 ```

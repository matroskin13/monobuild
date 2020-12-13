# Monobuild

A simple tool for GO monorepo:

- Automatically analyze package dependencies (local files&go.mod) in monorepo, and report for changes (by git diff)

## Install

## Usage

Imagine we have a project with simple structure:

```bash
project
├── packageB
│   └── greeting.go
├── packageA
   └── cmd
       └── main.go
```
       
packageB/greeting.go:

```go
package packageB

func Hello() string {
	return "Hello world"
}
```

packageA/cmd/main.go:

```go
package main

import (
	"fmt"
	"project/packageB"
)

func main() {
	fmt.Println("ok", packageB.Hello())
}
```

### Changes command

Change the greeting "Hello world" to "Hello Earth", and run command:

```shell script
monobuild --module packageA/cmd changes

== Changes for package ==
dependency file: /goexample/packageB/greeting.go
``` 

The Monobuild automatically analyze your module for dependencies, and check dependencies for any version changes.

### Use config

Create a .monobuild.yml in root directory of monorepo:

```yaml
packages:
  packageA:
    entry: cmd
```

And run any command:

```shell script
monobuild changes

== Changes for packageA ==
dependency file: /Users/valentin/goexample/packageB/greeting.go
```
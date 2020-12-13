# Monobuild

A simple tool for GO monorepo:

- Analyze package dependencies (local files&go.mod) in monorepo, and report for changes (by git diff)
- Automatically push changed packages to docker registry
- (InProgress) Generate gitlab CI config
- (InProgress) Watch mode for development

## Install

## Build only changes packages

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

Package A depends on package B, and we want to build Package A when changes code of Package A and code of Package B. 

Monobuild will automatically detect the dependencies for package A by git diff, and build it.

First, make simple config for our application:

```yaml
packages:
  packageA:
    entry: cmd
    build:
      docker:
        image: "package-a"
```

and run command:

```shell script
$ monobuild build
$ docker run package-a
Hello world
```


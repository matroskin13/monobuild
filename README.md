# Monobuild

A simple tool for GO monorepo:

- Automatically analyze package dependencies (local files&go.mod) in monorepo, and report for changes (by git diff)
- Generate a template for your CI engine

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

### Use template for integration with CI

For example, let's integrate with Gitlab CI, and create .service-ci.yml:

```yaml
build:{{.serviceName}}:
  stage: build
  script:
    - echo {{.serviceName}}
```

And run generate template:

```shell script
monobuild changes  --use-template goexample/.service-ci.yml --out-template goexample/.full-ci.yml
```

As a result, Monobuild will generate a combined CI template including only the changed services:

```yaml
build:packageA:
  script:
    - echo packageA

build:packageC:
  script:
    - echo packageC
```

Let's finally integrate with Gitlab CI yml, change your gitlab-ci.yml:

```yaml
generate-config:
  stage: build
  script:
    - monobuild changes  --use-template goexample/.service-ci.yml --out-template goexample/.full-ci.yml
  artifacts:
    paths:
      - goexample/.full-ci.yml

trigger-changes:
  stage: build
  trigger:
    include:
      - artifact: goexample/.full-ci.yml
        job: generate-config
```

For more info read https://docs.gitlab.com/ee/ci/parent_child_pipelines.html#dynamic-child-pipelines
package build

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/matroskin13/monobuild/internal/tar"
	"io"
	"log"
)

var defaultGoDockerfile = `FROM golang:1.13 as builder

ARG SERVICE_PATH

WORKDIR /usr/src
COPY go.mod .
COPY go.sum .
RUN GOPROXY=${PROXY} go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./${SERVICE_PATH}

FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates tzdata

WORKDIR /usr/app
COPY --from=builder /usr/src/service .
CMD ["./service"]
`

type Docker struct {
	cli   *client.Client
	debug bool
}

func NewDocker() *Docker {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	return &Docker{cli: cli}
}

func (d *Docker) Build(ctx context.Context, dir string, packageEntry string, outputImage string) error {
	var buf bytes.Buffer

	additionalFile := []tar.CustomFile{
		{Name: "Dockerfile", Body: []byte(defaultGoDockerfile)},
	}

	if err := tar.CompressTarDir(dir, additionalFile, &buf); err != nil {
		return fmt.Errorf("cannot compress dir: %w", err)
	}

	buildArgs := make(map[string]*string)

	buildArgs["SERVICE_PATH"] = &packageEntry

	res, err := d.cli.ImageBuild(ctx, &buf, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{outputImage},
		BuildArgs:  buildArgs,
	})
	if err != nil {
		return fmt.Errorf("cannot build image: %w", err)
	}

	if d.debug {
		return writeToLog(res.Body)
	}

	defer res.Body.Close()

	return nil
}

func writeToLog(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		log.Println(string(n))
	}
	return nil
}

package env

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

func ParseTemplateWithEnv(in string) (string, error) {
	envMap, _ := envToMap()
	t := template.Must(template.New("tmpl").Parse(in))

	var buf bytes.Buffer

	if err := t.Execute(&buf, envMap); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func envToMap() (map[string]string, error) {
	envMap := make(map[string]string)
	var err error

	for _, v := range os.Environ() {
		split_v := strings.Split(v, "=")
		envMap[split_v[0]] = split_v[1]
	}

	return envMap, err
}

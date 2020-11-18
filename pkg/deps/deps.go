package deps

import (
	"golang.org/x/tools/go/packages"
	"strings"
)

func GetDepsAsFiles(inputPackagePath string) ([]string, error) {
	cfg := &packages.Config{Mode: packages.NeedSyntax | packages.NeedImports | packages.NeedName | packages.NeedFiles}
	pkgs, err := packages.Load(cfg, inputPackagePath)
	if err != nil {
		return nil, err
	}

	result := extractPackageFiles(pkgs[0])

	for _, imp := range pkgs[0].Imports {
		result = append(result, extractPackageFiles(imp)...)
	}

	return result, nil
}

func extractPackageFiles(pack *packages.Package) []string {
	var result []string

	for _, f := range pack.GoFiles {
		// TODO hack for go core files
		if !strings.Contains(f, "/go/src") {
			result = append(result, f)
		}
	}

	return result
}

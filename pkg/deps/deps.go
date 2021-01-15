package deps

import (
	"fmt"
	"github.com/matroskin13/monobuild/pkg/git"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"strings"
)

func GetDepsAsFiles(inputPackagePath string) ([]string, error) {
	pkg, err := loadFirstPackage(inputPackagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load packages: %w", err)
	}

	result := extractPackageFiles(pkg)

	for _, imp := range pkg.Imports {
		result = append(result, extractPackageFiles(imp)...)
	}

	return result, nil
}

func PackageChangeDeps(inputPackagePath string, revision string) (bool, []*modfile.Require, error) {
	pkg, err := loadFirstPackage(inputPackagePath)
	if err != nil {
		return false, nil, fmt.Errorf("cannot load packages: %w", err)
	}

	currentGoModFile, err := ioutil.ReadFile(pkg.Module.GoMod)
	if err != nil {
		return false, nil, err
	}

	// TODO need for git show (relative path)
	relativeGoModPath := strings.Replace(pkg.Module.GoMod, pkg.Module.Dir+"/", "", 1)

	oldGoModFile, err := git.GetOldFile(inputPackagePath, relativeGoModPath, revision)
	if err != nil {
		return false, nil, fmt.Errorf("cannot get old version of go.mod: %w", err)
	}

	currentGoMod, err := modfile.Parse("go.mod", currentGoModFile, nil)
	if err != nil {
		return false, nil, fmt.Errorf("cannot convert current file to go.mod: %w", err)
	}

	oldGoMod, err := modfile.Parse("go.mod", oldGoModFile, nil)
	if err != nil {
		return false, nil, fmt.Errorf("cannot convert old file to go.mod: %w", err)
	}

	diff := DiffRequire(currentGoMod.Require, oldGoMod.Require)

	return len(diff) > 0, diff, nil
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

func loadFirstPackage(inputPackagePath string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedModule,
		Dir:  inputPackagePath,
	}
	pkgs, err := packages.Load(cfg, inputPackagePath)
	if err != nil {
		return nil, err
	}

	return pkgs[0], nil
}

package deps

import (
	"golang.org/x/mod/modfile"
)

func DiffRequire(slice1, slice2 []*modfile.Require) []*modfile.Require {
	var diff []*modfile.Require

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1.Mod.String() == s2.Mod.String() {
					found = true
					break
				}
			}

			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

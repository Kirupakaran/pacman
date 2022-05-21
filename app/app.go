package app

import (
	"encoding/json"
	"log"
	"os"
)

type PackageDependencies struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func IsValidDir(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func Parse(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}

	contents, err := f.Readdirnames(0)
	repoPkgs := make(map[string]PackageDependencies)

	for _, subdir := range contents {
		data, err := os.ReadFile(dir + "/" + subdir + "/package.json")
		var pkgDeps PackageDependencies

		if err == nil {
			err := json.Unmarshal(data, &pkgDeps)
			if err != nil {
				log.Println("error parsing package.json for :", subdir, err)
			}
			repoPkgs[subdir] = pkgDeps
		}
	}
}

func extractCommonDeps(repoPkgs map[string]PackageDependencies) {
	commonDeps := make(map[string]string)
	commonDevDeps := make(map[string]string)

	for repo, pkgs := range repoPkgs {
	}

}

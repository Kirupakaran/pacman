package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Jeffail/gabs/v2"
)

type PackageDependencies struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Package struct {
	name     string
	versions map[string][]string
	isDev    bool
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

	ExtractPackages(repoPkgs)
}

func ExtractPackages(repoPkgs map[string]PackageDependencies) {
	allPkgs := make(map[string]Package)

	for repo, pkgs := range repoPkgs {
		log.Println("Extracting packages from repo : ", repo)
		log.Println("Number of dependencies : ", len(pkgs.Dependencies))
		log.Println("Number of dev dependencies : ", len(pkgs.DevDependencies))
		allPkgs = transform(allPkgs, repo, pkgs.Dependencies, false)
		allPkgs = transform(allPkgs, repo, pkgs.DevDependencies, true)
	}

	Prepare(allPkgs)
}

func Prepare(packages map[string]Package) {
	packument := gabs.New()
	aliasRegexp := regexp.MustCompile(`[.~^]`)
	for _, pkg := range packages {
		var outerKey string

		if pkg.isDev {
			outerKey = "devDependencies"
		} else {
			outerKey = "dependencies"
		}

		for version := range pkg.versions {
			if len(pkg.versions) == 1 {
				packument.Set(version, outerKey, pkg.name)
			} else {
				alias := pkg.name + aliasRegexp.ReplaceAllString(version, "")
				packument.Set("npm:"+pkg.name+"@"+version, outerKey, alias)
			}
		}
	}

	fmt.Println(packument)
}

func transform(packages map[string]Package, repo string, dependencies map[string]string, isDev bool) map[string]Package {
	for pkg, version := range dependencies {
		if existingPkg, exists := packages[pkg]; exists {
			// package already exists in common deps
			if _, exists := existingPkg.versions[version]; exists {
				// package and exact version exists; add repo to version list
				existingPkg.versions[version] = append(existingPkg.versions[version], repo)
			} else {
				// package exists but not the version; add new version info
				existingPkg.versions[version] = []string{repo}
			}
		} else {
			// add a new common dependency
			newPkg := new(Package)
			newPkg.name = pkg
			newPkg.versions = make(map[string][]string)
			newPkg.versions[version] = []string{repo}
			if isDev {
				newPkg.isDev = true
			}
			packages[pkg] = *newPkg
		}
	}

	return packages
}

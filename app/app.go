package app

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/Masterminds/semver/v3"
)

type PackageDependencies struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Package struct {
	Name     string
	Versions map[string][]string
	IsDev    bool
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

	extractPackages(repoPkgs)
}

func extractPackages(repoPkgs map[string]PackageDependencies) {
	allPkgs := make(map[string]Package)

	for repo, pkgs := range repoPkgs {
		log.Println("Extracting packages from repo : ", repo)
		log.Println("Number of dependencies : ", len(pkgs.Dependencies))
		log.Println("Number of dev dependencies : ", len(pkgs.DevDependencies))
		allPkgs = transform(allPkgs, repo, pkgs.Dependencies, false)
		allPkgs = transform(allPkgs, repo, pkgs.DevDependencies, true)
	}
	writeEncodedMapToFile(allPkgs)
	writePackagesWRepoToFile(allPkgs)
	writeBasePackageJsonToFile(allPkgs)
}

func writeEncodedMapToFile(packages map[string]Package) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	encodeErr := e.Encode(packages)
	if encodeErr != nil {
		log.Fatal("Failed to encode map to binary", encodeErr)
	}
	ioErr := os.WriteFile("packages.gob", b.Bytes(), 0666)
	if ioErr != nil {
		log.Fatal("Failed to write file", ioErr)
	}
}

func writePackagesWRepoToFile(packages map[string]Package) {
	pkgJson := gabs.New()
	pkgJson.Array("dependencies")
	pkgJson.Array("devDependencies")
	for _, pkg := range packages {
		jsonObj := gabs.New()
		for version, repos := range pkg.Versions {
			jsonObj.Set(repos, pkg.Name, version)
		}
		if pkg.IsDev {
			pkgJson.ArrayAppend(jsonObj, "devDependencies")
		} else {
			pkgJson.ArrayAppend(jsonObj, "dependencies")
		}
	}

	err := os.WriteFile("packages_list.json", pkgJson.Bytes(), 0666)
	if err != nil {
		log.Fatal("Failed to create packages_list.json")
	}
}

func writeBasePackageJsonToFile(packages map[string]Package) {
	pkgJson := gabs.New()
	aliasRegexp := regexp.MustCompile(`[.~^]`)
	for _, pkg := range packages {
		var outerKey string

		if pkg.IsDev {
			outerKey = "devDependencies"
		} else {
			outerKey = "dependencies"
		}

		for version := range pkg.Versions {
			if len(pkg.Versions) == 1 {
				pkgJson.Set(version, outerKey, pkg.Name)
			} else {
				alias := pkg.Name + ":" + aliasRegexp.ReplaceAllString(version, "")
				pkgJson.Set("npm:"+pkg.Name+"@"+version, outerKey, alias)
			}
		}
	}

	err := os.WriteFile("package.json", pkgJson.Bytes(), 0666)
	if err != nil {
		log.Fatal("Failed to create package.json")
	}
}

func transform(packages map[string]Package, repo string, dependencies map[string]string, isDev bool) map[string]Package {
	for pkg, version := range dependencies {
		// remove version modifiers
		digitRegexp := regexp.MustCompile(`^[0-9]+$`)
		if !digitRegexp.MatchString(version[0:1]) {
			version = version[1:]
		}
		if existingPkg, exists := packages[pkg]; exists {
			// package already exists in common deps
			if _, exists := existingPkg.Versions[version]; exists {
				// package and exact version exists; add repo to version list
				existingPkg.Versions[version] = append(existingPkg.Versions[version], repo)
			} else {
				// package exists but not the version; add new version info
				existingPkg.Versions[version] = []string{repo}
			}
		} else {
			// add a new common dependency
			newPkg := new(Package)
			newPkg.Name = pkg
			newPkg.Versions = make(map[string][]string)
			newPkg.Versions[version] = []string{repo}
			if isDev {
				newPkg.IsDev = true
			}
			packages[pkg] = *newPkg
		}
	}

	return packages
}

func Unify(isMinor bool) {
	var modifier string

	backupFiles()
	packages := readEncodedMapFromFile()

	if isMinor {
		modifier = "^"
	} else {
		modifier = "~"
	}

	for _, pkg := range packages {
		if len(pkg.Versions) == 1 {
			continue
		}

		versions := make([]*semver.Version, 0, len(pkg.Versions))
		for version := range pkg.Versions {
			v, err := semver.NewVersion(version)
			if err != nil {
				log.Printf("Failed to parse version %s for package %s\n", version, pkg.Name)
				continue
			}
			versions = append(versions, v)
		}

		sort.Sort(semver.Collection(versions))

		for i := 0; i < len(versions); i++ {
			c, err := semver.NewConstraint(modifier + versions[i].String())
			if err != nil {
				log.Printf("Failed to parse contraint %s\n", versions[i])
				continue
			}
			for j := i + 1; j < len(versions); j++ {
				if valid, _ := c.Validate(versions[j]); valid {
					log.Printf("Found similar version %s to %s for package %s\n", versions[j], versions[i], pkg.Name)
					copy(pkg.Versions[versions[j].String()], pkg.Versions[versions[i].String()])
					delete(pkg.Versions, versions[i].String())
					i += 1
				}
			}
		}
	}

	writeEncodedMapToFile(packages)
	writePackagesWRepoToFile(packages)
	writeBasePackageJsonToFile(packages)
}

func readEncodedMapFromFile() map[string]Package {
	data, ioErr := os.ReadFile("packages.gob")
	if os.IsNotExist(ioErr) {
		log.Fatal("packages.gob file missing. Run parse to create it")
	}

	var decodedMap map[string]Package
	r := bytes.NewReader(data)
	d := gob.NewDecoder(r)

	decodeErr := d.Decode(&decodedMap)
	if decodeErr != nil {
		log.Fatal("Failed to decode packages.gob", decodeErr)
	}

	return decodedMap
}

func backupFiles() {
	file, err := os.Stat("package.json")
	if err == nil && file != nil {
		os.Rename("package.json", "package.json"+"_"+time.Now().Format("20060102150405"))
	}

	file, err = os.Stat("packages_list.json")
	if err == nil && file != nil {
		os.Rename("packages_list.json", "packages_list.json"+"_"+time.Now().Format("20060102150405"))
	}
}

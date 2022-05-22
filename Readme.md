# pacman
Stealin' from Arch, until I find a better one :)

Manage your node microservices dependencies easily.

## Commands

*pacman parse <dir>*
Given a directory path where sub-directories are multiple node repos,  this command will parse package.json from all the repos and create a unified package.json in the root level with list of all dependencies and dev dependencies.

An example directory structure:
```
npm
|__ dumbledore
|__ npm-auth-ws
|__ registry-frontdoor
|__ user-acl-two
|__ wubwub
```

This command will create a package.json with all the dependencies from dumbledore, npm-auth-ws, registry-frontdoor, user-acl-two, wubwub. Now that we have a common package.json, it can be used to analyze, update, add, remove dependencies easily. Also committing this package.json in a repo will give dependabot alerts for all the packages in a single place.

Note that, repos may have different versions of dependencies and this command will create dependencies with alias to manage them easily. For example, if user-acl-two has `lodash: "4.17.19` and wubwub has `lodash: "4.21.0`, the common package.json will have
```
"lodash:41719": "npm:lodash@4.17.19",
"lodash:4210": "npm:lodash@4.21.0",
```

*pacman unify --minor (optional)*
This command will try to unify the dependencies with multiple versions. By default, it will unify to a common patch version. If `--minor` is passed, it'll unify to a minor version.

For example,
```
"lodash:41719": "npm:lodash@4.17.1",
"lodash:41719": "npm:lodash@4.17.19",
```

will be unified to 
```
"lodash:41719": "npm:lodash@4.17.19",
```

*pacman update <repo path>*

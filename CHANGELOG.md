## 0.8.3 (2025-05-13)

### Fix

- **internal/pkg/cmd/get.go**: Enhance output handling for structured formats
- **internal/pkg/cmd/get.go**: Enhance output handling for structured formats
- **internal/pkg/cmd/get.go**: Enhance output handling for structured formats

## 0.8.2 (2025-03-12)

### Fix

- **.goreleaser.yaml**: add checksum

## 0.8.1 (2024-11-21)

### Fix

- **.goreleaser.yaml, Makefile**: Update build configuration for go-mod and version export locations

## 0.8.0 (2024-11-21)

### Feat

- change tag for alias and change gittag format (#65)
- add flag to init (#60)
- actualizar readmemd (#62)

### Refactor

- **internal/pkg/cmd/bump.go**: Update pre-bump and post-bump scripts execution flow (#66)
- **docs**: Improve examples in README and CLI (#64)
- improve message of launching hooks (#59)

## 0.8.0-a3 (2024-11-12)

### Fix

- **conventionalcommits**: Update Commit Determination Functionality

### Refactor

- **internal/conventionalcommits/conventionalcommits.go**: Update common change determination logic
- **internal/conventionalcommits/conventionalcommits.go**: Simplify package using new parser structure of go-conventionalcommits

## 0.8.0-a2 (2024-11-11)

## 0.8.0-a1 (2024-11-08)

## 0.8.0-a0 (2024-11-08)

### Feat

- hooks (#51)
- **Makefile**: Introduce version bumping options and improve release process (#47)

### Fix

- output (#46)
- **internal/config/config.go**: Filter config versions by specified fields (#45)

### Refactor

- cmd retrieveCommandFlags (#44)
- **Makefile**: Rearrange build targets, use CURRENT_VERSION for go application
- **internal/config/config.go**: Introduce field filtering in printConfigVersionsPlain function (#42)
- **internal/config/config.go**: Introduce field filtering in printConfigVersionsPlain function (#42)

## 0.7.0 (2024-11-04)

### Feat

- project information

## 0.6.4 (2024-10-30)

### Fix

- **.goreleaser.yaml**: Update GoReleaser configuration - Set owner, adjust builds, and add Docker images.

## 0.6.3 (2024-10-30)

### Refactor

- **.goreleaser.yaml**: Update build process using gommitizen, remove ai-commit, create Docker images and modify build retrieveCommandFlags

## 0.6.2 (2024-10-30)

### Refactor

- **.goreleaser.yaml**: Update GitHub repository owner, set CGO_ENABLED to 0 for all builds, correct version export in gommitizen binary, update Docker image templates and build retrieveCommandFlags

## 0.6.1 (2024-10-30)

### Refactor

- **myFile.go**: Refactor function for improved readability and maintainability
- **go.mod**: Update module path to github.com/freepik-company/gommitizen

## 0.6.0 (2024-10-30)

## 0.5.2 (2024-09-03)

## 0.5.1 (2024-03-04)

### Fix

- Refactor version handling

## 0.5.0 (2024-01-10)

## 0.4.13 (2024-01-10)

## 0.4.12 (2024-01-10)

## 0.4.11 (2024-01-10)

## 0.4.10 (2024-01-10)

## 0.4.9 (2024-01-10)

## 0.4.8 (2024-01-10)

## 0.4.7 (2024-01-10)

## 0.4.6 (2024-01-09)

## 0.4.5 (2024-01-09)

## 0.4.4 (2024-01-09)

## 0.4.3 (2024-01-09)

## 0.4.2 (2024-01-09)

### Fix

- show message to run gommitizen in the root path of git

## 0.4.1 (2024-01-09)

### Fix

- update commit message and tag when it's root When gommitizen it's initialized in the root of the git repository the tag message is updated to avoid concatenate a dot at the end of the tag string and set the commit to a more legible string
- show a best error message It shows how you could fix the problem by teaching you the requierements of gommitizen about to the minimun number of commits it needs

## 0.4.0 (2024-01-08)

### Feat

- add version command

### Fix

- update version

## 0.3.2 (2023-12-05)

### Fix

- fix unit tests in pipeline

## 0.3.1 (2023-12-05)

### Fix

- fix unit tests in pipeline

## 0.3.0 (2023-12-05)

### Feat

- add new dependencies for testing
- add test target to makefile
- git interface for mocking
- unit tests
- mock for git object

### Refactor

- rename git object to gitHandler to avoid collides with package name
- delete nor needed file
- update code for unit testing

## 0.2.3 (2023-11-29)

### Fix

- **GH-Actions**: change event to respond to

## 0.2.2 (2023-11-29)

### Fix

- **GH-Actions**: change event to listen on release

## 0.2.1 (2023-11-29)

### Fix

- **GH-Action**: change event on which the event is triggered

## 0.2.0 (2023-11-29)

### Feat

- add prefix to config (#16)
- add scan target to Makefile to perform a scan with sonnarqube and lint
- **version.go**: add support for manual increment of version
- support for multiples version files Example of commandUsage: file .version.json ``` {   "version": "0.12.0",   "commit": "37149b77ce36d40422441ad20b999156ba107e18",   "version_files": [     "Chart.yaml:version",     "other-version.txt:version"   ] } ```
- add changelog refactor version to use changelog and split update function into two to allow for more flexibility and to allow for the changelog to be updated. Translates comments and literals to english.
- add support of changelog and increment in bump command
- add cobra to handle command input and some extra features like autogenerated help or cli completion

### Fix

- add filePath to version while LoadVersionData
- fix adding commit message to the list of changes Also translate comments and literals to English
- fix initialization with HEAD^
- fix git diff comment when serveral files changed

### Refactor

- **version**: split the code in several files for readability
- improve text output
- create a version object from constructor
- **cmd/bump.go**: replace Printf with Errorf when errors

## 0.1.3 (2023-11-24)

### Fix

- **GH-Action**: change binary name

## 0.1.2 (2023-11-24)

### Fix

- **GH-Actions**: fix a typo in the copy process

## 0.1.1 (2023-11-24)

### Fix

- **GH-Actions**: add release autogeneration

## 0.1.0 (2023-11-24)

### Fix

- **GH-Actions**: add release autogeneration

## 0.1.0 (2023-11-24)

### Fix

- type of increment related to its priority

### Refactor

- **main.go**: fix text for help
- fix name of application
- rename application

# Gommitizen

Gommitizen is a command-line tool that helps manage the versioning of a software project. This tool is able to manage serveral projects in a same repository with their different versions each.

The tool analyzes commit messages in a Git repository, looking for certain prefixes in the messages to determine the type of changes made. The prefixes include "BREAKING CHANGE:", "feat:", and "fix:", between others, which likely refer to breaking changes, new features, and bug fixes, respectively.

Version information is stored in a `VersionData` structure, which includes the current version, the current commit, and a list of version files. This information can be used to generate a changelog, determine the next version of the software, or perform other tasks related to version management.

## Compilation

To compile the project, run the following command:

```bash
make build
```

This will generate a binary in the `bin/` directory.

## Docker

You can build a Docker image of Gommitizen with the following command:

```bash
make docker
```

After building the image, you can run it with:

```bash
docker run -it gommitizen:<tag> help
```

Replace `<tag>` with the tag of the image you have built.

## Installation

To install Gommitizen in your `$GOPATH`, first build the project with `make build`, then run:

```bash
make install
```

_*TODO*: Hombrew install_

## Code Analysis with SonarQube

To start a code analysis with SonarQube, run:

```bash
make scan
```

Make sure to have SonarQube running before starting the analysis.

There is another way to start a code analysis with SonarQube, run the target `start-sonar` in the `Makefile`.

If you want to stop the SonarQube server, run the target `stop-sonar` in the `Makefile`.

## Testing

To run the tests, run:

```bash
make test
```

There is an alias to run the tests with `tests`.

### Coverage

To run the tests with coverage, run:

```bash
make coverage
```

### Development

If you want to regenerate the mocks, run:

```bash
make mocks
```

## Usage

To use Gommitizen, run:

```bash
gommitizen <command> [flags]
```

### Commands

The following commands are available:

- `bump`: Bumps the version of a project.
- `completion`: Generates the completion script for the specified shell.
- `help`: Shows the help message.
- `init`: Initializes the versioning of a project.

#### Global flags

The following flags are available for all commands:

- `-h` or `--help`: Shows the help message.
- `-t` or `--toggle`: Help message for toggle.

#### Init flags

The following flags are available for the `init` command:

- `-d` or `--directory`: The directory where the project is located. If not specified, the current directory is used.
- `-p` or `--prefix`: The prefix of the tag message. If not specified, an empty prefix is used.

#### Bump flags

The following flags are available for the `bump` command:

- `-d` or `--directory`: The directory where the project is located. If not specified, it scans all the directories in the current directory to look for projects with a .version file.
- `-c` or `--changelog`: It generates a changelog with the changes made since the last version.
- `-i` or `--increment`: The type of increment to make. It can be `major`, `minor`, or `patch`. If it is specified the automatic detection of version is not run.  

## Types of commits

There are two types of commits: version commits and regular commits. Version commits are those that change the version of the software, while regular commits are those that do not change the version.

Those that change the version of the software are those that have a commit message with a prefix that indicates the type of change. The prefixes are the following:

- `BREAKING CHANGE:` or `bc`: Indicates a breaking change in the software.
- `feat:`: Indicates a new feature in the software.
- `fix:`: Indicates a bug fix in the software.

Those that do not change the version of the software are those that have a commit message with a prefix that indicates the type of change. The prefixes are the following:

- `perf:`: Indicates a performance improvement in the software.
- `refactor:`: Indicates a code refactoring in the software.
- `docs:`: Indicates a documentation change in the software.
- `test:`: Indicates a test change in the software.
- `chore:`: Indicates a change in the build process or auxiliary tools in the software.
- `ci:`: Indicates a change in the CI configuration files and scripts in the software.
- `style:`: Indicates a change in the style of the code in the software.
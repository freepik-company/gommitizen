# Gommitizen

Gommitizen is a command-line tool that helps manage the versioning of a software project. This tool is able to manage serveral projects in a same repository with their different versions each.

The tool analyzes commit messages in a Git repository, looking for certain prefixes in the messages to determine the type of changes made. The prefixes include "BREAKING CHANGE:", "feat:", and "fix:", which likely refer to breaking changes, new features, and bug fixes, respectively.

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

`TODO: Hombrew install`

## Code Analysis with SonarQube

To start a code analysis with SonarQube, run:

```bash
make scan
```

Make sure to have SonarQube running before starting the analysis.

## Command line utility

TODO
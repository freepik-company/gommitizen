# Gommitizen

Gommitizen is a command-line tool that helps manage the versioning of a software project. This tool is able to manage serveral projects in a same repository with their different versions each.

The tool analyzes commit messages in a Git repository, looking for certain prefixes in the messages to determine the type of changes made. The prefixes include "BREAKING CHANGE:", "feat:", and "fix:", between others, which likely refer to breaking changes, new features, and bug fixes, respectively.

Version information is stored in a `VersionData` structure, which includes the current version, the current commit, and a list of version files. This information can be used to generate a changelog, determine the next version of the software, or perform other tasks related to version management.

## Installation

To install **gommitizen**, run the following command:
```bash 
curl -s https://raw.githubusercontent.com/freepikcompany/gommitizen/main/scripts/get-gommitizen.sh | sudo bash
```

This script will download the latest release of the **gommitizen** binary and install it in the `/usr/local/bin` directory.

To verify the installation, run the following command:

```bash
gommitizen --version
```

You should see the version of **gommitizen** that you installed.

## Compilation

To compile the project, run the following command:

```bash
make bin
```

This will generate a binary in the `bin/` directory.

```bash
make install 
```

This will install the binary in the `/usr/local/bin` directory.

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

### Docker

To run Gommitizen in a Docker container, run:

```bash
docker run --rm \
  -e GIT_USER_NAME=user.name \
  -e GIT_USER_EMAIL=user@email \
  -v $(pwd):/source \
  ghcr.io/freepikcompany/gommitizen:<tag> [flags]
```

Replace  `<tag>` with the tag of the image you want to use. Select the command and flags you want to use.

### Examples

#### Init

To initialize the versioning of a project, run:

```bash
gommitizen init -d <directory> -p <prefix>
```

This will create a `.version.json` file in the given directory with the version `0.0.0`.

#### Bump

To bump the version of a project, run:

```bash
gommitizen bump
```

This will bump the version of all projects in the current directory.

If you want to bump the version of a specific project, run:

```bash
gommitizen bump -d <directory>
```

This will bump the version of the project in the given directory.

if you want to bump the version of projects and generate a changelog, run:

```bash
gommitizen bump -c
```

This will bump the version of the projects and generate a changelog with the changes made since the last version.

If you want to bump the version of project to a major version, run:

```bash
gommitizen bump -i major
```

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

## Version files structure

Each project in the monorepo has a `.version.json` file that contains the version of the software.

The version files are structured as follows:

```json
{
    "version": "0.18.1",
    "commit": "72929b90547b8527e22e402b6784e0c7f5812428",
    "version_files": [
        "Chart.yaml:version",
        "other-version.txt:version",
        "a-file-that-need-a-regex.txt:^version=([0-9]+\\.[0-9]+\\.[0-9]+)$"
    ],
    "prefix": ""
}
```

The `version` field contains the current version of the software. The `commit` field contains the commit where the version was changed. The `version_files` field contains the list of files that contain the version of the software and the bump process will upgrade too. The `prefix` field contains the prefix of the tag message that changed the version of the software.

The `version` and `commit` fields are managed by Gommitizen. The `version_files` and `prefix` fields are managed by the user.

`version_files` is a list of strings. Each string contains the path of the file and the name of the variable that contains the version. The path and the name of the variable are separated by a colon (`:`). The path is relative to the root of the project. Tha name of the variable can be replace by a regular expression to find the version in the file (remember to scape the special characters and group the version part of the expression with parentheses like in the example).
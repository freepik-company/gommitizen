# Gommitizen

Gommitizen is a command-line tool that helps manage the versioning of a software project. This tool is 
able to manage serveral projects in a same repository with their different versions each. It supports the conventional 
commits specification (https://www.conventionalcommits.org/en/v1.0.0/) to determine the increment of the version for 
each project.

## Installation

To install **gommitizen**, run the following command:
```bash 
curl -s https://raw.githubusercontent.com/freepik-company/gommitizen/main/scripts/get-gommitizen.sh | sudo bash
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

Next a list of the available commands and their description:
- `bump`: Make a version bump

- `get`: Give a list of projects, their versions and other information

- `init`: Start a repository to use gommitizen

### bump command

Increment the version of the project according to the conventional commits specification.

**Flags:**
- `-c`, `--changelog`: generate the changelog for the newest version

- `-i`, `--increment`: manually specify the desired increment {MAJOR, MINOR, PATCH}



**Examples of usage:**

```shell
# To bump the version of a project, run:
gommitizen bump
# This will bump the version of all projects in the current directory.

# If you want to bump the version of a specific project, run:
gommitizen bump -d <directory>
# This will bump the version of the project in the given directory.

# If you want to bump the version of projects and generate a changelog, run:
gommitizen bump -c
# This will bump the version of the projects and generate a changelog with the changes made since the last version.

# If you want to bump the version of project to a major version, run:
gommitizen bump -i MAJOR

```





### get command

Show information about the projects in the repository. It can show the version, the prefix, the commit 
information and all the information saved in the config file.

**Flags:**
- `-o`, `--output`: Select the output format {json, yaml, plain}

- `-p`, `--prefix`: A prefix to look for a project to show information



**Examples of usage:**

```shell
# To show all information in yaml format, run:
gommitizen get all -o yaml
# To show the version of the projects in plain format, run:
gommitizen get version -o plain
# or just:
gommitizen get version

```



**Subcommands:**



- `all`: Get all the information of the projects in the repository. It will show the version, the prefix, the commit
information and all the information saved in the config file.

- `commit`: Get the commit information of the projects in the repository.

- `prefix`: Get the prefix of the projects in the repository.

- `version`: Get the version of the projects in the repository. It will show the version of the projects and the prefix.

### init command

Initialize the repository to use gommitizen. It will create a file with the version of the project and 
the first commit of the project.

**Flags:**
- `-`, `--bump-changelog`: Update changelog on bump

- `-p`, `--prefix`: Select a prefix for the version file



**Examples of usage:**

```shell
# To initialize the versioning of a project, run: 
gommitizen init -d <directory> -p <prefix>`
# This will create a .version.json file in the given directory with the version 0.0.0.
```






### Docker

To run Gommitizen in a Docker container, run:

```bash
docker run --rm \
  -e GIT_USER_NAME=user.name \
  -e GIT_USER_EMAIL=user@email \
  -v $(pwd):/source \
  ghcr.io/freepik-company/gommitizen:<tag> [retrieveCommandFlags]
```

Replace  `<tag>` with the tag of the image you want to use. Select the command and retrieveCommandFlags you want to use.

Example:
```bash
docker run --rm \
  -e GIT_USER_NAME=user.name \
  -e GIT_USER_EMAIL=user@email \
  -v $(pwd):/source \
  ghcr.io/freepik-company/gommitizen:latest [retrieveCommandFlags]
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
    "prefix": "my_prj"
}
```

The `version` field contains the current version of the software. The `commit` field contains the commit where the version was changed. The `version_files` field contains the list of files that contain the version of the software and the bump process will upgrade too. The `prefix` field contains the prefix of the tag message that changed the version of the software.

The `version` and `commit` fields are managed by Gommitizen. The `version_files` and `prefix` fields are managed by the user.

`version_files` is a list of strings. Each string contains the path of the file and the name of the variable that contains the version. The path and the name of the variable are separated by a colon (`:`). The path is relative to the root of the project. Tha name of the variable can be replace by a regular expression to find the version in the file (remember to scape the special characters and group the version part of the expression with parentheses like in the example).

### Hooks

Example:

```json
{
    "version": "0.18.1",
    "commit": "72929b90547b8527e22e402b6784e0c7f5812428",
    "version_files": [
        "Chart.yaml:version",
        "other-version.txt:version",
        "a-file-that-need-a-regex.txt:^version=([0-9]+\\.[0-9]+\\.[0-9]+)$"
    ],
    "prefix": "my-prj",
    "hooks": {
        "pre_bump": "echo 'pre-bump hook'",
        "post_bump": "echo 'post-bump hook'",
        "post_changelog": "echo 'post-changelog hook'",
        "pre_changelog": "echo 'pre-changelog hook'"
    }
}
```

There are four hooks available:

- `pre_bump`: Runs before the bump process.
- `post_bump`: Runs after the bump process.
- `pre_changelog`: Runs before the changelog generation.
- `post_changelog`: Runs after the changelog generation.

The hooks are shell getCommands that are executed in the root of the project. These are all optional fields.

## Development

To run the project in development mode, run:

```bash
go run ./cmd/gommitizen/main.go
```

To run a new release locally, run:

```bash
make release
```

If you want to run the release in pipeline, run:

```bash
make bump
```

to bump the version of the project and changelog. Then push the changes and tag to the repository to trigger the pipeline. That will generate the release and publish the binaries and docker image.

If you want to increase the version manually, run:

```bash
cz bump --increment (MAJOR|MINOR|PATCH) --changelog
```

Then push the tag to the repository:

```bash
git push && git push --tags
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

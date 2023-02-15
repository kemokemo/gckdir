# gckdir (go check directory tool)

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![test-and-build](https://github.com/kemokemo/gckdir/actions/workflows/test-and-build.yml/badge.svg)](https://github.com/kemokemo/gckdir/actions/workflows/test-and-build.yml)

This is a CLI tool that compares the structure of two directories and the hash values of individual files to verify their identity.

The hash list can be output in json format to check the contents of a directory deployed on another PC, or the check results can be output to a report in html format.

## Install

### Homebrew

```sh
brew install kemokemo/tap/gckdir
```

### Scoop

First, add my scoop-bucket.

```sh
scoop bucket add kemokemo-bucket https://github.com/kemokemo/scoop-bucket.git
```

Next, install this app by running the following.

```sh
scoop install gckdir
```

### Binary

Get the latest version from [the release page](https://github.com/kemokemo/gckdir/releases/latest), and download the archive file for your operating system/architecture. Unpack the archive, and put the binary somewhere in your `$PATH`.

## Usage

### Generate

This is a way to generate a hash list. Stores directory structure and hash values for individual files.

```bash
$ gckdir generate path/to/source_directory hash_name.json
```

For more details, please see `gckdir generate --help`.

### Verify

The following is an example of comparing a deployed directory with a pre-generated hash list.

```bash
$ gckdir verify hash_name.json path/to/target_directory
```

Direct directory-to-directory comparisons are also possible.

```bash
$ gckdir verify path/to/source_directory path/to/target_directory
```

For more details, please see `gckdir verify --help`.

#### Create a result report

![verification_report](./images/verification_report.png)

You can create a verification result report with `--report` or `-r` option.
The report file name will be in the format `Result_{YYYYYMMDD}-{hhmmss}. {nanosecond}.html`.

```bash
$ gckdir verify --report hash_name.json path/to/target_directory
```

To specify the file name of the report file, use the `--output` or `-o` option as follows.

```bash
$ gckdir verify --report --output output_name.html hash_name.json path/to/target_directory
```

#### Open a result report after creating

![open_animation](./images/open_animation.gif)

If you want to check the result immediately on the browser, please use the `--open` or `-p` option. This option includes the `--report` option.

```bash
$ gckdir verify --open hash_name.json path/to/target_directory
```

#### Verify only the structure of files and directories

Use the `--no-hv` or `-nh` option if you do not care about the hash value and only want to check the placement of files and folders.

```bash
$ gckdir verify --open --no-hv hash_name.json path/to/target_directory
```

#### Ignore files of other software

To ignore the presence or absence of a particular file, use the `--no-uv` or `-nu` option.

```bash
$ gckdir verify --report --no-uv hash_name.json path/to/target_directory
```

## Contribution

Please feel free to send me a pull request. :smile:

1. Fork ([https://github.com/kemokemo/gckdir/fork](https://github.com/kemokemo/gckdir/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Run `golangci-lint run` and confirm that it passes.
1. Create a new Pull Request

## Author

[kemokemo](https://github.com/kemokemo)

## License

[MIT](https://github.com/kemokemo/gckdir/blob/main/LICENSE)

## Special Thanks

This application uses the following excellent projects.

* [github.com/tcnksm/gcli](https://github.com/tcnksm/gcli) - [MIT](https://github.com/tcnksm/gcli/blob/master/LICENSE)
* [github.com/urfave/cli](https://github.com/urfave/cli) - [MIT](https://github.com/urfave/cli/blob/master/LICENSE)
* [github.com/ahmetb/go-linq](https://github.com/ahmetb/go-linq) - [Apache License 2.0](https://github.com/ahmetb/go-linq/blob/master/LICENSE)
* [github.com/skratchdot/open-golang](https://github.com/skratchdot/open-golang) - [MIT](https://github.com/skratchdot/open-golang/blob/master/LICENSE-MIT)
* [github.com/twbs/bootstrap](https://github.com/twbs/bootstrap) - [MIT](https://github.com/twbs/bootstrap/blob/master/LICENSE)

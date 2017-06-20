# gckdir

This is go-check directory tool.

## Description

You can generate a hash list of specified directory in the json format.  
With the hash list, you can compare the target directory's structure and each hash value
of files.

## Usage

### Generate

```bash
$ gckdir generate path/to/source_directory hash_name.json
```

For more details, please see `gckdir generate --help`.

### Compare

```bash
$ gckdir compare hash_name.json path/to/target_directory
```
or
```bash
$ gckdir compare path/to/source_directory path/to/target_directory
```

For more details, please see `gckdir compare --help`.

**Appendix**

```bash
$ gckdir compare hash_name.json path/to/target_directory --report
```

You can create a comparison result report with `--report` option.

![comparison_report](./images/comparison_report.png)

If you want to check the result immediately on the browser, please use the `--open` option with `--report` option.

```bash
$ gckdir compare hash_name.json path/to/target_directory --report --open
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/KemoKemo/gckdir
```

## Contribution

1. Fork ([https://github.com/KemoKemo/gckdir/fork](https://github.com/KemoKemo/gckdir/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[KemoKemo](https://github.com/KemoKemo)

## License

MIT

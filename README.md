# xpcolor

xpcolor is a `crossplane` wrapper used to add colors to your Crossplane CLI output.

## What's this?

xpcolor is a wrapper that colorizes your `crossplane` command output and does nothing else.
It internally calls the `crossplane` command and colorizes the output, so you can
alias xpcolor as a substitute for `crossplane`:

```sh
alias crossplane="xpcolor"
```

## Features

- Does not tamper with the output. It only adds colors to the `crossplane` output.
- Dynamic TTY detection to send plaintext when called programmatically.
- Colorized `resource trace`, `resource validate`, `cluster top`, rendered YAML, help and errors.

## Installation

```sh
go install github.com/tolgayakar/xpcolor@latest
```

## Usage

Use `xpcolor` anywhere you would use `crossplane`:

```sh
xpcolor resource trace mydatabase.example.org/my-db
xpcolor composition render xr.yaml composition.yaml functions.yaml
```

## Acknowledgements

xpcolor is heavily inspired by [kubecolor](https://github.com/kubecolor/kubecolor),
which does the same for `kubectl`.

## License

[Apache License 2.0](LICENSE)

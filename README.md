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
- Colorized `resource trace` (default and wide), `resource validate`, `cluster top`,
  YAML and JSON output (`composition render`, `-o yaml`, `-o json`, ...), help and errors.
- Dark and light themes, every color configurable with environment variables.

## Installation

```sh
go install github.com/tolgayakar/xpcolor/cmd/xpcolor@latest
```

## Usage

Use `xpcolor` anywhere you would use `crossplane`:

```sh
xpcolor resource trace mydatabase.example.org/my-db
xpcolor composition render xr.yaml composition.yaml functions.yaml
```

Aliases work too, including ones built on top of `crossplane`:

```sh
alias crossplane="xpcolor"
alias crt="crossplane resource trace"
```

`resource trace --watch` and commands that don't print one of the formats above are
passed through untouched.

## Configuration

| Flag | Environment variable | Description |
|---|---|---|
| `--force-colors` | `XPCOLOR_FORCE_COLORS` | Colorize even when the output is not a terminal. |
| `--plain` | `NO_COLOR` | Never colorize. |
| `--xpcolor-theme=dark\|light` | `XPCOLOR_PRESET` | Theme preset, `dark` by default. |
| `--light-background` | `XPCOLOR_LIGHT_BACKGROUND` | Same as `--xpcolor-theme=light`. |
| `--xpcolor-version` | | Print the xpcolor version. |
| | `XPCOLOR_CROSSPLANE` | Path to the `crossplane` binary. |

Arguments after `--` are passed to `crossplane` as they are.

### Colors

Every theme color can be overridden with `XPCOLOR_THEME_<GROUP>_<NAME>`:

```sh
export XPCOLOR_THEME_STATUS_FALSE="fg=white:bg=red"
export XPCOLOR_THEME_DATA_NULL=none
```

A color is a `:` separated list of names (`red`, `hicyan`, `gray`), attributes
(`bold`, `dim`, `italic`, `underline`), 256-color codes (`208`), hex colors (`#c2270a`)
and backgrounds (`bg=red`). `none` disables a color.

| Key | Default (dark) | Used for |
|---|---|---|
| `base.primary` | `magenta` | |
| `base.secondary` | `cyan` | flags, versions, `cluster top` types |
| `base.success` | `green` | |
| `base.warning` | `yellow` | |
| `base.danger` | `red` | |
| `base.info` | `white` | |
| `base.muted` | `gray:italic` | separators, comments |
| `base.key` | `hicyan` | |
| `data.key` | `base.key` | YAML and JSON keys |
| `data.string` | `hiyellow` | |
| `data.number` | `base.primary` | |
| `data.true` / `data.false` | `base.success` / `base.danger` | |
| `data.null` | `base.muted` | |
| `status.true` / `status.false` / `status.unknown` | `base.success` / `base.danger` / `base.warning` | conditions, validate results |
| `table.header` | `bold` | table headers, help sections |
| `trace.tree` | `base.muted` | `├─ └─` |
| `trace.kind` | `base.key` | resource kinds |
| `trace.namespace` | `base.muted` | `(namespace)` |

Keys that default to another key follow it, so changing `base.success` also changes
`data.true` and `status.true` unless they are set themselves.

## Acknowledgements

xpcolor is heavily inspired by [kubecolor](https://github.com/kubecolor/kubecolor),
which does the same for `kubectl`.

## License

[Apache License 2.0](LICENSE)

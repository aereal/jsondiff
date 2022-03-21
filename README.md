[![status][ci-status-badge]][ci-status]
[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]

# jsondiff

jsnodiff provides functions to calculate JSON objects differences with [gojq][] filter.

## Synopsis

See examples on [pkg.go.dev][pkg-go-dev].

## Installation

```sh
go get github.com/aereal/jsondiff
```

## CLI

```sh
go install github.com/aereal/jsondiff/cmd/jsondiff@latest
jsondiff -only '.d' ./testdata/from.json ./testdata/to.json
# --- from.json
# +++ to.json
# @@ -1,2 +1,2 @@
# -4
# +3
```

## License

See LICENSE file.

[pkg-go-dev]: https://pkg.go.dev/github.com/aereal/jsondiff
[pkg-go-dev-badge]: https://pkg.go.dev/badge/aereal/jsondiff
[ci-status-badge]: https://github.com/aereal/jsondiff/workflows/CI/badge.svg?branch=main
[ci-status]: https://github.com/aereal/jsondiff/actions/workflows/CI
[gojq]: https://github.com/itchyny/gojq

# Go bindings for Roc

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/roc-project/roc-go/roc) [![Build Status](https://travis-ci.org/roc-project/roc-go.svg?branch=master)](https://travis-ci.org/roc-project/roc-go) [![Coverage Status](https://coveralls.io/repos/github/roc-project/roc-go/badge.svg?branch=master)](https://coveralls.io/github/roc-project/roc-go?branch=master)

_Work in progress!_

## Dependencies

You will need to have libroc and libroc-devel (headers) installed. Refer to official build [instructions](https://roc-project.github.io/roc/docs/building.html) on how to install libroc. There is no official distribution for any OS as of now, you will need to install from source.

## Installation

```
go get github.com/roc-project/roc-go/roc
```

## Development

Check for compilation and linter errors:

```
make check
```

Run tests:

```
make test
make race # run tests under race detector
```

Format code:

```
make fmt
```

## Authors

See [here](https://github.com/roc-project/roc-go/graphs/contributors).

## License

[MIT](LICENSE)

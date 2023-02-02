# Go bindings for Roc Toolkit

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/roc-streaming/roc-go/roc) [![Build](https://github.com/roc-streaming/roc-go/workflows/build/badge.svg)](https://github.com/roc-streaming/roc-go/actions) [![Coverage Status](https://coveralls.io/repos/github/roc-streaming/roc-go/badge.svg?branch=main)](https://coveralls.io/github/roc-streaming/roc-go?branch=main) [![GitHub release](https://img.shields.io/github/release/roc-streaming/roc-go.svg)](https://github.com/roc-streaming/roc-go/releases) [![Matrix chat](https://matrix.to/img/matrix-badge.svg)](https://app.element.io/#/room/#roc-streaming:matrix.org)

This library provides Go (golang) bindings for [Roc Toolkit](https://github.com/roc-streaming/roc-toolkit), a toolkit for real-time audio streaming over the network.

## About Roc

Compatible senders and receivers include:

* [command-line tools](https://roc-streaming.org/toolkit/docs/tools/command_line_tools.html)
* [sound server modules](https://roc-streaming.org/toolkit/docs/tools/sound_server_modules.html) (PulseAudio, PipeWire)
* [C library](https://roc-streaming.org/toolkit/docs/api.html)
* [Java bindings](https://github.com/roc-streaming/roc-java/) and [Android app](https://github.com/roc-streaming/roc-droid) that uses them

Key features:

* real-time streaming with guaranteed latency;
* restoring lost packets using Forward Erasure Correction codes;
* converting between the sender and receiver clock domains;
* CD-quality audio;
* multiple profiles for different CPU and latency requirements;
* portability;
* relying on open, standard protocols.

## Documentation

Documentation for the bindings is availabe on [pkg.go.dev](https://pkg.go.dev/github.com/roc-streaming/roc-go/roc).

Documentation for the underlying C API can be found [here](https://roc-streaming.org/toolkit/docs/api.html).

## Versioning

Go bindings and the C library both use [semantic versioning](https://semver.org/).

Rules prior to 1.0.0 release:

* According to semantic versioning, there is no compatibility promise until 1.0.0 is released. Small breaking changes are possible. For convenience, breaking changes are introduced only in minor version updates, but not in patch version updates.

Rules starting from 1.0.0 release:

* The first two components (major and minor) of the bindings and the C library versions correspond to each other. The third component (patch) is indepdendent.

  **Bindings are compatible with the C library if its major version is the same, and minor version is the same or higher.**

  For example, version 1.2.3 of the bindings would be compatible with 1.2.x and 1.3.x, but not with 1.1.x (minor version is lower) or 2.x.x (major version is different).

## Installation

You will need to have Roc Toolkit library and headers installed system-wide. Refer to official build [instructions](https://roc-streaming.org/toolkit/docs/building.html) on how to install libroc from source.

After installing libroc, you can install bindings using regular `go get`:

```
go get github.com/roc-streaming/roc-go/roc
```

## Development

Run all checks:

```
make
```

Only run specific checks:

```
make build
make lint
make test
```

Update modules:

```
make tidy
```

Format code:

```
make fmt
```

## Authors

See [here](https://github.com/roc-streaming/roc-go/graphs/contributors).

## License

Bindings are licensed under [MIT](LICENSE).

For details on Roc Toolkit licensing, see [here](https://roc-streaming.org/toolkit/docs/about_project/licensing.html).

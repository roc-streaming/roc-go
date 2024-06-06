# Go bindings for Roc Toolkit

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/roc-streaming/roc-go/roc) [![Build](https://github.com/roc-streaming/roc-go/workflows/build/badge.svg)](https://github.com/roc-streaming/roc-go/actions) [![Coverage Status](https://coveralls.io/repos/github/roc-streaming/roc-go/badge.svg?branch=main)](https://coveralls.io/github/roc-streaming/roc-go?branch=main) [![GitHub release](https://img.shields.io/github/release/roc-streaming/roc-go.svg)](https://github.com/roc-streaming/roc-go/releases) [![Matrix chat](https://matrix.to/img/matrix-badge.svg)](https://app.element.io/#/room/#roc-streaming:matrix.org)

This library provides Go (golang) bindings for [Roc Toolkit](https://github.com/roc-streaming/roc-toolkit), a toolkit for real-time audio streaming over the network.

## About Roc

Key features of Roc Toolkit:

* real-time streaming with guaranteed latency;
* robust work on unreliable networks like Wi-Fi, due to use of Forward Erasure Correction codes;
* CD-quality audio;
* multiple profiles for different CPU and latency requirements;
* relying on open, standard protocols, like RTP and FECFRAME;
* interoperability with both Roc and third-party software.

Compatible Roc Toolkit senders and receivers include:

* [cross-platform command-line tools](https://roc-streaming.org/toolkit/docs/tools/command_line_tools.html)
* [modules for sound servers](https://roc-streaming.org/toolkit/docs/tools/sound_server_modules.html) (PulseAudio, PipeWire)
* [C library](https://roc-streaming.org/toolkit/docs/api.html) and [bindings for other languages](https://roc-streaming.org/toolkit/docs/api/bindings.html)
* [end-user apps](https://roc-streaming.org/toolkit/docs/tools/applications.html)

## Documentation

Documentation for the bindings is available on [pkg.go.dev](https://pkg.go.dev/github.com/roc-streaming/roc-go/roc).

Documentation for the underlying C API can be found [here](https://roc-streaming.org/toolkit/docs/api.html).

## Quick start

#### Sender

```go
import (
	"github.com/roc-streaming/roc-go/roc"
)

context, err := roc.OpenContext(roc.ContextConfig{})
if err != nil {
	panic(err)
}
defer context.Close()

sender, err := roc.OpenSender(roc.SenderConfig{
	FrameEncoding: roc.MediaEncoding{
		Rate:     44100,
		Format:   roc.FormatPcmFloat32,
		Channels: roc.ChannelLayoutStereo,
	},
	FecEncoding: roc.FecEncodingRs8m,
	ClockSource: roc.ClockSourceInternal,
})
if err != nil {
	panic(err)
}
defer sender.Close()

sourceEndpoint, err := roc.ParseEndpoint("rtp+rs8m://192.168.0.1:10001")
if err != nil {
	panic(err)
}

repairEndpoint, err := roc.ParseEndpoint("rs8m://192.168.0.1:10002")
if err != nil {
	panic(err)
}

controlEndpoint, err := roc.ParseEndpoint("rtcp://192.168.0.1:10003")
if err != nil {
	panic(err)
}

err = sender.Connect(roc.SlotDefault, roc.InterfaceAudioSource, sourceEndpoint)
if err != nil {
	panic(err)
}

err = sender.Connect(roc.SlotDefault, roc.InterfaceAudioRepair, repairEndpoint)
if err != nil {
	panic(err)
}

err = sender.Connect(roc.SlotDefault, roc.InterfaceAudioControl, controlEndpoint)
if err != nil {
	panic(err)
}

for {
	samples := make([]float32, 320)

	/* fill samples */

	err = sender.WriteFloats(samples)
	if err != nil {
		panic(err)
	}
}
```

#### Receiver

```go
import (
	"github.com/roc-streaming/roc-go/roc"
)

context, err := roc.OpenContext(roc.ContextConfig{})
if err != nil {
	panic(err)
}
defer context.Close()

receiver, err := roc.OpenReceiver(roc.ReceiverConfig{
	FrameEncoding: roc.MediaEncoding{
		Rate:     44100,
		Format:   roc.FormatPcmFloat32,
		Channels: roc.ChannelLayoutStereo,
	},
	ClockSource: roc.ClockSourceInternal,
})
if err != nil {
	panic(err)
}
defer receiver.Close()

sourceEndpoint, err := roc.ParseEndpoint("rtp+rs8m://0.0.0.0:10001")
if err != nil {
	panic(err)
}

repairEndpoint, err := roc.ParseEndpoint("rs8m://0.0.0.0:10002")
if err != nil {
	panic(err)
}

controlEndpoint, err := roc.ParseEndpoint("rtcp://0.0.0.0:10003")
if err != nil {
	panic(err)
}

err = receiver.Bind(roc.SlotDefault, roc.InterfaceAudioSource, sourceEndpoint)
if err != nil {
	panic(err)
}

err = receiver.Bind(roc.SlotDefault, roc.InterfaceAudioRepair, repairEndpoint)
if err != nil {
	panic(err)
}

err = receiver.Bind(roc.SlotDefault, roc.InterfaceAudioControl, controlEndpoint)
if err != nil {
	panic(err)
}

for {
	samples := make([]float32, 320)

	err = receiver.ReadFloats(samples)
	if err != nil {
		panic(err)
	}

	/* process samples */
}
```

## Versioning

Go bindings and the C library both use [semantic versioning](https://semver.org/).

Bindings are **compatible** with the C library when:

- **major** version of bindings **is same** as major version of C library
- **minor** version of bindings **is same or higher** as minor version of C library

Patch versions of bindings and C library are independent.

For example, version 1.2.3 of the bindings would be compatible with 1.2.x and 1.3.x, but not with 1.1.x (minor version is lower) or 2.x.x (major version is different).

## Installation

You will need to have Roc Toolkit library and headers installed system-wide. Refer to official build [instructions](https://roc-streaming.org/toolkit/docs/building.html) on how to install it.

After installing libroc, you can install bindings using regular `go get`:

```
go get github.com/roc-streaming/roc-go/roc
```

## Development

Install development dependencies:

* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

* [stringer](https://github.com/golang/tools)

    `go install golang.org/x/tools/cmd/stringer@latest`

Run all checks:

```
make
```

Run only specific checks:

```
make gen|build|lint|test|testall
```

Update modules:

```
make tidy
```

Format code:

```
make fmt
```

## Release

To release a new version:

 * Create git tag

    ```
    ./tag.py --push <remote> <version>
    ```

    e.g.

    ```
    ./tag.py --push origin 1.2.3
    ```

    Or use **tag.py** without **--push** to only create a tag locally, and then push it manually.

* Wait until "Release" CI job completes and creates GitHub release draft.

* Edit GitHub release created by CI and publish it.

## Authors

See [here](https://github.com/roc-streaming/roc-go/graphs/contributors).

## License

Bindings are licensed under [MIT](LICENSE).

For details on Roc Toolkit licensing, see [here](https://roc-streaming.org/toolkit/docs/about_project/licensing.html).

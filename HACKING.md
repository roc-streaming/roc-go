# Hacking guide

## Toolkit version

Make sure you've installed matching Roc Toolkit version. If you've installed it via package manager, it may be outdated, and you may need to install from sources.

Versioning policy and installation instructions are documented in README.

## Go dependencies

Additional dependencies needed for development:

* [golangci-lint](https://golangci-lint.run/welcome/install/#local-installation)

## Makefile targets

Update and build everything:

```
make
```

Run specific step:

```
make gen|build|lint|test|test_all
```

Update go modules:

```
make tidy
```

Format code:

```
make fmt
```

## Making release

To release a new version:

 * Create git tag

    ```
    ./make_tag.py --push <remote> <version>
    ```

    e.g.

    ```
    ./make_tag.py --push origin 1.2.3
    ```

    You can omit `--push origin` to only create a tag locally without pushing it to github.

* Wait until "Release" CI job completes and creates GitHub release draft.

* Edit GitHub release created by CI and publish it.

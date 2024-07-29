# leakless

Golang doesn't provide a way for third-party packages to implicitly register a callback when the main process exits to do cleanups. If you fail to explicitly add proper code in your main to kill sub-processes that are created by third-party packages they may keep running after the main process exists or crashes.

Leakless ensures to kill the sub-process when the main process exits.

How it works:

1. Main process outputs a standalone executable file and executes it as a guard process.
1. A TCP connection is created between main process and the guard process.
1. The guard process starts the sub-process.
1. If the TCP connection is closed, the guard process will kill the sub-process.

This design ensures it works consistently across different platforms, the CI tests Linux, Mac, and Windows.

If you don't trust the executable, you can build it yourself from the source code by running `go generate` at the root of this repo, then use the [replace](https://golang.org/ref/mod#go-mod-file-replace) to use your own module. Usually, it won't be a concern, all the executables are committed by this [Github Action](https://github.com/ysmood/leakless/actions?query=workflow%3ARelease), the Action will print the hash of the commit, you can compare it with the repo.

Not using the PID is because after a process exits, a newly created process may have the same PID.

## How to Use

See the [examples](example_test.go).

## Custom build for `GOOS` or `GOARCH`

Such as if you want to support FreeBSD, you can clone this project and modify the [targets.go](cmd/pack/targets.go) to something like:

```go
var targets = []utils.Target{
    "freebsd/amd64",
}
```

Then run `go generate` and use [replace](https://golang.org/ref/mod#go-mod-file-replace) in the project that will use leakless.
You can keep this fork of leakless to serve your own interest.

# leakless

Run sub-process and make sure to kill it when the parent process exits.
The way how it works is to check parent tcp connection with an uuid.

Not using the pid is because after a process exits, newly created process may have the same pid.

# How to Use

See the [examples](example_test.go)

## Deploy

```
godev build -d ./cmd/leakless
```

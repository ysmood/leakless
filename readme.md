# leakless

Run sub-process and make sure to kill it when the parent process exits.
The way how it works is to check parent tcp connection with a uuid.

Not using the pid is because after a process exits, newly created process may have the same pid.

## Deploy

```
godev build -d ./cmd/leakless
```

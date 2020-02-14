# leakless

Run sub-process and make sure to kill it when the parent process exit.
The way how it works is to check parent tcp connection.

## Deploy

```
godev build -d ./cmd/leakless
```

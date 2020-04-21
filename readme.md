# leakless

Run sub-process and make sure to kill it when the parent process exits.
The way how it works is to output a standalone executable file to guard the subprocess and check parent TCP connection with a UUID.
So that it works consistently on Linux, Mac, and Windows.

Not using the PID is because after a process exits, a newly created process may have the same PID.

## How to Use

See the [examples](example_test.go).

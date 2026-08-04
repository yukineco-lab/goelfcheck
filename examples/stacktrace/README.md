# Stack Trace Trimpath Demo

This sample prints `runtime.Caller()` frames and then intentionally panics so
you can compare file paths with and without `-trimpath`.

## Build and Run

```sh
go build -o stacktrace-default .
./stacktrace-default
```

The stack trace normally includes the local checkout path:

```text
/home/user/goelfcheck/examples/stacktrace/main.go
/home/user/goelfcheck/examples/stacktrace/main.go:15
```

Build with `-trimpath`:

```sh
go build -trimpath -o stacktrace-trimpath .
./stacktrace-trimpath
```

The stack trace uses the module/package path instead:

```text
goelfcheck/examples/stacktrace/main.go
goelfcheck/examples/stacktrace/main.go:15
```

The same difference appears in `runtime.Caller()` output before the panic:

```text
skip=0 func=main.printCallerFrames file=/home/user/goelfcheck/examples/stacktrace/main.go line=...
```

With `-trimpath`:

```text
skip=0 func=main.printCallerFrames file=goelfcheck/examples/stacktrace/main.go line=...
```

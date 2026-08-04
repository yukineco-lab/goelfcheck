# Stack Trace Trimpath Demo

This sample intentionally panics so you can compare stack traces with and without
`-trimpath`.

## Build and Run

```sh
go build -o stacktrace-default .
./stacktrace-default
```

The stack trace normally includes the local checkout path:

```text
/home/user/goelfcheck/examples/stacktrace/main.go:15
```

Build with `-trimpath`:

```sh
go build -trimpath -o stacktrace-trimpath .
./stacktrace-trimpath
```

The stack trace uses the module/package path instead:

```text
goelfcheck/examples/stacktrace/main.go:15
```

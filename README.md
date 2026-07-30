# goelfcheck

`goelfcheck` is a small Go-only CLI for inspecting Linux ELF files built from Go
programs. It focuses on information-leak checks instead of acting as a full
`readelf` replacement.

## Usage

```sh
go run ./cmd/goelfcheck <ELF file>
go run ./cmd/goelfcheck -v <ELF file>
```

The command exits with status `1` when it finds Build ID, VCS metadata, DWARF
debug information, symbol tables, or absolute path-like strings.

## Recommended Build

```sh
go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w -buildid="
```

## Checks

- ELF header, program header, and section header summary
- Go version and `debug/buildinfo` settings
- Go and GNU Build ID notes
- Go-specific sections such as `.gopclntab` and `.go.buildinfo`
- VCS settings such as `vcs.revision`, `vcs.time`, and `vcs.modified`
- DWARF compile units, source files, directories, and producers
- Static and dynamic symbol tables
- Absolute path-like strings such as `/home/...`, `/Users/...`, `C:\Users\...`,
  and `D:\...`

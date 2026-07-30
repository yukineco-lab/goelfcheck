package report

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"goelfcheck/internal/inspect"
)

func Print(w io.Writer, r inspect.Result, verbose bool) {
	fmt.Fprintln(w, "====================================================")
	fmt.Fprintln(w, "Go Binary Inspection Tool")
	fmt.Fprintln(w, "====================================================")
	fmt.Fprintf(w, "Target         : %s\n\n", r.Path)

	printSummary(w, r)
	printChecks(w, r)
	printOverall(w, r)
	printRecommendedBuild(w)

	if verbose {
		printVerbose(w, r)
	}

	if len(r.Diagnostics) > 0 {
		fmt.Fprintln(w, "\n[Diagnostics]")
		for _, d := range r.Diagnostics {
			fmt.Fprintf(w, "  - %s\n", d)
		}
	}
}

func printSummary(w io.Writer, r inspect.Result) {
	fmt.Fprintln(w, "[ELF]")
	fmt.Fprintf(w, "  Class          : %s\n", r.ELF.Class)
	fmt.Fprintf(w, "  Architecture   : %s\n", r.ELF.Machine)
	fmt.Fprintf(w, "  Endian         : %s\n", endianLabel(r.ELF.Data))
	fmt.Fprintf(w, "  Type           : %s\n", r.ELF.Type)
	fmt.Fprintf(w, "  PIE            : %s\n", present(r.ELF.IsPIE))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "[Go Build]")
	fmt.Fprintf(w, "  Go Version     : %s\n", valueOrNA(r.Build.GoVersion))
	fmt.Fprintf(w, "  GOOS           : %s\n", buildSetting(r.Build.Settings, "GOOS"))
	fmt.Fprintf(w, "  GOARCH         : %s\n", buildSetting(r.Build.Settings, "GOARCH"))
	fmt.Fprintf(w, "  CGO_ENABLED    : %s\n", buildSetting(r.Build.Settings, "CGO_ENABLED"))
	fmt.Fprintf(w, "  Build ID       : %s\n", present(r.Build.HasBuildID()))
	fmt.Fprintf(w, "  VCS Info       : %s\n", present(r.Build.HasVCSInfo()))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "[Debug]")
	fmt.Fprintf(w, "  DWARF          : %s\n", present(r.DWARF.Present))
	fmt.Fprintf(w, "  Symbols        : %s\n", present(r.Symbols.HasSymbols()))
	fmt.Fprintf(w, "  Absolute Paths : %s\n", present(len(r.Strings.AbsolutePaths) > 0))
	fmt.Fprintln(w)
}

func printChecks(w io.Writer, r inspect.Result) {
	fmt.Fprintln(w, "[Security Check]")
	fmt.Fprintln(w)

	if len(r.Strings.AbsolutePaths) == 0 {
		fmt.Fprintln(w, "[OK] No absolute source path found")
	} else {
		fmt.Fprintln(w, "[NG] Absolute path-like strings detected")
		fmt.Fprintln(w, "     -> Use -trimpath and avoid embedding local file paths")
		for _, p := range firstN(r.Strings.AbsolutePaths, 10) {
			fmt.Fprintf(w, "        %s\n", p)
		}
		if len(r.Strings.AbsolutePaths) > 10 {
			fmt.Fprintf(w, "        ... and %d more\n", len(r.Strings.AbsolutePaths)-10)
		}
	}
	fmt.Fprintln(w)

	if r.Build.HasBuildID() {
		fmt.Fprintln(w, "[NG] Build ID detected")
		fmt.Fprintln(w, "     -> Use -ldflags=\"-buildid=\"")
	} else {
		fmt.Fprintln(w, "[OK] Build ID not detected")
	}
	fmt.Fprintln(w)

	if r.Build.HasVCSInfo() {
		fmt.Fprintln(w, "[NG] VCS information detected")
		fmt.Fprintln(w, "     -> Use -buildvcs=false")
	} else {
		fmt.Fprintln(w, "[OK] VCS information not detected")
	}
	fmt.Fprintln(w)

	if r.DWARF.Present {
		fmt.Fprintln(w, "[NG] DWARF debug information detected")
		fmt.Fprintln(w, "     -> Use -ldflags=\"-w\"")
	} else {
		fmt.Fprintln(w, "[OK] DWARF debug information not detected")
	}
	fmt.Fprintln(w)

	if r.Symbols.HasSymbols() {
		fmt.Fprintln(w, "[NG] Symbol table detected")
		fmt.Fprintln(w, "     -> Use -ldflags=\"-s\"")
	} else {
		fmt.Fprintln(w, "[OK] Symbol table not detected")
	}
}

func printOverall(w io.Writer, r inspect.Result) {
	risk := "LOW"
	switch {
	case len(r.Strings.AbsolutePaths) > 0 || r.DWARF.Present || r.Build.HasVCSInfo():
		risk = "HIGH"
	case r.Build.HasBuildID() || r.Symbols.HasSymbols():
		risk = "MEDIUM"
	}

	result := "OK"
	if r.HasFindings() {
		result = "NG"
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Overall")
	fmt.Fprintln(w, "--------")
	fmt.Fprintf(w, "Result : %s\n", result)
	fmt.Fprintf(w, "Risk   : %s\n", risk)
}

func printRecommendedBuild(w io.Writer) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Recommended build")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "go build \\")
	fmt.Fprintln(w, "    -trimpath \\")
	fmt.Fprintln(w, "    -buildvcs=false \\")
	fmt.Fprintln(w, "    -ldflags=\"-s -w -buildid=\"")
}

func printVerbose(w io.Writer, r inspect.Result) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: ELF Header]")
	fmt.Fprintf(w, "  OS ABI         : %s\n", r.ELF.OSABI)
	fmt.Fprintf(w, "  Entry          : 0x%x\n", r.ELF.Entry)
	fmt.Fprintf(w, "  GNU_STACK      : %s\n", present(r.ELF.HasGNUStack))
	fmt.Fprintf(w, "  RELRO          : %s\n", present(r.ELF.HasRelRO))

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: Program Headers]")
	for _, p := range r.ELF.Programs {
		fmt.Fprintf(w, "  %-14s flags=%-8s off=0x%x vaddr=0x%x filesz=%d memsz=%d align=%d\n",
			p.Type, p.Flags, p.Off, p.Vaddr, p.Filesz, p.Memsz, p.Align)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: Section Headers]")
	for _, s := range r.ELF.Sections {
		fmt.Fprintf(w, "  %-24s type=%-16s flags=%-8s off=0x%x addr=0x%x size=%d\n",
			s.Name, s.Type, s.Flags, s.Off, s.Addr, s.Size)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: Go Sections]")
	if len(r.ELF.GoSections) == 0 {
		fmt.Fprintln(w, "  n/a")
	} else {
		for _, s := range r.ELF.GoSections {
			fmt.Fprintf(w, "  %-24s size=%d\n", s.Name, s.Size)
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: BuildInfo]")
	fmt.Fprintf(w, "  Go Build ID    : %s\n", valueOrNA(shorten(r.Build.GoBuildID, 80)))
	fmt.Fprintf(w, "  GNU Build ID   : %s\n", valueOrNA(shorten(r.Build.GNUBuildID, 80)))
	if r.Build.ReadError != "" && r.Build.GoVersion == "" {
		fmt.Fprintf(w, "  Read Error     : %s\n", r.Build.ReadError)
	}
	for _, s := range r.Build.Settings {
		fmt.Fprintf(w, "  %-14s : %s\n", s.Key, s.Value)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: DWARF Compile Units]")
	if !r.DWARF.Present {
		fmt.Fprintln(w, "  n/a")
	} else {
		for _, cu := range firstUnits(r.DWARF.CompileUnits, 20) {
			fmt.Fprintf(w, "  Unit           : %s\n", valueOrNA(cu.Name))
			fmt.Fprintf(w, "    Directory    : %s\n", valueOrNA(cu.Directory))
			fmt.Fprintf(w, "    Producer     : %s\n", valueOrNA(cu.Producer))
			for _, file := range firstN(cu.SourceFiles, 8) {
				fmt.Fprintf(w, "    Source       : %s\n", filepath.ToSlash(file))
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: Symbols]")
	fmt.Fprintf(w, "  .symtab        : %d\n", r.Symbols.StaticCount)
	for _, s := range r.Symbols.StaticSample {
		fmt.Fprintf(w, "    %s\n", s)
	}
	fmt.Fprintf(w, "  .dynsym        : %d\n", r.Symbols.DynamicCount)
	for _, s := range r.Symbols.DynamicSample {
		fmt.Fprintf(w, "    %s\n", s)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "[Verbose: Absolute Paths]")
	if len(r.Strings.AbsolutePaths) == 0 {
		fmt.Fprintln(w, "  n/a")
	} else {
		for _, p := range r.Strings.AbsolutePaths {
			fmt.Fprintf(w, "  %s\n", p)
		}
	}
}

func endianLabel(data fmt.Stringer) string {
	switch data.String() {
	case "ELFDATA2LSB":
		return "Little Endian"
	case "ELFDATA2MSB":
		return "Big Endian"
	default:
		return data.String()
	}
}

func present(ok bool) string {
	if ok {
		return "Present"
	}
	return "Not detected"
}

func valueOrNA(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func buildSetting(settings []inspect.BuildSetting, key string) string {
	for _, s := range settings {
		if s.Key == key {
			return s.Value
		}
	}
	return "n/a"
}

func firstN(values []string, n int) []string {
	if len(values) <= n {
		return values
	}
	return values[:n]
}

func firstUnits(values []inspect.CompileUnit, n int) []inspect.CompileUnit {
	if len(values) <= n {
		return values
	}
	return values[:n]
}

func shorten(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

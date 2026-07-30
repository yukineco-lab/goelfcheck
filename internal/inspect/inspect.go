package inspect

import (
	"debug/elf"
	"fmt"
)

func File(path string) (Result, error) {
	f, err := elf.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("open ELF: %w", err)
	}
	defer f.Close()

	result := Result{Path: path}
	result.ELF = inspectELF(f)
	result.Build = inspectBuildInfo(path, f)
	result.DWARF = inspectDWARF(f)
	result.Symbols = inspectSymbols(f)

	scan, err := scanStrings(path)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, fmt.Sprintf("string scan failed: %v", err))
	} else {
		result.Strings = scan
	}

	return result, nil
}

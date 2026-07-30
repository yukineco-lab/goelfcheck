package inspect

import "debug/elf"

const symbolSampleLimit = 20

func inspectSymbols(f *elf.File) SymbolInfo {
	var out SymbolInfo

	if symbols, err := f.Symbols(); err == nil {
		out.StaticCount = len(symbols)
		out.StaticSample = sampleSymbols(symbols)
	}
	if symbols, err := f.DynamicSymbols(); err == nil {
		out.DynamicCount = len(symbols)
		out.DynamicSample = sampleSymbols(symbols)
	}

	return out
}

func sampleSymbols(symbols []elf.Symbol) []string {
	var sample []string
	for _, sym := range symbols {
		if sym.Name == "" {
			continue
		}
		sample = append(sample, sym.Name)
		if len(sample) >= symbolSampleLimit {
			break
		}
	}
	return sample
}

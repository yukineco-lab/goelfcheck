package inspect

import "debug/elf"

type Result struct {
	Path        string
	ELF         ELFInfo
	Build       BuildInfo
	DWARF       DWARFInfo
	Symbols     SymbolInfo
	Strings     StringScan
	Diagnostics []string
}

func (r Result) HasFindings() bool {
	return r.Build.HasBuildID() ||
		r.Build.HasVCSInfo() ||
		r.DWARF.Present ||
		r.Symbols.HasSymbols() ||
		len(r.Strings.AbsolutePaths) > 0
}

type ELFInfo struct {
	Class        elf.Class
	Data         elf.Data
	Type         elf.Type
	Machine      elf.Machine
	OSABI        elf.OSABI
	Entry        uint64
	Sections     []SectionInfo
	Programs     []ProgramInfo
	GoSections   []SectionInfo
	HasGNUStack  bool
	HasRelRO     bool
	IsExecutable bool
	IsPIE        bool
}

type SectionInfo struct {
	Name  string
	Type  elf.SectionType
	Flags elf.SectionFlag
	Addr  uint64
	Off   uint64
	Size  uint64
}

type ProgramInfo struct {
	Type   elf.ProgType
	Flags  elf.ProgFlag
	Vaddr  uint64
	Off    uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

type BuildInfo struct {
	GoVersion  string
	Settings   []BuildSetting
	GoBuildID  string
	GNUBuildID string
	ReadError  string
}

func (b BuildInfo) HasBuildID() bool {
	return b.GoBuildID != "" || b.GNUBuildID != ""
}

func (b BuildInfo) HasVCSInfo() bool {
	for _, s := range b.Settings {
		switch s.Key {
		case "vcs", "vcs.revision", "vcs.time", "vcs.modified":
			if s.Value != "" {
				return true
			}
		}
	}
	return false
}

type BuildSetting struct {
	Key   string
	Value string
}

type DWARFInfo struct {
	Present      bool
	CompileUnits []CompileUnit
	ReadError    string
}

type CompileUnit struct {
	Name        string
	Directory   string
	Producer    string
	SourceFiles []string
}

type SymbolInfo struct {
	StaticCount   int
	DynamicCount  int
	StaticSample  []string
	DynamicSample []string
}

func (s SymbolInfo) HasSymbols() bool {
	return s.StaticCount > 0 || s.DynamicCount > 0
}

type StringScan struct {
	AbsolutePaths []string
}

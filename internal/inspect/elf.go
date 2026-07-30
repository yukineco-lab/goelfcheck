package inspect

import (
	"debug/elf"
	"strings"
)

func inspectELF(f *elf.File) ELFInfo {
	info := ELFInfo{
		Class:        f.Class,
		Data:         f.Data,
		Type:         f.Type,
		Machine:      f.Machine,
		OSABI:        f.OSABI,
		Entry:        f.Entry,
		IsExecutable: f.Type == elf.ET_EXEC,
		IsPIE:        f.Type == elf.ET_DYN,
	}

	for _, s := range f.Sections {
		section := SectionInfo{
			Name:  s.Name,
			Type:  s.Type,
			Flags: s.Flags,
			Addr:  s.Addr,
			Off:   s.Offset,
			Size:  s.Size,
		}
		info.Sections = append(info.Sections, section)
		if isGoSection(s.Name) {
			info.GoSections = append(info.GoSections, section)
		}
		if s.Name == ".got" || s.Name == ".got.plt" {
			info.HasRelRO = true
		}
	}

	for _, p := range f.Progs {
		info.Programs = append(info.Programs, ProgramInfo{
			Type:   p.Type,
			Flags:  p.Flags,
			Vaddr:  p.Vaddr,
			Off:    p.Off,
			Filesz: p.Filesz,
			Memsz:  p.Memsz,
			Align:  p.Align,
		})
		if p.Type == elf.PT_GNU_STACK {
			info.HasGNUStack = true
		}
		if p.Type == elf.PT_GNU_RELRO {
			info.HasRelRO = true
		}
	}

	return info
}

func isGoSection(name string) bool {
	return name == ".gopclntab" ||
		name == ".go.buildinfo" ||
		strings.HasPrefix(name, ".go.") ||
		name == ".typelink" ||
		name == ".itablink" ||
		name == ".noptrdata" ||
		name == ".noptrbss"
}

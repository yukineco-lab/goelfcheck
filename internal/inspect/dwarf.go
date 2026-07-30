package inspect

import (
	"debug/dwarf"
	"debug/elf"
)

func inspectDWARF(f *elf.File) DWARFInfo {
	d, err := f.DWARF()
	if err != nil {
		return DWARFInfo{Present: false, ReadError: err.Error()}
	}

	out := DWARFInfo{Present: true}
	reader := d.Reader()
	for {
		entry, err := reader.Next()
		if err != nil {
			out.ReadError = err.Error()
			return out
		}
		if entry == nil {
			break
		}
		if entry.Tag != dwarf.TagCompileUnit {
			continue
		}

		unit := CompileUnit{
			Name:      attrString(entry, dwarf.AttrName),
			Directory: attrString(entry, dwarf.AttrCompDir),
			Producer:  attrString(entry, dwarf.AttrProducer),
		}
		if lr, err := d.LineReader(entry); err == nil && lr != nil {
			for _, file := range lr.Files() {
				if file != nil && file.Name != "" {
					unit.SourceFiles = append(unit.SourceFiles, file.Name)
				}
			}
		}
		out.CompileUnits = append(out.CompileUnits, unit)
	}

	return out
}

func attrString(entry *dwarf.Entry, attr dwarf.Attr) string {
	value := entry.Val(attr)
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

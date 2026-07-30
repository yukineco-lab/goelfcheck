package inspect

import (
	"bytes"
	"debug/buildinfo"
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func inspectBuildInfo(path string, f *elf.File) BuildInfo {
	var out BuildInfo

	if bi, err := buildinfo.ReadFile(path); err == nil {
		out.GoVersion = bi.GoVersion
		for _, s := range bi.Settings {
			out.Settings = append(out.Settings, BuildSetting{Key: s.Key, Value: s.Value})
		}
	} else {
		out.ReadError = err.Error()
	}

	out.GoBuildID = readGoBuildID(f)
	out.GNUBuildID = readGNUBuildID(f)
	return out
}

func readGoBuildID(f *elf.File) string {
	if s := f.Section(".note.go.buildid"); s != nil {
		if id := readNoteDescString(s, f.Data); id != "" {
			return id
		}
	}
	return ""
}

func readGNUBuildID(f *elf.File) string {
	if s := f.Section(".note.gnu.build-id"); s != nil {
		if data := readNoteDescBytes(s, f.Data, "GNU"); len(data) > 0 {
			return hex.EncodeToString(data)
		}
	}
	return ""
}

func readNoteDescString(s *elf.Section, dataOrder elf.Data) string {
	data := readNoteDescBytes(s, dataOrder, "")
	if len(data) == 0 {
		return ""
	}
	return string(bytes.TrimRight(data, "\x00"))
}

func readNoteDescBytes(s *elf.Section, dataOrder elf.Data, wantedName string) []byte {
	data, err := s.Data()
	if err != nil {
		return nil
	}

	var order binary.ByteOrder = binary.LittleEndian
	if dataOrder == elf.ELFDATA2MSB {
		order = binary.BigEndian
	}

	for off := 0; off+12 <= len(data); {
		namesz := int(order.Uint32(data[off:]))
		descsz := int(order.Uint32(data[off+4:]))
		off += 12
		if off+align4(namesz)+align4(descsz) > len(data) {
			return nil
		}
		name := string(bytes.TrimRight(data[off:off+namesz], "\x00"))
		off += align4(namesz)
		desc := data[off : off+descsz]
		off += align4(descsz)
		if wantedName == "" || wantedName == name {
			return desc
		}
	}
	return nil
}

func align4(n int) int {
	return (n + 3) &^ 3
}

func settingValue(settings []BuildSetting, key string) string {
	for _, s := range settings {
		if s.Key == key {
			return s.Value
		}
	}
	return ""
}

func formatSetting(settings []BuildSetting, key string) string {
	if v := settingValue(settings, key); v != "" {
		return v
	}
	return "n/a"
}

func shortID(id string) string {
	if len(id) <= 18 {
		return id
	}
	return fmt.Sprintf("%s...", id[:18])
}

package main

//go:generate stringer -type FileType,Class,Data,ProgramHeaderFlag,ProgramHeaderType,SectionHeaderFlag,SectionHeaderType -output string.go

type ELFFile struct {
	Header         *Header64
	ProgramHeaders []ProgramHeader64
	SectionHeaders []SectionHeader64
}

var MagicBytes = [4]byte{0x7f, 'E', 'L', 'F'}

// ELFIdentifier is the header which is independent from 32bit/64bit (Class),
// OS (OSABI), and endianess (Data).
type ELFIdentifier struct {
	Magic      [4]byte
	Class      Class // 32bit or 64bit
	Data       Data  // little or big endian
	Version    byte  // always 1
	OSABI      byte  // operating system, linux 0x03
	ABIVersion byte  //
	Padding    [7]byte
}

type Header64 struct {
	ELFIdentifier
	Type    FileType
	Machine uint16 // Machine specifies ISA (e.g. 0x03 for x86)
	Version uint32 // Version is always set to 1
	Entry   uint64 // Entry point for executable files

	ProgramHeaderOffset uint64 // where in the file do the program headers start
	SectionHeaderOffset uint64 // where in the file do the section headers start
	Flags               uint32
	EhSize              uint16
	ProgramHeaderSize   uint16
	ProgramHeaderCount  uint16
	SectionHeaderSize   uint16
	SectionHeaderCount  uint16

	SectionHeaderStringIndex uint16
}

type FileType uint16

const (
	ET_NONE   FileType = 0x00
	ET_REL    FileType = 0x01   // Relocatable file.
	ET_EXEC   FileType = 0x02   // Executable file.
	ET_DYN    FileType = 0x03   // Shared object.
	ET_CORE   FileType = 0x04   // Core file.
	ET_LOOS   FileType = 0xFE00 // Reserved inclusive range. Operating system specific.
	ET_HIOS   FileType = 0xFEFF
	ET_LOPROC FileType = 0xFF00 // Reserved inclusive range. Processor specific.
	ET_HIPROC FileType = 0xFFFF
)

type Class byte

const (
	ELFCLASSNONE Class = iota
	ELFCLASS32
	ELFCLASS64
)

type Data byte

const (
	ELFDATANONE Data = iota
	ELFDATA2LSB      // little-endian
	ELFDATA2MSB      // big-endian
)

type ProgramHeader64 struct {
	Type            ProgramHeaderType
	Flags           ProgramHeaderFlag
	Offset          uint64
	VirtualAddress  uint64
	PhysicalAddress uint64
	FileSize        uint64
	MemorySize      uint64
	Align           uint64
}

type ProgramHeaderFlag uint32

const (
	PF_X ProgramHeaderFlag = 0x1
	PF_W ProgramHeaderFlag = 0x2
	PF_R ProgramHeaderFlag = 0x4
)

type ProgramHeaderType uint32

const (
	PT_NULL ProgramHeaderType = iota
	PT_LOAD
	PT_DYNAMIC
	PT_INTERP
	PT_NOTE
	PT_SHLIB
	PT_PHDR

	PT_LOPROC ProgramHeaderType = 0x70000000
	PT_HIPROC ProgramHeaderType = 0x7fffffff
)

type SectionHeader64 struct {
	// Name specifies the name of the section. Its value is an index into
	// the section header string table section [see "String Table'' below],
	// giving the location of a null-terminated string
	Name uint32

	// Type categorizes the section's contents and semantics.
	Type SectionHeaderType

	Flags SectionHeaderFlag

	Address uint64

	Offset uint64

	Size uint64

	Link uint32

	Info uint32

	AddressAlign uint64
	EntSize      uint64
}

type SectionHeaderType uint32

const (
	SHT_NULL SectionHeaderType = iota
	SHT_PROGBITS
	SHT_SYMTAB
	SHT_STRTAB
	SHT_RELA
	SHT_HASH
	SHT_DYNAMIC
	SHT_NOTE
	SHT_NOBITS
	SHT_REL
	SHT_SHLIB
	SHT_DYNSYM

	SHT_LOPROC SectionHeaderType = 0x70000000
	SHT_HIPROC SectionHeaderType = 0x7fffffff
	SHT_LOUSER SectionHeaderType = 0x80000000
	SHT_HIUSER SectionHeaderType = 0xffffffff
)

type SectionHeaderFlag uint64

const (
	SHF_WRITE            SectionHeaderFlag = 0x1        // Writable
	SHF_ALLOC            SectionHeaderFlag = 0x2        // Occupies memory during execution
	SHF_EXECINSTR        SectionHeaderFlag = 0x4        // Executable
	SHF_MERGE            SectionHeaderFlag = 0x10       // Might be merged
	SHF_STRINGS          SectionHeaderFlag = 0x20       // Contains null-terminated strings
	SHF_INFO_LINK        SectionHeaderFlag = 0x40       // 'sh_info' contains SHT index
	SHF_LINK_ORDER       SectionHeaderFlag = 0x80       // Preserve order after combining
	SHF_OS_NONCONFORMING SectionHeaderFlag = 0x100      // Non-standard OS specific handling required
	SHF_GROUP            SectionHeaderFlag = 0x200      // Section is member of a group
	SHF_TLS              SectionHeaderFlag = 0x400      // Section hold thread-local data
	SHF_MASKOS           SectionHeaderFlag = 0x0FF00000 // OS-specific
	SHF_MASKPROC         SectionHeaderFlag = 0xF0000000 // Processor-specific
	SHF_ORDERED          SectionHeaderFlag = 0x4000000  // Special ordering requirement (Solaris)
	SHF_EXCLUDE          SectionHeaderFlag = 0x8000000  // Section is excluded unless referenced or allocated (Solaris)
)

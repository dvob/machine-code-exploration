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
	Entry   uint64 // Entry point for executable files or zero for relocatable files, shared objects and the like

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

// SectionHeader64 see chapter 4 - 10 in System V Application Binary Interface
type SectionHeader64 struct {
	// Name specifies the name of the section. Its value is an index into
	// the section header string table section [see "String Table'' below],
	// giving the location of a null-terminated string
	Name uint32

	// Type categorizes the section's contents and semantics.
	Type SectionHeaderType

	Flags SectionHeaderFlag

	// Address at which the section's first byte should reside in memory.
	// Only applicable if the section will appear in the memory image of a
	// process. therwise, the member contains 0.
	Address uint64

	// Offset gives the byte offset from the beginning of the file to the
	// rst byte in the section. One section type, SHT_NOBITS described
	// below, occupies no space in the file, and its sh_offset member
	// locates the conceptual placement in the file.
	Offset uint64

	// Size gives the section’s size in bytes. Unless the section type is
	// SHT_NOBITS, the section occupies sh_size bytes in the file. A
	// section of type SHT_NOBITS may have a non-zero size, but it occupies
	// no space in the file.
	Size uint64

	// Link holds a section header table index link, whose interpretation
	// depends on the section type. A table below describes the values.
	Link uint32

	// Info holds extra information, whose interpretation depends on the
	// section type.
	Info uint32

	// Some sections have address alignment constraints. For example, if a
	// section holds a doubleword, the system must ensure doubleword
	// alignment for the entire section. That is, the value of sh_addr must
	// be congruent to 0, modulo the value of sh_addralign. Currently, only
	// 0 and positive integral powers of two are allowed. Values 0 and 1
	// mean the section has no alignment constraints.
	AddressAlign uint64

	// Some sections hold a table of fixed-size entries, such as a symbol
	// table. For such a section, this member gives the size in bytes of
	// each entry. The member contains 0 if the section does not hold a
	// table of fixed-size entries
	EntSize uint64
}

type SectionHeaderType uint32

const (
	SHT_NULL SectionHeaderType = iota

	// The section holds information defined by the program, whose format
	// and meaning are determined solely by the program. For example used
	// for the code and variables (.text, .data, .bss).
	SHT_PROGBITS

	// SYMTAB and DYNSYM hold a symbol table. Currently, an object file may
	// have only one section of each type, but this restriction may be
	// relaxed in the future. Typically, SHT_SYMTAB provides symbols for
	// link editing, though it may also be used for dynamic linking. As a
	// complete symbol table, it may contain many symbols unnecessary for
	// dynamic linking. Consequently, an object file may also contain a
	// SHT_DYNSYM section, which holds a minimal set of dynamic linking
	// symbols, to save space. These point to a symbol tables.
	SHT_SYMTAB

	// The section holds a string table. An object file may have multiple
	// string table sections.
	SHT_STRTAB

	// The section holds relocation entries with explicit addends. An
	// object file may have multiple relocation sections
	SHT_RELA

	// The section holds a symbol hash table. All objects participating in
	// dynamic linking must contain a symbol hash table. Currently, an
	// object file may have only one hash table, but this restriction may
	// be relaxed in the future
	SHT_HASH

	// The section holds information for dynamic linking. Currently, an
	// object file may have only one dynamic section, but this restriction
	// may be relaxed in the future
	SHT_DYNAMIC

	// The section holds information that marks the ﬁle in some way.
	SHT_NOTE

	// A section of this type occupies no space in the file but other- wise
	// resembles SHT_PROGBITS. Although this section contains no bytes, the
	// sh_offset member contains the conceptual file offset. Typically used
	// by uninitialized variables (.bss).
	SHT_NOBITS

	// The section holds relocation entries without explicit addends. An
	// object file may have multiple relocation sections.
	SHT_REL

	// This section type is reserved but has unspecified semantics.
	// Programs that contain a section of this type do not conform to the
	// ABI.
	SHT_SHLIB

	// See SHT_SYMTAB
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

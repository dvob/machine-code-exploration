package main

//go:generate stringer -type FileType,Class,Data,ProgramHeaderFlag,ProgramHeaderType,SectionHeaderFlag,SectionHeaderType -output string.go

// ELFFile combines the various information a ELF file could contain. But this
// struct can't be read using binary.Read as only the header is guaranteed be
// be at the beginning.
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

// NewSectionHeaderTable64 returns a section header table with the special null
// entry at the beginning
func NewSectionHeaderTable64() []SectionHeader64 {
	return []SectionHeader64{
		{},
	}
}

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

	SHT_INIT_ARRAY    SectionHeaderType = 14
	SHT_FINI_ARRAY    SectionHeaderType = 15
	SHT_PREINIT_ARRAY SectionHeaderType = 16
	SHT_GROUP         SectionHeaderType = 17
	SHT_SYMTAB_SHNDX  SectionHeaderType = 18
	SHT_LOOS          SectionHeaderType = 0x60000000
	SHT_HIOS          SectionHeaderType = 0x6fffffff
	SHT_LOPROC        SectionHeaderType = 0x70000000
	SHT_HIPROC        SectionHeaderType = 0x7fffffff
	SHT_LOUSER        SectionHeaderType = 0x80000000
	SHT_HIUSER        SectionHeaderType = 0xffffffff
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

// NewSymbolTable64 creates a symbol table with the first entry set to the specail null value
func NewSymbolTable64() []Symbol64 {
	return []Symbol64{
		{},
	}
}

// Symbol64 represents a 64-bit ELF symbol table entry.
// NOTE: Field order differs from ELF32! In ELF64, the smaller fields (Info,
// Other, SectionHeaderIndex) come before Value/Size for proper 8-byte alignment.
type Symbol64 struct {
	// Name holds an index into the object file’s symbol string table. The
	// table (section index) is stored in the symbol tables sections Link
	// field. Usually it is .strtab for .symtab and .dynstr for .dyntab. If
	// the value is zero the symbol table entry has no name.
	Name uint32

	// Info specifies the symbol's type and binding attributes:
	//   - Binding: the upper 4 bits. See SymbolBinding
	//   - Type: lower 4 bits. See SymbolType
	Info uint8

	// Other originally had no special meaning. Latter its 3 lower bits got used for visibility.
	Other uint8

	// SectionHeaderIndex indicates which section this symbol is defined in:
	//   - SHN_UNDEF (0): Undefined symbol (needs to be resolved by linker)
	//   - SHN_ABS (0xfff1): Absolute value (Value is not affected by relocation)
	//   - SHN_COMMON (0xfff2): Common block (unallocated data, like extern int x;)
	//   - 1..n: Index into section header table
	//
	// Example: SectionHeaderIndex = 1 means defined in sections[1] (often .text)
	SectionHeaderIndex uint16

	// Value holds the symbol's value, meaning depends on context:
	//   - Relocatable files (ET_REL): offset from beginning of section
	//   - Executable/shared (ET_EXEC/ET_DYN): virtual address
	Value uint64

	// Size gives the associated object's size in bytes:
	//   - Functions: number of bytes of code
	//   - Variables: number of bytes of data
	//   - 0: no size or unknown size
	//
	// Example: A 4-byte int variable has Size = 4
	Size uint64
}

func (s Symbol64) SymbolBinding() SymbolBinding {
	return SymbolBinding(s.Info >> 4)
}

func (s Symbol64) SymbolType() SymbolType {
	return SymbolType(s.Info & 0xf)
}

func (s Symbol64) SymbolVisibility() SymbolVisbility {
	return SymbolVisbility(s.Other & 0x3)
}

type SymbolType uint8

const (
	STT_NOTYPE  SymbolType = 0  // Symbol type not specified
	STT_OBJECT  SymbolType = 1  // Data object (variable)
	STT_FUNC    SymbolType = 2  // Function or executable code
	STT_SECTION SymbolType = 3  // Section symbol (used for relocations)
	STT_FILE    SymbolType = 4  // Source file name symbol
	STT_COMMON  SymbolType = 5  // Uninitialized common block
	STT_TLS     SymbolType = 6  // Thread-Local Storage entity
	STT_LOOS    SymbolType = 10 // OS-specific semantics
	STT_HIOS    SymbolType = 12 // OS-specific semantics
	STT_LOPROC  SymbolType = 13 // Processor-specific semantics
	STT_HIPROC  SymbolType = 15 // Processor-specific semantics
)

type SymbolBinding uint8

const (
	STB_LOCAL  SymbolBinding = 0  // Local symbol (not visible outside object file)
	STB_GLOBAL SymbolBinding = 1  // Global symbol (visible to all object files)
	STB_WEAK   SymbolBinding = 2  // Weak symbol (like global but lower precedence)
	STB_LOOS   SymbolBinding = 10 // OS-specific semantics
	STB_HIOS   SymbolBinding = 12 // OS-specific semantics
	STB_LOPROC SymbolBinding = 13 // Processor-specific semantics
	STB_HIPROC SymbolBinding = 15 // Processor-specific semantics
)

type SymbolVisbility uint8

const (
	STV_DEFAULT   SymbolVisbility = 0
	STV_INTERNAL  SymbolVisbility = 1
	STV_HIDDEN    SymbolVisbility = 2
	STV_PROTECTED SymbolVisbility = 3
)

// NewSymbolInfo combines binding and type into Info byte
func NewSymbolInfo(binding SymbolBinding, symbolType SymbolType) uint8 {
	return (uint8(binding) << 4) | (uint8(symbolType) & 0xf)
}

type SectionIndex uint16

// Special section indices
const (
	SHN_UNDEF     SectionIndex = 0      // Undefined section
	SHN_LORESERVE SectionIndex = 0xff00 // Lower bound of reserved indexes
	SHN_LOPROC    SectionIndex = 0xff00 // Processor-specific semantics
	SHN_HIPROC    SectionIndex = 0xff1f // Processor-specific semantics
	SHN_LOOS      SectionIndex = 0xff20 // OS-specific semantics
	SHN_HIOS      SectionIndex = 0xff3f // OS-specific semantics
	SHN_ABS       SectionIndex = 0xfff1 // Absolute values
	SHN_COMMON    SectionIndex = 0xfff2 // Common block
	SHN_XINDEX    SectionIndex = 0xffff // Escape value for large section index
	SHN_HIRESERVE SectionIndex = 0xffff // Upper bound of reserved indexes
)

type SectionGroupFlag uint32

// Section group flags
const (
	GRP_COMDAT   SectionGroupFlag = 0x1        // COMDAT group
	GRP_MASKOS   SectionGroupFlag = 0x0ff00000 // OS-specific semantics
	GRP_MASKPROC SectionGroupFlag = 0xf0000000 // Processor-specific semantics
)

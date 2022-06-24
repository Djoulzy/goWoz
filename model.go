package gowoz

import (
	"os"

	"github.com/bits-and-blooms/bitset"
)

const (
	INFO_CHUNK_ID = 0x4F464E49
	TMAP_CHUNK_ID = 0x50414D54
	TRKS_CHUNK_ID = 0x534B5254
	META_CHUNK_ID = 0x4154454D
)

var (
	DiskType   = []string{"Unknown", "5.25", "3.5"}
	BootSector = []string{"Unknown", "16-sector", "13-sector", "Both"}
)

type WOZHeader struct {
	Format   [4]byte
	HighBits byte
	LFCRLF   [3]byte
	CRC      [4]byte
}

type WOZChunkHeader struct {
	ID   string
	Size uint32
}

type WOZChunkMeta struct {
	Header   WOZChunkHeader
	Metadata string
}

type WOZInfoChunk struct {
	Header             WOZChunkHeader
	Version            uint8
	DiskType           uint8
	WriteProtected     uint8
	Synchronized       uint8
	Cleaned            uint8
	Creator            [32]byte
	DiskSides          uint8
	BootSectorFormat   uint8
	OptimalBitTiming   uint8
	CompatibleHardware uint16
	RequiredRAM        uint16
	LargestTrack       uint16
	FLUXBlock          uint16
	LargestFluxTrack   uint16
}

type WOZTMapChunk struct {
	Header WOZChunkHeader
	Map    map[float32]byte
}

type WOZTrackDesc struct {
	StartBlock uint16
	BlockCount uint16
	BitCount   uint32
}

type WOZTRKSChunk struct {
	Header WOZChunkHeader
	Tracks [160]WOZTrackDesc
	Data   [160]*bitset.BitSet
}

type WOZFileFormat struct {
	fdesc  *os.File
	Header WOZHeader
	INFO   WOZInfoChunk
	TMAP   WOZTMapChunk
	META   WOZChunkMeta
	TRKS   WOZTRKSChunk
}

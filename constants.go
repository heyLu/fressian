package fressian

// Byte codes for the types fressian supports
const (
	PRIORITY_CACHE_PACKED_START = 0x80
	PRIORITY_CACHE_PACKED_END   = 0xA0
	STRUCT_CACHE_PACKED_START   = 0xA0
	STRUCT_CACHE_PACKED_END     = 0xB0
	LONG_ARRAY                  = 0xB0
	DOUBLE_ARRAY                = 0xB1
	BOOLEAN_ARRAY               = 0xB2
	INT_ARRAY                   = 0xB3
	FLOAT_ARRAY                 = 0xB4
	OBJECT_ARRAY                = 0xB5
	MAP                         = 0xC0
	SET                         = 0xC1
	CODE_UUID                   = 0xC3
	REGEX                       = 0xC4
	URI                         = 0xC5
	BIGINT                      = 0xC6
	BIGDEC                      = 0xC7
	INST                        = 0xC8
	SYM                         = 0xC9
	KEY                         = 0xCA
	GET_PRIORITY_CACHE          = 0xCC
	PUT_PRIORITY_CACHE          = 0xCD
	PRECACHE                    = 0xCE
	FOOTER                      = 0xCF
	FOOTER_MAGIC                = 0xCFCFCFCF
	BYTES_PACKED_LENGTH_START   = 0xD0
	BYTES_PACKED_LENGTH_END     = 0xD8
	BYTES_CHUNK                 = 0xD8
	BYTES                       = 0xD9
	STRING_PACKED_LENGTH_START  = 0xDA
	STRING_PACKED_LENGTH_END    = 0xE2
	STRING_CHUNK                = 0xE2
	STRING                      = 0xE3
	LIST_PACKED_LENGTH_START    = 0xE4
	LIST_PACKED_LENGTH_END      = 0xEC
	LIST                        = 0xEC
	BEGIN_CLOSED_LIST           = 0xED
	BEGIN_OPEN_LIST             = 0xEE
	STRUCTTYPE                  = 0xEF
	STRUCT                      = 0xF0
	META                        = 0xF1
	ANY                         = 0xF4
	TRUE                        = 0xF5
	FALSE                       = 0xF6
	NULL                        = 0xF7
	INT                         = 0xF8
	FLOAT                       = 0xF9
	DOUBLE                      = 0xFA
	DOUBLE_0                    = 0xFB
	DOUBLE_1                    = 0xFC
	END_COLLECTION              = 0xFD
	RESET_CACHES                = 0xFE
	INT_PACKED_1_START          = 0xFF
	INT_PACKED_1_END            = 0x40
	INT_PACKED_2_START          = 0x40
	INT_PACKED_2_ZERO           = 0x50
	INT_PACKED_2_END            = 0x60
	INT_PACKED_3_START          = 0x60
	INT_PACKED_3_ZERO           = 0x68
	INT_PACKED_3_END            = 0x70
	INT_PACKED_4_START          = 0x70
	INT_PACKED_4_ZERO           = 0x72
	INT_PACKED_4_END            = 0x74
	INT_PACKED_5_START          = 0x74
	INT_PACKED_5_ZERO           = 0x76
	INT_PACKED_5_END            = 0x78
	INT_PACKED_6_START          = 0x78
	INT_PACKED_6_ZERO           = 0x7A
	INT_PACKED_6_END            = 0x7C
	INT_PACKED_7_START          = 0x7C
	INT_PACKED_7_ZERO           = 0x7E
	INT_PACKED_7_END            = 0x80
)

const (
	BYTE_CHUNK_SIZE        = 65535
	STRING_CHUNK_MAX_SIZE  = 65536
	STRING_PACKED_MAX_SIZE = STRING_PACKED_LENGTH_END - STRING_PACKED_LENGTH_START
	BYTES_PACKED_MAX_SIZE  = BYTES_PACKED_LENGTH_END - BYTES_PACKED_LENGTH_START
	LIST_PACKED_MAX_SIZE   = LIST_PACKED_LENGTH_END - LIST_PACKED_LENGTH_START
	STRUCT_CACHE_MAX_SIZE  = STRUCT_CACHE_PACKED_END - STRUCT_CACHE_PACKED_START
)

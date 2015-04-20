# Reading fressian data in Go

**Caution:** This is very early stages. It can read fressian data, but
not all of fressian is supported and the API is definitely *not* final
in any way. Say hi if you try to use it. :)

## Usage

`go get github.com/heyLu/fressian`

- create a reader with `fressian.NewReader(r, nil)`
- use `.ReadValue()` to read the next object
- see [./fsn](./fsn/main.go) for an example

## TODO

- improving the API
    - better error handling (e.g. don't crash on errors)
    - maybe export `Read{Int,...}` and friends
- implement the remaining bytecodes
- writing fressian data

## Examples

- numbers from 0 to 63 (0x00 .. 0x3F) are encoded directly
- numbers from -4096 to 4095 are encoded using 2 bytes (`INT_2_PACKED`)

    the first byte is used to encode the high 8 bits, the second byte
    encodes the lower 8 bits

    the number is calculated using `(hi - 0x50 << 8) | lo`.

- floating point numbers (32 bit = 4 bytes), `(float 1.2345)`

                  FLOAT
                  ^
        00000000  f9 3f 9e 04 19                                    |.?...|
        00000005
- double-precision (64 bit = 8 bytes) floating point numbers, `3.257329852835`

                  DOUBLE
                  ^
        00000000  fa 40 0a 0f 02 f4 31 af  c1                       |.@....1..|
        00000009

    `0.0` and `1.0` have special encodings, `0xFB` and `0xFC`
- `[1 2 3 4 5]`

                  LIST_PACKED_LENGTH_START + 5 = 0xe4 + 5
                  ^
        00000000  e9 01 02 03 04 05                                 |......|
        00000006
- `[1 2 3 4 5 "hello"]`

                  LIST_PACKED_LENGTH_START + 6 = 0xe4 + 6
                  ^
                  |                 STRING_PACKED_LENGTH_START + 5 = 0xda + 5
                  |                 ^
        00000000  ea 01 02 03 04 05 df 68  65 6c 6c 6f              |.......hello|
        0000000c
- `{"hey" 3, "ho" 2, "answer" 42}`

                     LIST_PACKED_LENGTH_START + 6 = 0xe4 + 6
                     ^
                     |  STRING_PACKED_LENGTH_START + 3
                     |  ^   h  e  y
                  MAP|  |            3  STRING_PACKED_LENGTH + 2
                  ^  |  |               ^   h  o  2     a  n  s  w
        00000000  c0 ea dd 68 65 79 03 dc  68 6f 02 e0 61 6e 73 77  |...hey..ho..answ|
        00000010  65 72 2a                                          |er*|
        00000013   |
                   e  r 42

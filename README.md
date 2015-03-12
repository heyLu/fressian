# Exploring fressian

Before starting to play, run `make deps` to fetch various dependencies.

## Using `Datomic/fressian`

```
$ pushd fressian; mvn compile; popd
...
$ clj -cp fressian/target/classes
Clojure 1.6.0
=> (import '(org.fressian FressianReader FressianWriter))
=> (import '(java.io File FileInputStream FileOutputStream))
=>
=> (def f (File. "example.fressian"))
=>
=> (def w (FressianWriter. (FileOutputStream f)))
=> (.writeObject w [1 2 3 4 5])
=> (.close w)

=> (def r (FressianReader. (FileInputStream. f)))
=> (.readObject r)
[1 2 3 4 5]
=> (.close r)
```

Then you can inspect `example.fressian` using `hexdump -C
example.fressian`, for example.

## Examples

- numbers from 0 to 63 (0x00 .. 0x3F) are encoded directly
- numbers from -4096 to 4095 are encoded using 2 bytes (`INT_2_PACKED`)

    the first byte is used to encode the high 8 bits, the second byte
    encodes the lower 8 bits

    the number is calculated using `(hi - 0x50 << 8) | lo`.
- `(.writeObject w [1 2 3 4 5])`, `.readObject` returns an `ArrayList`

                  LIST_PACKED_LENGTH_START + 5 = 0xe4 + 5
                  ^
        00000000  e9 01 02 03 04 05                                 |......|
        00000006
- `(.writeObject w [1 2 3 4 5 "hello"])`, also an `ArrayList`

                  LIST_PACKED_LENGTH_START + 6 = 0xe4 + 6
                  ^
                  |                 STRING_PACKED_LENGTH_START + 5 = 0xda + 5
                  |                 ^
        00000000  ea 01 02 03 04 05 df 68  65 6c 6c 6f              |.......hello|
        0000000c

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

- `(.writeObject w [1 2 3 4 5])`, `.readObject` returns an `ArrayList`

        00000000  e9 01 02 03 04 05                                 |......|
        00000006
- `(.writeObject w [1 2 3 4 5 "hello"])`, also an `ArrayList`

        00000000  ea 01 02 03 04 05 df 68  65 6c 6c 6f              |.......hello|
        0000000c

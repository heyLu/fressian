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

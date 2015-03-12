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
=> (def w (FressianWriter. (FileOutputStream (File. "example.fressian"))))
=> (.writeObject w [1 2 3 4 5])
=> (.close w)

=> (def r (FressianReader. (FileInputStream. (File. "example.fressian"))))
=> (.readObject r)
[1 2 3 4 5]
=> (.close r)
```

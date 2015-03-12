build: misc/adler32

misc/adler32: misc/adler32.c
	clang -lz $< -o $@

deps:
	[ -e fressian ] || git clone git://github.com/Datomic/fressian
	[ -e data.fressian ] || git clone git://github.com/clojure/data.fressian
	curl -s http://zlib.net/zlib-1.2.8.tar.gz | tar -xzv
	curl -sO https://tools.ietf.org/rfc/rfc1950.txt

clean:
	rm -rf fressian data.fressian zlib* rfc*
	rm -f misc/adler32

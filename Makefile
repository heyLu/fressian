build: adler32

adler32: adler32.c
	clang -lz $< -o $@

deps:
	[ -e fressian ] || git clone git://github.com/Datomic/fressian
	curl -s http://zlib.net/zlib-1.2.8.tar.gz | tar -xzv
	curl -sO https://tools.ietf.org/rfc/rfc1950.txt

clean:
	rm -rf fressian zlib* rfc*
	rm -f adler32

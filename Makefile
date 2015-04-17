all:

deps:
	[ -e deps/fressian ] || git clone git://github.com/Datomic/fressian deps/fressian

clean:
	rm -rf fressian

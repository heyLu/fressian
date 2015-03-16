all:

deps:
	[ -e fressian ] || git clone git://github.com/Datomic/fressian

clean:
	rm -rf fressian

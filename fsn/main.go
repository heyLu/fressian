package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	fressian "github.com/heyLu/fressian"
)

func prettySprint(value interface{}) string {
	switch value := value.(type) {
	case bool, byte, int, float32, float64, string:
		return fmt.Sprintf("%#v", value)
	case fressian.Key:
		if value.Namespace == "" {
			return fmt.Sprintf(":%s", value.Name)
		} else {
			return fmt.Sprintf(":%s/%s", value.Namespace, value.Name)
		}
	default:
		log.Fatalf("unexpected value: %#v", value)
		return ""
	}
}

func prettyPrint(indent string, value interface{}) {
	switch value := value.(type) {
	case time.Time:
		fmt.Printf("%s#inst \"%s\"\n", indent, value.Format(time.RFC3339))

	case *url.URL:
		fmt.Printf("%s#uri \"%s\"\n", indent, value)

	case fressian.StructAny:
		fmt.Printf("%s#%s [\n", indent, value.Tag)
		for _, val := range value.Values {
			prettyPrint(indent+"  ", val)
		}
		fmt.Printf("%s]\n", indent)

	case map[interface{}]interface{}:
		fmt.Printf("%s{\n", indent)
		for key, val := range value {
			switch val.(type) {
			case bool, byte, int, float32, float64, string, fressian.Key:
				fmt.Printf("%s%s %s\n", indent+"  ", prettySprint(key), prettySprint(val))
			default:
				prettyPrint(indent+"  ", key)
				prettyPrint(indent+"    ", val)
			}
		}
		fmt.Printf("%s}\n", indent)

	case []interface{}:
		fmt.Printf("%s[\n", indent)
		for _, val := range value {
			prettyPrint(indent+"  ", val)
		}
		fmt.Printf("%s]\n", indent)

	case []bool:
		fmt.Printf("%s#booleans [", indent)
		length := len(value)
		for i, val := range value {
			if i != length-1 {
				fmt.Printf("%t, ", val)
			} else {
				fmt.Printf("%t]\n", val)
			}
		}

	case []byte:
		fmt.Printf("%s#bytes [", indent)
		length := len(value)
		for i, val := range value {
			if i != length-1 {
				fmt.Printf("0x%x, ", val)
			} else {
				fmt.Printf("0x%x]\n", val)
			}
		}

	case []int:
		fmt.Printf("%s#ints [", indent)
		length := len(value)
		for i, val := range value {
			if i != length-1 {
				fmt.Printf("%d, ", val)
			} else {
				fmt.Printf("%d]\n", val)
			}
		}

	case []float32:
		fmt.Printf("%s#floats [", indent)
		length := len(value)
		for i, val := range value {
			if i != length-1 {
				fmt.Printf("%f, ", val)
			} else {
				fmt.Printf("%f]\n", val)
			}
		}

	case []float64:
		fmt.Printf("%s#doubles [", indent)
		length := len(value)
		for i, val := range value {
			if i != length-1 {
				fmt.Printf("%f, ", val)
			} else {
				fmt.Printf("%f]\n", val)
			}
		}

	case fressian.Key:
		fmt.Printf("%s%s\n", indent, prettySprint(value))

	default:
		if value == nil {
			fmt.Printf("%snil\n", indent)
		} else {
			fmt.Printf("%s%#v\n", indent, value)
		}
	}
}

func isGzipped(r *bufio.Reader) bool {
	magic, err := r.Peek(2)
	if err != nil {
		log.Fatal(err)
	}
	return magic[0] == 0x1F && magic[1] == 0x8B
}

var pretty = flag.Bool("p", false, "pretty print the value read")

func main() {
	flag.Parse()

	var f io.Reader
	if flag.NArg() == 0 {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	}
	f = bufio.NewReader(f)

	if isGzipped(f.(*bufio.Reader)) {
		var err error
		f, err = gzip.NewReader(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := fressian.NewReader(f, nil)
	obj, err := r.ReadObject()
	if err != nil {
		log.Fatal(err)
	}

	if *pretty {
		prettyPrint("", obj)
	} else {
		fmt.Printf("%#v\n", obj)
	}
}

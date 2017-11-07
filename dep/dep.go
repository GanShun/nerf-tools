package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nerfirmware/tools/pkg/guid"
)

const (
	glen = 16
)

var (
	opcodes = map[byte]string{
		6: "TRUE",
		2: "PUSH",
		3: "AND",
		4: "OR",
		8: "END",
	}
)

func readGUID(r io.Reader) guid.GUID {
	var g guid.GUID
	var b []byte
	b = make([]byte, glen, glen)
	if n, err := r.Read(b[:]); n != glen || err != nil {
		log.Fatalf("Reading GUID: got %d of %d bytes: %v", n, glen, err)
	}
	if err := g.UnmarshalBinary(b); err != nil {
		log.Fatalf("Unmarshalling GUID, err: %v", err)
	}
	return g
}

func main() {
	var (
		hdr [4]byte
		op  [1]byte
	)

	if n, err := os.Stdin.Read(hdr[:]); n != len(hdr) || err != nil {
		log.Fatalf("Reading header: %v", err)
	}

	if hdr[3] != 0x13 {
		log.Fatalf("Header type is 0x%x, not 0x13", hdr[3])
	}

	l := int(hdr[0]) + (int(hdr[1]) << 8) + (int(hdr[2]) << 16) - len(hdr)
	//fmt.Fprintf(os.Stderr, "hdr: %v; len is 0x%x bytes\n", hdr, l)

	for i := 0; i < l; {
		n, err := os.Stdin.Read(op[:])
		i += n
		if err != nil && (err != io.EOF || n == 0) {
			if err == io.EOF {
				log.Fatalf("Header length %d is longer than the file length %d", l, i)
			}
			log.Fatalf("dep file: %v", err)
		}

		fmt.Printf("%v", opcodes[op[0]])

		if op[0] == 2 {
			g := readGUID(os.Stdin)
			fmt.Printf(" %v", g.String())
			i += glen
		}

		fmt.Printf("\n")
	}
}

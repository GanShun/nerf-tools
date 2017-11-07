package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nerfirmware/tools/pkg/guid"
)

var (
	opcodes = map[string]byte{
		"TRUE": 6,
		"PUSH": 2,
		"AND":  3,
		"OR":   4,
		"END":  8,
	}
	stackdepth int
)

func main() {
	var (
		b    bytes.Buffer
		done bool
	)

	for !done {
		var op, g string

		depGUID := new(guid.GUID)
		n, err := fmt.Scanln(&op, &g)
		if err == io.EOF {
			break
		}
		if n == 0 {
			continue
		}
		//fmt.Printf("%v %v\n", op, g)
		opcode, ok := opcodes[op]
		if !ok {
			log.Fatalf("Opcode %v not known", opcode)
		}
		b.Write([]byte{opcode})
		if op == "PUSH" {
			if err := depGUID.Parse(g); err != nil {
				log.Fatalf("Error parsing guid: %v", err)
			}
			if binGUID, err := depGUID.MarshalBinary(); err != nil {
				log.Fatalf("Error marshalling binary: %v", err)
			} else {
				b.Write(binGUID)
			}
		}
		switch op {
		case "TRUE":
			stackdepth++
		case "PUSH":
			stackdepth++
		case "AND":
			stackdepth--
		case "OR":
			stackdepth--
		case "END":
			done = true
		}

	}
	l := b.Len() + 4
	hdr := append([]byte{byte(l), byte(l >> 8), byte(l >> 16), 0x13}, b.Bytes()...)
	if _, err := os.Stdout.Write(hdr); err != nil {
		log.Fatalf("%v", err)
	}
	if stackdepth != 1 {
		log.Fatalf("stackdepth is %d, should be 1", stackdepth)
	}

}

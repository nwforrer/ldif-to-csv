package main

import (
	"os"
	"log"
	"github.com/nwforrer/ldif-to-csv"
	"fmt"
)

func main() {
	prog := os.Args[0]
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: %s <input-file> [output-file]\n", prog)
		return
	}

	inFilename := args[0]

	inFile, err := os.Open(inFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	headers := []string{"cn", "owner"}

	ldifEntries := ldiftocsv.ReadLdifFile(inFile, headers)

	out := os.Stdout
	if len(args) > 1 {
		outFile, err := os.Create(args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		out = outFile
	}
	ldiftocsv.WriteToCsv(out, ldifEntries)
}

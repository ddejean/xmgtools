// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"xioxoz.fr/hwpreader/rtl"
)

var (
	file   = flag.String("f", "", "file to parse")
	offset = flag.Int64("o", int64(0), "offset in the file")
)

func main() {
	// Disable date and timestamps.
	log.SetFlags(0)

	// Parse the command line flags.
	flag.Parse()

	if file == nil || *file == "" {
		log.Fatal("input file required")
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Seek(*offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	s := &rtl.Switch{}
	if err := s.UnmarshalBinary(bufio.NewReader(f)); err != nil {
		log.Fatal(err)
	}

	log.Print(s)
}

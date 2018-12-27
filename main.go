package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//validFormats := []string{"INTEL", "MOTO", "TEX"}

	var outFile string
	flag.StringVar(&outFile, "o", "dump.bin", "File to dump extracted data.")
	var hexFormat string
	flag.StringVar(&hexFormat, "f", "INTEL", "Format of data. Must be INTEL/MOTO/TEX.")
	fmt.Println(os.Args)
	flag.Parse()
	inputFile := flag.Args()[0]

	fmt.Printf("Reading from %s in %s format. Dumping to %s\n", inputFile, hexFormat, outFile)

	if hexFormat != "INTEL" {
		fmt.Println("Unsupported format " + hexFormat)
		flag.Usage()
		os.Exit(1)
	}

	ifile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Failed to open input file. Err:", err)
		flag.Usage()
		os.Exit(2)
	}

	ofile, err := os.Create(outFile)
	if err != nil {
		fmt.Println("Failed to open output file. Err:", err)
		flag.Usage()
		os.Exit(3)
	}

	iHexParser := IntelHex{
		input:   ifile,
		output:  ofile,
		records: make([]*HexRecord, 0)}

	iHexParser.Parse()
	for record := range iHexParser.records {
		fmt.Printf("%+v\n", *(iHexParser.records[record]))
	}

	iHexParser.Dump()
	for addr, data := range iHexParser.buffers {
		fmt.Printf("%8x : %v\n", addr, data)
	}
}

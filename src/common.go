package main

import (
	"bufio"
	"os"
)

//HexRecord defines strcure of a line in HexDump.
type HexRecord struct {
	byteCount  int8
	offset     int16
	recordType int8
	data       []byte
	checksum   byte
	linearAddr int32
}

//HexFormat is an iterface for pasring hexdumps.
type HexFormat struct {
	input   *os.File
	output  *os.File
	records []*HexRecord
}

// Parse is an abstract definition to parse input file content.
func (hxfmt *HexFormat) Parse() bool {
	return false
}

// Dump dumps content to outfile.
func (hxfmt *HexFormat) Dump() int32 {
	var preRecord HexRecord
	writer := bufio.NewWriter(hxfmt.output)
	contentLength := 0
	for record := range hxfmt.records {

	}
}

package main

import (
	"os"
)

//HexRecord defines strcture of a line in HexDump.
type HexRecord struct {
	byteCount  uint8
	offset     uint16
	recordType uint8
	data       []byte
	checksum   byte
	linearAddr uint32
	lineNumber uint
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
/* func (hxfmt *HexFormat) Dump() int32 {
	var preRecord HexRecord
	writer := bufio.NewWriter(hxfmt.output)
	contentLength := 0
	for record := range hxfmt.records {

	}
	return
}
*/

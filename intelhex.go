package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
)

//RecordType : valid record types
type RecordType uint8

//RecordType : valid record types
const (
	DATA RecordType = iota
	EOF
	EXTSEGMENT
	EXTLINEAR
	STARTLINEAR
)

//RecordTypes : name of record types
var RecordTypes = [...]string{
	"DATA",
	"EOF",
	"EXTENDED_SEGMENT_ADDR",
	"EXTENDED_LINEAR_ADDR",
	"START_LINEAR_ADDR",
}

// String : string representation for RecordTypes
func (rType RecordType) String() string {
	return RecordTypes[rType]
}

//IntelHex : derivation of HexFormat
type IntelHex HexFormat

//Parse : Parse IntelHex file.
func (ihex *IntelHex) Parse() bool {
	scanner := bufio.NewScanner(ihex.input)
	//var segmentAddr int32
	var lineNumber uint

	for scanner.Scan() {
		line := scanner.Text()
		lineLength := len(line)

		if lineLength <= 10 {
			fmt.Print("WARN: Skipped line (Must have atleast 11 characters):", line)
			continue
		}

		if line[0] != ':' {
			fmt.Println("WARN", "Skipped line (Missing start marker):", line)
			continue
		}

		record := &HexRecord{}
		record.lineNumber = lineNumber
		lineNumber++

		bValue := line[1:3]
		iTemp, err := strconv.ParseUint(bValue, 16, 8)
		if err == nil {
			record.byteCount = uint8(iTemp)
		} else {
			fmt.Println("WARN: Skipped; Hex (", bValue, ") is not a valid record size in line:", line)
			continue
		}

		bValue = line[3:7]
		iTemp, err = strconv.ParseUint(bValue, 16, 16)
		if err == nil {
			record.offset = uint16(iTemp)
		} else {
			fmt.Println("WARN: Skipped; Hex (", bValue, ") is not a valid record offset in line:", line, ". Err:", err)
			continue
		}

		bValue = line[7:9]
		iTemp, err = strconv.ParseUint(bValue, 16, 8)
		if err == nil {
			record.recordType = uint8(iTemp)
		} else {
			fmt.Println("WARN: Skipped; Hex (", bValue, ") is not a valid record type in line:", line, ". Err:", err)
			continue
		}

		bValue = line[(lineLength - 2):lineLength]
		bCheck, err := hex.DecodeString(bValue)
		if err == nil {
			record.checksum = bCheck[0]
		} else {
			fmt.Println("WARN:", bValue, "is not a valid checksum. Err:", err)

		}

		if record.byteCount > 0 {
			dataLen := int(record.byteCount) * 2

			if lineLength == dataLen+11 {
				dataBytes, err := hex.DecodeString(line[9 : dataLen+9])
				if err != nil {
					fmt.Println("Failed to decode data bytes, Err:", err)
					continue
				}
				record.data = dataBytes
			} else if lineLength < (int(record.byteCount)*2 + 11) {
				fmt.Println("WARN: Data bytes are less than specified length of", record.byteCount)
				continue
			} else {
				fmt.Println("WARN: Data bytes are more than specified length of", record.byteCount)
				continue
			}
		}

		ihex.records = append(ihex.records, record)
	}
	return true
}

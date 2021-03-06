package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

//IntelHex : derivation of HexFormat
type IntelHex HexFormat

//Parse : Parse IntelHex file.
func (ihex *IntelHex) Parse() bool {
	scanner := bufio.NewScanner(ihex.input)

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

// Dump : Dump parsed records
func (ihex *IntelHex) Dump() int {
	var segmentAddr uint32
	var linearAddr uint32
	var finalAddr uint32
	ihex.buffers = make(map[uint32][]byte)

	for _, record := range ihex.records {
		//Extended segment address record
		if record.recordType == 2 {
			segmentAddr = uint32(binary.BigEndian.Uint16(record.data))
			segmentAddr = segmentAddr * 16
		}
		// Linear address record
		if record.recordType == 4 {
			linearAddr = uint32(binary.BigEndian.Uint16(record.data))
			linearAddr = linearAddr << 16
		}
		// Data record
		if record.recordType == 0 {
			finalAddr = linearAddr + segmentAddr + uint32(record.offset)
			ihex.buffers[finalAddr] = record.data
		}
	}
	return len(ihex.buffers)
}

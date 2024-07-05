package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"unicode"
)

// main is the entry point of the program.
// It accepts command line arguments and sorts a slice of integers.
// The sorted elements are printed to the standard output using fmt.Println.
func main() {

	// ./intel-nvm-tool \
	//		-i <input stock nvm file> \
	//		-m <mac address / mac address range> \
	//		-o <output nvm file prefix>

	var inputFile, outputPrefix string
	var macStartAddr, macEndAddr net.HardwareAddr
	var err error
	args := os.Args[1:]

	for x := 0; x < len(args); x += 2 {

		if x > len(args)-2 {
			break
		}

		switch args[x] {
		case "-i":
			inputFile = args[x+1]
		case "-o":
			outputPrefix = args[x+1]
		case "-ms":
			if macStartAddr, err = net.ParseMAC(args[x+1]); err != nil {
				fmt.Printf("-ms address invalid: %s\n", err.Error())
				os.Exit(-1)
			}
		case "-me":
			if macEndAddr, err = net.ParseMAC(args[x+1]); err != nil {
				fmt.Printf("-me address invalid: %s\n", err.Error())
				os.Exit(-1)
			}
		}
	}

	if inputFile == "" || outputPrefix == "" || macStartAddr == nil {
		fmt.Println("-i -o and -ms are mandatory parameters to run")
		fmt.Printf("Usage: \t ./intel-nvm-tool \\\n\t\t\t-i <input stock nvm file> \\\n\t\t\t-m <mac address / mac address range> \\\n\t\t\t-o <output nvm file prefix>")
		os.Exit(-1)
	}

	file, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("unable to open input file: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("input file checksum: %X\n", file[126:128])

	macStart := binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(macStartAddr)...))
	macEnd := macStart
	if len(macEndAddr) == 6 {
		macEnd = binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(macEndAddr)...))
	}

	// Split based on function f
	//macParts := strings.FieldsFunc(macAddressRange, splitAny)
	//fmt.Printf("MAC address range: %#x - %#x\n", macStart, macEnd)
	fmt.Printf("MAC address range: %s - %s (%d addresses)\n",
		macStartAddr.String(),
		macEndAddr.String(),
		macEnd-macStart,
	)

	macBuffer := make([]byte, 8)
	for x := macStart; x < macEnd; x++ {
		binary.BigEndian.PutUint64(macBuffer, x)
		file = append(macBuffer[2:8], file[6:]...)

		var checksum uint16
		for y := 0; y < 0x7d; y += 2 {
			checksum += binary.BigEndian.Uint16(file[y : y+2])
			checksum += uint16(file[y])
		}
		checksumCorrection := int64(MAGIC_BABA) - int64(checksum)
		if checksumCorrection < 0 {
			checksumCorrection += 256 * 256
		}
		checksum = uint16(checksumCorrection)

		checksumBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(checksumBytes, checksum)
		file[126] = checksumBytes[0]
		file[127] = checksumBytes[1]

		err = os.WriteFile(fmt.Sprintf("%s-%s.bin", outputPrefix, macStartAddr.String()), file, 0644)
		if err != nil {
			fmt.Printf("unable to write output file: %s\n", err.Error())
			os.Exit(-1)
		}
		break
	}
}

func splitAny(sep rune) bool {
	return !unicode.IsLetter(sep) && !unicode.IsNumber(sep)
}

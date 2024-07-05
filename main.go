package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// This tool changes the MAC address in an Intel Foxville NVM file and recalculates
// the checksum for a particular mac address or range

func main() {

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

	if inputFile == "" || macStartAddr == nil {
		fmt.Println("-i and -ms are mandatory parameters to run")
		fmt.Printf("Usage: \t ./intel-nvm-tool \\\n\t-i <input stock nvm file> \\\n")
		fmt.Printf("\t-ms <first mac address> \\\n\t[-me <last mac address>] \\\n")
		fmt.Printf("\t[-o <output nvm file prefix>]\n")
		os.Exit(-1)
	}

	if outputPrefix == "" {
		outputPrefix = "foxville-nvm"
	}
	file, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("unable to open input file: %s\n", err.Error())
		os.Exit(-1)
	}

	macStart := binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(macStartAddr)...))
	macEnd := macStart
	if len(macEndAddr) == 6 {
		macEnd = binary.BigEndian.Uint64(append([]byte{0, 0}, []byte(macEndAddr)...))
	} else {
		macEndAddr = macStartAddr
	}

	fmt.Printf("MAC address range: %s - %s (%d address(es))\n",
		macStartAddr.String(),
		macEndAddr.String(),
		macEnd-macStart+1,
	)

	macBuffer := make([]byte, 8)
	for x := macStart; x <= macEnd; x++ {
		binary.BigEndian.PutUint64(macBuffer, x)
		file = append(macBuffer[2:8], file[6:]...)

		var checksum uint16
		for y := 0; y < 0x7d; y += 2 {
			checksum += binary.LittleEndian.Uint16(file[y : y+2])
		}
		checksumCorrection := int64(MAGIC_BABA) - int64(checksum)
		if checksumCorrection < 0 {
			checksumCorrection += 256 * 256
		}
		checksum = uint16(checksumCorrection)

		checksumBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(checksumBytes, checksum)
		file[126] = checksumBytes[0]
		file[127] = checksumBytes[1]

		err = os.WriteFile(fmt.Sprintf("%s-%s.bin",
			outputPrefix,
			fmt.Sprintf("%02X-%02X-%02X-%02X-%02X-%02X",
				macBuffer[2], macBuffer[3], macBuffer[4],
				macBuffer[5], macBuffer[6], macBuffer[7])),
			file,
			0644,
		)
		if err != nil {
			fmt.Printf("unable to write output file: %s\n", err.Error())
			os.Exit(-1)
		}
	}
}

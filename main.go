package i225_226_firmware_cutomisation_tool

import (
	"fmt"
	"os"
	"strings"
)

// main is the entry point of the program.
// It accepts command line arguments and sorts a slice of integers.
// The sorted elements are printed to the standard output using fmt.Println.
func main() {

	// ./intel-nvm-tool \
	//		-i <input stock nvm file> \
	//		-m <mac address / mac address range> \
	//		-o <output nvm file prefix>

	var inputFile, outputPrefix, macAddressRange string
	args := os.Args[1:]

	for x := 0; x < len(args); x++ {

		if x < len(args)-2 {
			break
		}

		switch args[x] {
		case "-i":
			inputFile = args[x+1]
		case "-o":
			outputPrefix = args[x+1]
		case "-m":
			macAddressRange = args[x+1]
		}
	}

	if inputFile == "" || outputPrefix == "" || macAddressRange == "" {
		println("-i -o and -m are mandatory parameters to run")
		os.Exit(-1)
	}

	file, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("unable to open input file: %s\n", err.Error())
		os.Exit(-1)
	}

	macAddresses := strings.Split(macAddressRange, ":")

}

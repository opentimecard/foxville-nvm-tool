### Command line tool to create Intel Foxville NVM files 

Tool to add mac address to an Intel Foxville NVM image and recalculate the checksum. Supports a single mac or a range for bulk mode

To install :
`go install github.com/opentimecard/foxville-nvm-tool@latest`

Usage :

```
./foxville-nvm-tool \
  -i <input nvm file> \
  -ms <first mac address> \
 [-me <last mac address>] \
 [-o <output nvm file prefix>]
```
 
Example :
`./foxville-nvm-tool -i FXVL_125B_LM_2MB_2.25.bin -ms 8C-1F-64-10-41-01 -me 8C-1F-64-10-41-C8 -o timebeat-mac-address`


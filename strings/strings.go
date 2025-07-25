package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"flag"
	"os"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "Input file with newline-separated strings")
	outputFile := flag.String("output", "obfuscated.go", "Output Go file with hashmap and decode function")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input file with -input flag")
		return
	}

	key, err := generateKey()
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}
	strings, err := readStringsFromFile(*inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	err = generateGoCode(strings, key, *outputFile)
	if err != nil {
		fmt.Println("Error generating Go code:", err)
		return
	}
	fmt.Printf("Successfully generated %s with %d obfuscated strings\n", *outputFile, len(strings))
}

func generateKey() ([]byte, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

func readStringsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { 
			result = append(result, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return result, nil
}

func obfuscateString(input string) string {
	data := []byte(input)
	obfuscated := make([]byte, len(data))
	
	for i, b := range data {
		transformed := b
		transformed = ((transformed << 3) | (transformed >> 5))
		transformed ^= 0xAA                                      
		transformed = ^transformed                              
		transformed += byte(i + 1)                             
		obfuscated[i] = transformed
	}
	
	return hex.EncodeToString(obfuscated)
}


func deobfuscateString(obfuscated string) (string, error) {
	data, err := hex.DecodeString(obfuscated)
	if err != nil {
		return "", err
	}
	
	original := make([]byte, len(data))
	for i, b := range data {
		transformed := b
		transformed -= byte(i + 1)            
		transformed = ^transformed                              
		transformed ^= 0xAA                                   
		transformed = ((transformed >> 3) | (transformed << 5)) 
		original[i] = transformed
	}
	
	return string(original), nil
}

func generateGoCode(strings []string, key []byte, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()
	fmt.Fprintf(file, "package obf\n\n")
	fmt.Fprintf(file, "import (\n")
	fmt.Fprintf(file, "\t\"encoding/hex\"\n")
	fmt.Fprintf(file, "\t\"fmt\"\n")
	fmt.Fprintf(file, ")\n\n")
	fmt.Fprintf(file, "var obfuscatedStrings = []string{\n")
	for _, str := range strings {
		obfuscated := obfuscateString(str)
		fmt.Fprintf(file, "\t%q,\n", obfuscated)
	}
	fmt.Fprintf(file, "}\n\n")
	fmt.Fprintf(file, "func Decode(obfuscated string) (string, error) {\n")
	fmt.Fprintf(file, "\tdata, err := hex.DecodeString(obfuscated)\n")
	fmt.Fprintf(file, "\tif err != nil {\n")
	fmt.Fprintf(file, "\t\treturn \"\", err\n")
	fmt.Fprintf(file, "\t}\n\n")
	fmt.Fprintf(file, "\toriginal := make([]byte, len(data))\n")
	fmt.Fprintf(file, "\tfor i, b := range data {\n")
	fmt.Fprintf(file, "\t\t// Reverse the transformations in opposite order\n")
	fmt.Fprintf(file, "\t\ttransformed := b\n")
	fmt.Fprintf(file, "\t\ttransformed -= byte(i + 1)                               // Subtract position offset\n")
	fmt.Fprintf(file, "\t\ttransformed = ^transformed                               // Bitwise NOT\n")
	fmt.Fprintf(file, "\t\ttransformed ^= 0xAA                                      // XOR with constant\n")
	fmt.Fprintf(file, "\t\ttransformed = ((transformed >> 3) | (transformed << 5)) // Rotate right 3 bits\n")
	fmt.Fprintf(file, "\t\toriginal[i] = transformed\n")
	fmt.Fprintf(file, "\t}\n\n")
	fmt.Fprintf(file, "\treturn string(original), nil\n")
	fmt.Fprintf(file, "}\n\n")
	fmt.Fprintf(file, "// Get returns the decoded string at the given index\n")
	fmt.Fprintf(file, "func Get(index int) (string, error) {\n")
	fmt.Fprintf(file, "\tif index < 0 || index >= len(obfuscatedStrings) {\n")
	fmt.Fprintf(file, "\t\treturn \"\", fmt.Errorf(\"index %%d out of range [0, %%d)\", index, len(obfuscatedStrings))\n")
	fmt.Fprintf(file, "\t}\n")
	fmt.Fprintf(file, "\treturn Decode(obfuscatedStrings[index])\n")
	fmt.Fprintf(file, "}\n\n")
	fmt.Fprintf(file, "// Example usage\n")
	fmt.Fprintf(file, "func main() {\n")
	fmt.Fprintf(file, "\t// Example: Get decoded string by index\n")
	fmt.Fprintf(file, "\tif decoded, err := Get(0); err == nil {\n")
	fmt.Fprintf(file, "\t\tfmt.Printf(\"Decoded string at index 0: %%s\\n\", decoded)\n")
	fmt.Fprintf(file, "\t} else {\n")
	fmt.Fprintf(file, "\t\tfmt.Printf(\"Error getting string at index 0: %%v\\n\", err)\n")
	fmt.Fprintf(file, "\t}\n\n")
	fmt.Fprintf(file, "\t// Show all decoded strings\n")
	fmt.Fprintf(file, "\tfmt.Println(\"\\nAll strings decoded from obfuscated data:\")\n")
	fmt.Fprintf(file, "\tfor i := range obfuscatedStrings {\n")
	fmt.Fprintf(file, "\t\tif decoded, err := Get(i); err == nil {\n")
	fmt.Fprintf(file, "\t\t\tfmt.Printf(\"[%%d]: %%s\\n\", i, decoded)\n")
	fmt.Fprintf(file, "\t\t} else {\n")
	fmt.Fprintf(file, "\t\t\tfmt.Printf(\"[%%d]: Error - %%v\\n\", i, err)\n")
	fmt.Fprintf(file, "\t\t}\n")
	fmt.Fprintf(file, "\t}\n")
	fmt.Fprintf(file, "}\n")
	return nil
}
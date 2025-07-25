package enc

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"lawk/pkg/walk"
	"github.com/tyler-smith/go-bip39"
)


func DecryptFiles() error {
	masterKey, err := generateMasterKey()
	if err != nil {
		return fmt.Errorf("failed to generate master key: %w", err)
	}
	files, err := walk.Walk()
	if err != nil {
		return fmt.Errorf("failed to walk directories: %w", err)
	}
	for i, filePath := range files {
		fmt.Printf("[%d/%d] Decrypting: %s\n", i+1, len(files), filePath)
		err := decryptFileFromMnemonic(filePath, masterKey)
		if err != nil {
			log.Printf("Failed to decrypt %s: %v", filePath, err)
			continue
		}
	}
	return nil
}

func decryptFileFromMnemonic(filePath string, masterKey []byte) error {
	encryptedData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	if len(encryptedData) == 0 {
		log.Printf("Skipping empty file: %s", filePath)
		return nil
	}
	content := string(encryptedData)
	if !strings.HasPrefix(content, "# mnemonic seed phrase encryption test :3") {

		log.Printf("Skipping non-encrypted file: %s", filePath)
		return nil
	}
	log.Printf("Processing encrypted file: %s", filePath)
	hasher := sha256.New()
	hasher.Write(masterKey)
	hasher.Write([]byte(filePath))
	fileSeed := hasher.Sum(nil)
	originalData, err := convertFromMnemonic(content, fileSeed)
	if err != nil {
		return fmt.Errorf("failed to convert from mnemonic: %w", err)
	}
	log.Printf("Successfully decrypted %d bytes for file: %s", len(originalData), filePath)
	err = ioutil.WriteFile(filePath, originalData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}
	return nil
}

func convertFromMnemonic(mnemonicContent string, seed []byte) ([]byte, error) {
	lines := strings.Split(mnemonicContent, "\n")
	var mnemonicLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			mnemonicLines = append(mnemonicLines, line)
		}
	}
	if len(mnemonicLines) == 0 {
		return nil, fmt.Errorf("no mnemonic phrases found in file")
	}
	log.Printf("Found %d mnemonic lines to process", len(mnemonicLines))
	var originalData []byte
	for chunkIndex, mnemonicPhrase := range mnemonicLines {
		entropy, err := bip39.EntropyFromMnemonic(mnemonicPhrase)
		if err != nil {
			return nil, fmt.Errorf("failed to convert mnemonic to entropy: %w", err)
		}
		if len(entropy) != 16 {
			return nil, fmt.Errorf("unexpected entropy length: %d, expected 16", len(entropy))
		}
		hasher := sha256.New()
		hasher.Write(seed)
		hasher.Write([]byte{byte(chunkIndex)})
		keyStream := hasher.Sum(nil)

	
		paddedChunk := make([]byte, 16)
		for j := range paddedChunk {
			paddedChunk[j] = entropy[j] ^ keyStream[j]
		}
		var originalChunk []byte
		if chunkIndex == len(mnemonicLines)-1 {
			lastByte := paddedChunk[15]
			if lastByte > 0 && lastByte < 16 {
				originalChunk = paddedChunk[:lastByte]
			} else {
				originalChunk = paddedChunk
			}
		} else {
			originalChunk = paddedChunk
		}
		originalData = append(originalData, originalChunk...)
	}
	return originalData, nil
}

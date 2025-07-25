# lawk - mnemonic encryption ransomware simulation

a proof-of-concept ransomware that encrypts files by converting them into bip39 mnemonic seed phrases, designed to explore entropy reduction during encryption operations.

## concept

traditional ransomware creates high-entropy encrypted data that looks obviously malicious to detection systems. this project takes a different approach, converting encrypted data into human-readable mnemonic phrases that appear as legitimate cryptocurrency seed phrases rather than suspicious encrypted blobs.

## technical implementation

### encryption process

1. **deterministic key generation**: creates a master key from hostname + os info + salt using sha256
2. **file-specific seeding**: each file gets a unique seed derived from master key + file path
3. **chunked processing**: splits file data into 16-byte chunks for processing
4. **xor encryption**: each chunk is xored with a deterministic keystream
5. **mnemonic conversion**: encrypted chunks become bip39 mnemonic phrases (12 words per chunk)


### decryption process

1. **header detection**: identifies encrypted files by custom header
2. **mnemonic parsing**: extracts seed phrases from file content
3. **entropy reconstruction**: converts mnemonics back to 16-byte entropy values
4. **keystream regeneration**: recreates the same deterministic keys used for encryption
5. **xor reversal**: entropy xor keystream = original data
6. **chunk reassembly**: handles padding and reconstructs original file content

### why this approach is interesting

**entropy reduction**: instead of high-entropy random-looking data, encrypted files contain structured english words that:
- pass basic entropy checks
- look like legitimate cryptocurrency content
- blend in with normal user data (crypto wallets, trading notes, etc.)

**detection evasion**: security tools looking for entropy spikes or suspicious file patterns might miss files that appear to contain wallet mnemonics

## usage

```bash
# encrypt files
cd cmd
go build -o lawk.exe
./lawk.exe

# decrypt files  
./lawk.exe -decrypt
```

## technical notes

- uses bip39 standard for mnemonic generation (compatible with real crypto tools)
- deterministic encryption allows decryption without key storage
- concurrent file processing for performance
- handles padding for variable-length chunks
- preserves file permissions and structure
- minimal string obfuscation applied to target extensions and directories
## disclaimer

this is educational/research code demonstrating novel encryption concealment techniques. not intended for malicious use.


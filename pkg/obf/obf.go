package obf

import (
	"encoding/hex"
	"fmt"
)

var obfuscatedStrings = []string{
	"7880d112fb34dd",
	"7830510243842dfed7",
	"7830f12a3b34657ed7",
	"3028817bcb24ed86",
	"703031723b845b7fcf28f18a",
	"1f513932037c5b7fcf28f18a",
	"78c831da4b349d",
	"25d0d11a",
	"25608902cb84",
	"6f20f91a034c",
	"6f20f93a634c",
	"25f899fa",
	"25d8796a",
	"254079",
	"257831529b",
	"25e8419a",
	"25e8417a13",
	"259839d2",
	"253079da",
	"25f0d97a",
	"25f0d9d2",
	"2550d1ea",
	"25d8617283d4",
	"252801424b84cdd6",
	"25383172",
	"25506972",
	"25d8699a",
	"25d88142",
	"251081a2",
	"25d8dfc8",
	"255081ca",
	"2550312a6b",
	"2508d1322b",
	"258029ea",
	"2518f9da63d4d5f67f",
	"2560f1d2",
	"2550c9827b842dfe276841da",
}

func Decode(obfuscated string) (string, error) {
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

func Get(index int) (string, error) {
	if index < 0 || index >= len(obfuscatedStrings) {
		return "", fmt.Errorf("index %d out of range [0, %d)", index, len(obfuscatedStrings))
	}
	return Decode(obfuscatedStrings[index])
}


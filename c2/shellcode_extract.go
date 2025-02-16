package main

import (
	"debug/pe"
	"encoding/hex"
	"fmt"
	"os"
)

func extractShellCode(imagePath string) {
	file, err := pe.Open(imagePath)
	if err != nil {
		panic(err)

	}
	defer file.Close()

	var textSection *pe.Section

	for _, section := range file.Sections {
		if section.Name == ".text" {
			textSection = section
			break
		}
	}

	if textSection == nil {
		fmt.Println("Text section not found")
		return
	}

	textData := make([]byte, textSection.Size)
	_, err = textSection.ReadAt(textData, 0)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println(hex.Dump(textData))

	err = os.WriteFile("shellcode.bin", textData, 0644)
	if err != nil {
		fmt.Println("error while writting in shellcode.bin:", err)
		return
	}

}

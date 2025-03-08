package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func main() {
	payloadType := flag.String("t", "rev-https", "Payload type")
	remoteHost := flag.String("H", "localhost", "Remote host")
	remotePort := flag.Int("p", 443, "Remote port")
	optFN := flag.String("o", "client.exe", "Output file name")
	format := flag.String("f", "exe", "Output file format")
	isConfFile := flag.Bool("c", false, "Enable config file")

	flag.Parse()

	fmt.Println("\033[1m\033[34m[>] Payload type:\033[0m", *payloadType)
	fmt.Println("\033[1m\033[34m[>] Remote Host:\033[0m", *remoteHost)
	fmt.Println("\033[1m\033[34m[>] Remote Port:\033[0m", *remotePort)
	fmt.Println("\033[1m\033[34m[>] Output File Name:\033[0m", *optFN)
	fmt.Println("\033[1m\033[34m[>] Output File Format:\033[0m", *format)
	fmt.Println("\033[1m\033[34m[>] Is Conf. File Used:\033[0m", *isConfFile)

	err := os.MkdirAll("./temp/client", os.ModePerm)
	if err != nil {
		fmt.Println("\033[1m\033[31m[-] Failed to create directory:\033[0m", err)
		return
	}

	source := "../client"
	destination := "./temp/client"

	fmt.Println("\033[1m\033[33m[$] Copying directory from\033[0m", source, "\033[1m\033[33mto\033[0m", destination)

	err = copyDir(source, destination)
	if err != nil {
		fmt.Println("\033[1m\033[31m[-] Error while copying directory:\033[0m", err)
		return
	}

	fmt.Println("\033[1m\033[32m[+] Directory successfully copied to\033[0m", destination)
	
}

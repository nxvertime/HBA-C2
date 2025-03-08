package main

import (
	"flag"
	"fmt"
)

func main() {
	payloadType := flag.String("t", "rev-https", "Payload type")
	remoteHost := flag.String("H", "localhost", "Remote host")
	remotePort := flag.Int("p", 443, "Remote port")
	optFN := flag.String("o", "client.exe", "Output file name")
	isConfFile := flag.Bool("c", false, "Enable config file")
	flag.Parse()
	fmt.Println("[>] Payload type: ", *payloadType)
	fmt.Println("[>] Remote Host: ", *remoteHost)
	fmt.Println("[>] Remote Port: ", *remotePort)
	fmt.Println("[>] Output File Name: ", *optFN)
	fmt.Println("[>] Is Conf. File Used: ", *isConfFile)

}

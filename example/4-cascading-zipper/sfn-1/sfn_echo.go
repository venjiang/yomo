package main

import (
	"log"
	"os"

	"github.com/yomorun/yomo"
)

func main() {
	sfn := yomo.NewStreamFunction(
		"echo-sfn",
		yomo.WithZipperAddr("localhost:9001"),
		yomo.WithObserveDataTags(0x33),
		yomo.WithCredential("token:z1"),
	)
	defer sfn.Close()

	// set handler
	sfn.SetHandler(handler)

	// start
	err := sfn.Connect()
	if err != nil {
		log.Fatalf("[sfn] connect err=%v", err)
		os.Exit(1)
	}

	select {}
}

func handler(data []byte) (byte, []byte) {
	val := string(data)
	log.Printf(">> [sfn] got tag=0x33, data=%s", val)
	val="🍌" + val
	return 0x34, []byte(val)
}
package main

import (
	"log"
	"os"

	"github.com/yomorun/yomo"
)

func main() {
	zipper := yomo.NewZipperWithOptions(
		"zipper-3",
		yomo.WithZipperAddr("localhost:9003"),
		// yomo.WithAuth("token", "z2"),
	)
	defer zipper.Close()

	zipper.ConfigWorkflow("zipper_3_wf.yaml")

	// add Downstream Zipper
	zipper.AddDownstreamZipper(yomo.NewDownstreamZipper(
		"zipper-1",
		yomo.WithZipperAddr("localhost:9001"),
		yomo.WithCredential("token:z1"),
	))
	zipper.AddDownstreamZipper(yomo.NewDownstreamZipper(
		"zipper-2",
		yomo.WithZipperAddr("localhost:9002"),
		yomo.WithCredential("token:z2"),
	))
	// start zipper service
	log.Printf("Server has started!, pid: %d", os.Getpid())
	go func() {
		err := zipper.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	select {}
}

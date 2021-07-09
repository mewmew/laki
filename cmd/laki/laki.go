package main

import (
	"log"

	"github.com/mewmew/laki/vk"
)

func main() {
	if err := vk.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
}

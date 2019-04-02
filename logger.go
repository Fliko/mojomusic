package main

import (
	"io"
	"log"
	"os"
)

func logger(err error, msg string) {
	f, err := os.OpenFile("orders.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stderr, f)
	log.SetOutput(wrt)

	//	t := time.Now()

	log.Println(err, msg)
}

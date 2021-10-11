package main

import (
	"log"
	"smitClient/server"
)

func main()  {

	h, err := server.New(":8088", 44100)

	if err != nil {
		log.Fatalln(err)
	}

	h.Run()

}

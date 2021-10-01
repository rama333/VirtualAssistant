package main

import server "VirtualAssistant/server"

func main()  {
	s := server.NewServer(":8081", 44100)

	s.Start()

}

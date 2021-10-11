package server

import (
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)
var (
	sampleRate = 44100
	seconds = 5
)


type Server struct {
	Address string
	SampleRate int
}


func New(address string, sampleRate int) (*Hub, error)  {

	flag.Parse()
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})


	go func() {
		//defer server.wg.Done()
		for {
			logrus.Info("start ws")
			if err := http.ListenAndServe(address, nil); err != nil{
				if err == http.ErrServerClosed {
					return
				}
				logrus.Info("failed to start")
			}

			time.Sleep(3 * time.Second)

			logrus.Info("restarting ws")
		}
	}()


	return hub, nil

}

//func (s * Server) Start() {
//
//	fmt.Println("Launching server...")
//	numBatch := 0
//	buffer := make([]float32, 0)
//
//	ln, _ := net.Listen("tcp", s.Address)
//
//	conn, _ := ln.Accept()
//
//
//
//
//}

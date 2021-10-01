package server

import (
	"bufio"
	"fmt"
	"github.com/cryptix/wav"
	"github.com/wcharczuk/go-chart/v2"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)
var (
	sampleRate = 44100
	seconds = 5
)


type Server struct {
	Address string
	SampleRate int
}

func NewServer(address string, sampleRate int) (*Server)  {

	return NewServer(address, sampleRate)
}

func (s * Server) Start() {

	fmt.Println("Launching server...")
	numBatch := 0
	buffer := make([]float32, 0)

	ln, _ := net.Listen("tcp", s.Address)

	conn, _ := ln.Accept()


	for {
		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			log.Panic(err)
		}

		//fmt.Print("Message Received:", string(message))
		// Процесс выборки для полученной строки
		//conn.Write([]byte(ok + "\n"))


		buf := strings.Split(message[0:len(message)-2], " ")

		for i := 0; i < len(buf); i++ {
			log.Println(buf[i])
			fl, err := strconv.ParseFloat(buf[i], 10)

			if err != nil {
				log.Println(buf[i])
				log.Println(err, buf[i])
			}
			buffer = append(buffer, float32(fl))
		}

		log.Println(len(buffer))


		audioBytes := make([]byte, len(buffer)*4)

		for i, n := 0, len(buffer); i < n; i++ {
			buffer[i] *= 32767

			audioBytes[i*2] = byte(buffer[i])
			t := int(buffer[i]) >> 8
			audioBytes[i*2+1] = byte(t)
		}
		//
		//log.Println(audioBytes)

		//waveFile, err := os.Create("wc3.wav")
		//
		//param := wave.WriterParam{
		//	Out:           waveFile,
		//	Channel:       1,
		//	SampleRate:    sampleRate,
		//	BitsPerSample: 8, // if 16, change to WriteSample16()
		//}
		//
		//
		//waveWriter, err := wave.NewWriter(param)
		//
		//
		//_, err = waveWriter.WriteSample8(nil) // WriteSample16 for 16 bits
		//
		//if err != nil{
		//	panic(err)
		//}

		x := make([]float64, len(buffer))
		y := make([]float64, len(buffer))

		b := make([]byte, len(buffer))
		//i := make([]int32, sampleRate * seconds)

		for i := 0; i < len(buffer); i++ {
			x[i] = float64(i)
			y[i] = float64(buffer[i])
			b[i] = byte(buffer[i])
		}

		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: x,
					YValues: y,
				},
			},
		}

		f, err := os.Create(filepath.Join("./data", filepath.Base(fmt.Sprintf("%d", numBatch)+".png")))

		defer f.Close()

		graph.Render(chart.PNG, f)

		////////////////
		f, err = os.Create(filepath.Join("./data", filepath.Base("output"+fmt.Sprintf("%d", numBatch)+".wav")))


		kk := len(buffer)
		// Create the headers for our new mono file
		meta := wav.File{
			Channels:        1,
			SampleRate:      uint32(kk),
			SignificantBits: 16,
		}

		writer, err := meta.NewWriter(f)
		if err != nil {
			log.Println(err)
		}

		// Write to file
		//	for k,_ := range b {
		//log.Println(i,t)
		//log.Println(i[k])
		k, err := writer.Write(audioBytes)
		log.Println(k)
		//checkErr(err)

		err = writer.Close()

		numBatch += 1
	}

}

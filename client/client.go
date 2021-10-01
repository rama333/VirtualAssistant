package main

import (
	"encoding/binary"
	"fmt"
	"github.com/cryptix/wav"
	"github.com/gordonklaus/portaudio"
	"github.com/wcharczuk/go-chart/v2"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const sampleRate = 44100
const seconds = 5

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8081")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	//stop1 := make(chan struct{})
	//
	//
	//
	//
	//go func() {
	//	for  {
	//
	//		fmt.Println("start send")
	//		fmt.Fprintf(conn, fmt.Sprintf("%s", "test" + "\n"))
	//		fmt.Println("stop send")
	//
	//		time.Sleep(10 *time.Second)
	//		stop1<- struct{}{}
	//	}
	//}()
	//
	//<-stop1
	//
	//return

	buffer := make([]float32, sampleRate*seconds)
	buffer1 := make([]float32, 0)
	//
	//
	//buffer = append(buffer, 5.5)
	//
	//t := fmt.Sprintf("%v", buffer)
	//t1 := t[1:len(t)-1]
	//
	//
	//
	//log.Println(t1)

	//fmt.Fprintf(conn, fmt.Sprintf("%v", buffer) + "\n")

	//conn, err := net.Dial("tcp", "127.0.0.1:8081")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	//
	//sig := make(chan os.Signal, 1)
	//signal.Notify(sig, os.Interrupt, os.Kill)

	portaudio.Initialize()
	defer portaudio.Terminate()
	//buffer := make([]float32, sampleRate * seconds)

	streamChan := make(chan float32)
	stop := make(chan struct{})
	//fmt.Fprintf(conn, fmt.Sprintf("%s", "t1"))

	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), func(in []float32) {
		//resp, err := http.Get("http://localhost:8080/audio")
		//chk(err)
		//body, _ := ioutil.ReadAll(resp.Body)
		//responseReader := bytes.NewReader(body)
		//binary.Read(responseReader, binary.BigEndian, &buffer)
		for i := range buffer {

			streamChan <- in[i]
			//log.Println("<-")

			// Отправляем в socket
			//fmt.Fprintf(conn, fmt.Sprintf("%s", buffer[i]) + "\n")
			//
			//bufio.NewReader(conn).ReadString('\n')
			//fmt.Print("Message from server: ", message)

		}
	})

	go func() {
		for {

			h :=  <-streamChan

			if h != 0 {

				buffer1 = append(buffer1, h)

			}


			if len(buffer1) == 88200 {

				log.Println("ok")

				t := fmt.Sprintf("%v", buffer1)
				t1 := t[1 : len(t)-1]
				fmt.Fprintf(conn, fmt.Sprintf("%s", t1) + "\n")

				//log.Println(t1)

				//	stop <- struct{}{}

				log.Println("before", len(buffer1))
				buffer1 = buffer1[:0]
				log.Println("after", len(buffer1))
			}

		}

	}()

	log.Println("stream start")
	chk(stream.Start())
	time.Sleep(time.Second * 5)
	chk(stream.Stop())
	log.Println("stream stop")

	//	log.Println(buffer)

	defer stream.Close()

	if err != nil {
		fmt.Println(err)
	}

	/////////////

	//portaudio.Initialize()
	//defer portaudio.Terminate()
	////var audio io.Reader
	////out := make([]int32, 8192)
	//stream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, len(buffer), &buffer)
	//chk(err)
	//defer stream.Close()
	//
	//chk(stream.Start())
	//
	//
	//
	//
	//stream.Write()
	////time.Sleep(time.Second * 4)
	//defer stream.Stop()
	//
	audioBytes := make([]byte, sampleRate*seconds*3)

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

	x := make([]float64, sampleRate*seconds)
	y := make([]float64, sampleRate*seconds)

	b := make([]byte, sampleRate*seconds)
	//i := make([]int32, sampleRate * seconds)

	for i := 0; i < sampleRate*seconds; i++ {
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
	f, _ := os.Create("output2.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
	////////////////
	f, err = os.Create("vdv.wav")

	// Create the headers for our new mono file
	meta := wav.File{
		Channels:        1,
		SampleRate:      sampleRate,
		SignificantBits: 16,
	}

	writer, err := meta.NewWriter(f)
	chk(err)

	// Write to file
	//	for k,_ := range b {
	//log.Println(i,t)
	//log.Println(i[k])
	k, err := writer.Write(audioBytes)
	log.Println(k)

	log.Println(len(buffer1))
	//checkErr(err)

	err = writer.Close()
	//checkErr(err)

	<-stop

}
func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

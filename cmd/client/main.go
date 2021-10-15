package main

import (
	"encoding/json"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/zenwerk/go-wave"
	"os"
	"path/filepath"
	"smitClient/synthesize"
	"sync"
	"time"
)

type RecognitionText struct {
	CurrentText string `json:"current_text"`
	FullText string `json:"full_text"`
}

func main() {

	wg := sync.WaitGroup{}

	wg.Add(1)

	waveFile, _ := os.Create(filepath.Join("./../../data", filepath.Base("testoutput"+fmt.Sprintf("%d", "5")+".wav")))

	var param = wave.WriterParam{
		Out:           waveFile,
		Channel:       1,
		SampleRate:    44100,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	waveWriter, err := wave.NewWriter(param)

	portaudio.Initialize()
	defer portaudio.Terminate()
	c, err := newClient(time.Second / 3)

	if err != nil {
		logrus.Fatal(err)
	}

	defer c.Close()

	err = c.Start()

	if err != nil {
		logrus.Error(err)
	}

	go func() {

		var s1 string

		fmt.Scanln(&s1)

		if s1 == "" {
			waveWriter.Close()
			os.Exit(0)
		}

	}()

	c.processAudio()

	//time.Sleep(4 * time.Second)
	defer waveWriter.Close()

	err = c.Stop()

	if err != nil {
		logrus.Error(err)
	}

	defer c.Stream.Close()

	//c.sendAudio()

	//e.Stream.Close()

	wg.Wait()
}

type Client struct {
	connect *websocket.Conn
	*portaudio.Stream
	buffer []byte
	i      int
}

type Audio struct {
	Data []byte `json:"data"`
}

func newConnect(URL string) (connect *websocket.Conn, err error) {
	connect, _, err = websocket.DefaultDialer.Dial(URL, nil)

	return
}

func newClient(delay time.Duration) (*Client, error) {

	con, err := newConnect("ws://192.168.114.145:8088/")

	if err != nil {
		return nil, err
	}

	framesPerBuffer := make([]byte, 64)

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(44100), len(framesPerBuffer), framesPerBuffer)

	e := &Client{
		buffer:  make([]byte, 64),
		connect: con,
	}

	go func(client *Client) {

		for {

			speech := synthesize.Speech{Folder: "", Language: synthesize.Russia, Synthesize: synthesize.MPlayer{}}

			var recData RecognitionText

			_, mes, err := client.connect.ReadMessage()

			if err != nil {
				logrus.Info(err)
			}

			err = json.Unmarshal(mes, &recData)

			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("current data: ", recData.CurrentText)
			logrus.Info("full text: ", recData.FullText)

			err = speech.Speak(recData.CurrentText)

			if err != nil {
				logrus.Error(err)
			}
		}
	}(e)

	e.Stream = stream
	e.buffer = framesPerBuffer

	//	e.Stream, err = portaudio.OpenStream(p, e.processAudio)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func (c *Client) processAudio() {

	a := 0

	for {
		a += 1

		err := c.Stream.Read()

		if err != nil {
			logrus.Error(err)
		}

		c.sendAudio()

		//if a > 5500 {
		//	break
		//}

		//_, err = waveWriter.Write(c.buffer) // WriteSample16 for 16 bits
		//if err != nil {
		//	logrus.Error(err)
		//}

		//logrus.Info(a)

	}

}

func (c *Client) sendAudio() {

	var audio Audio

	audio.Data = c.buffer

	data, err := json.Marshal(audio)

	if err != nil {
		logrus.Error(err)
	}

	err = c.connect.WriteMessage(websocket.TextMessage, data)

	//if err != nil {
	//	logrus.Error(err)
	//}

	//waveFile, _ := os.Create(filepath.Join("./../../data", filepath.Base("testoutput"+fmt.Sprintf("%d", "5")+".wav")))
	//
	//var param = wave.WriterParam{
	//	Out:           waveFile,
	//	Channel:       1,
	//	SampleRate:    44100,
	//	BitsPerSample: 8, // if 16, change to WriteSample16()
	//}
	//
	//waveWriter, err := wave.NewWriter(param)
	//
	//_, err = waveWriter.Write([]byte(c.buffer)) // WriteSample16 for 16 bits
	//if err != nil {
	//	logrus.Error(err)
	//}
	//
	//waveWriter.Close()

	//var audio Audio
	//
	//audio.Data = c.buffer
	//
	//logrus.Info("audio",audio)
	//
	////e.connect.WriteJSON(audio)
	//
	//
	//data, err := json.Marshal(audio)
	//
	//if err != nil {
	//	logrus.Error(err)
	//}
	//
	//err = c.connect.WriteMessage(websocket.TextMessage, data)
	//
	//if err != nil {
	//	logrus.Info(err)
	//}

	//err = c.connect.WriteMessage(websocket.TextMessage, []byte("test"))
	//if err != nil {
	//	logrus.Info(err)
	//}
}

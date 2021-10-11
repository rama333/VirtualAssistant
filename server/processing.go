package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zenwerk/go-wave"
	"os"
	"path/filepath"
	"regexp"
	"smitClient/recognition/kaldi_go"
	"smitClient/server/decoder"
	"sync"

	//	wave "github.com/zenwerk/go-wave"
)


type Audio struct {
	Data []byte `json:"data"`
}

type Processing struct {
	buffer [][]byte
	kaldi *kaldi_go.Encoder
	mutex sync.RWMutex
	fullText []string
}


func NewProcessing() (*Processing) {

	kaldi := kaldi_go.NewConfig(&kaldi_go.Config{
		"/home/ubuntu/kaldi/src/online2bin/online2-wav-nnet3-latgen-faster",
		false,
		3,
		1.0,
		13.0,
		6.0,
		7000,
		"/home/ubuntu/speech/kaldi-ru-0.9/exp/tdnn/conf/online.conf",
		"/home/ubuntu/speech/kaldi-ru-0.9/data/lang_test_rescore/words.txt",
		"/home/ubuntu/speech/kaldi-ru-0.9/exp/tdnn/final.mdl",
		"/home/ubuntu/speech/kaldi-ru-0.9/exp/tdnn/graph/HCLG.fst",
	})

	return &Processing{
		kaldi: kaldi,
		mutex: sync.RWMutex{},
		buffer: make([][]byte, 0, 100),
		fullText: make([]string, 0, 100),
	}
}

func(p *Processing) process(stream []byte) (string, error) {

	buffer := make([]byte, 0)

	var audio Audio
	err := json.Unmarshal(stream, &audio)

	if err != nil {
		logrus.Error(err)
	}

	buffer = audio.Data[:]

	//log.Println(len(buffer))
	p.buffer = append(p.buffer, buffer)


	if len(p.buffer) >= 1500 {

		fileName := uuid.New()
		waveFile, _ := os.Create(filepath.Join("/home/ubuntu/smit/data/", filepath.Base( fileName.String() +".wav")))

		var param = wave.WriterParam{
			Out:           waveFile,
			Channel:       1,
			SampleRate:    44100,
			BitsPerSample: 8, // if 16, change to WriteSample16()
		}

		waveWriter, err := wave.NewWriter(param)

		if err != nil {
			logrus.Error(err)
		}

		for _, byt := range p.buffer {
			waveWriter.Write(byt)
		}

		waveWriter.Close()

		p.mutex.Lock()

		defer p.mutex.Unlock()

		dec := decoder.NewDecoder()

		err = dec.Dec(filepath.Join("/home/ubuntu/smit/data/", filepath.Base( fileName.String() +".wav")), filepath.Join("/home/ubuntu/smit/data/", filepath.Base( fileName.String() +"dec.wav")))
		if err != nil{
			logrus.Panic(err)
		}

		text, err := p.kaldi.Recognition(filepath.Join("/home/ubuntu/smit/data/", filepath.Base( fileName.String() +"dec.wav")))

		if err != nil {
			logrus.Error(err)
		}

		p.buffer = p.buffer[:0]

		re := regexp.MustCompile(`\r?\n`)
		text = re.ReplaceAllString(text, "")

		if text != "" || text  != " " {
			p.fullText = append(p.fullText, text)
		}


		logrus.Printf("%v", text)

		return text, nil

	}


	return "ok", nil

}

package synthesize

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Speech struct {
	Folder     string
	Language   string
	Synthesize Synthesize
}

func (s *Speech) Speak(text string) (err error) {

	f, err := s.synthesizeSpeech(text)

	if err != nil {
		return err
	}

	err = s.Synthesize.Play(f)

	return
}

func (s *Speech) synthesizeSpeech(text string) (string, error) {

	URL := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), s.Language)

	res, err := http.Get(URL)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	fileName := path.Join("./", s.Folder) + uuid2.New().String() + ".mp3"

	f, err := os.Create(fileName)

	if err != nil {
		return "", err
	}

	defer f.Close()

	io.Copy(f, res.Body)

	return fileName, nil

}

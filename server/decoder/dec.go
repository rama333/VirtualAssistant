package decoder

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Decoder struct {
}

func NewDecoder() (*Decoder)  {

	return new(Decoder)
}


func (d Decoder) Dec(inputFile, outputFile string) (error)  {
	arg := []string{
		"/usr/bin/ffmpeg",
		"-i", inputFile,
		"-ac", "1", "-y",
		"-ar", "8000",
		"-acodec", "pcm_s16le",
		outputFile,
	}

	cmd := exec.Command(arg[0], arg[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}

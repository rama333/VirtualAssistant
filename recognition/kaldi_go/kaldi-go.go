package kaldi_go

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//var (
//DefaultPath = "/home/ubuntu/kaldi/src/online2bin/online2-wav-nnet3-latgen-faster"
//DoEndpointing = false
//AcousticScale = 1.0
//Beam = 13
//LatticeBeam = 6.0
//MaxActive = 7000
//ConfigPathOnline = ""
//WordSymbolTable = ""
//ConfigPathFinal = ""
//ConfigPathHCLG = ""
//PathWavFile = ""
//)

type Config struct {
	DefaultPath string
	DoEndpointing bool
	FrameSubsamplingFactor int
	AcousticScale float64
	Beam  float32
	LatticeBeam float32
	MaxActive int
	ConfigPathOnline string
	WordSymbolTable string
	ConfigPathFinal string
	ConfigPathHCLG string
	//PathWavFile string
}

type Encoder struct {
	config *Config
}

func NewConfig(config *Config) (*Encoder)  {

	c := &Encoder{config: config}

	return c
}

func ResultConfig(list ... string)  {
	logrus.Info(list)
}

func (e *Encoder) Recognition(file string) (string, error)  {

	t := time.Now()

	defer func() {
		log.Printf("время распознавания %s", time.Now().Sub(t).String())
	}()

	arg := []string{e.config.DefaultPath,
		"--do-endpointing=" + strconv.FormatBool(e.config.DoEndpointing),
		"--frame-subsampling-factor="+ strconv.Itoa(e.config.FrameSubsamplingFactor),
		"--acoustic-scale=" + fmt.Sprintf("%.1f", e.config.AcousticScale),
		"--beam=" + fmt.Sprintf("%.1f", e.config.Beam),
		"--minimize=" + "true",
		"--lattice-beam=" + fmt.Sprintf("%.1f", e.config.LatticeBeam),
		"--max-active=" + strconv.Itoa(e.config.MaxActive),
		"--config=" + e.config.ConfigPathOnline,
		"--word-symbol-table=" + e.config.WordSymbolTable,
		"" + e.config.ConfigPathFinal,
		"" + e.config.ConfigPathHCLG,
		"ark:echo utterance-id1 utterance-id1|",
		"scp:echo utterance-id1 "+ file+"|",
		"ark:- --acoustic-scale=0.1",
	}

	ResultConfig(arg[:] ...)

	cmd := exec.Command(arg[0], arg[1:] ...)
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()


	if err != nil {
		logrus.Info(  out.String() + out.String())
		log.Println(out.String() + stderr.String())
		return "", err
	}



	return strings.Split(strings.Split(stderr.String(),"utterance-id1 ")[3], "LOG")[0], nil
}




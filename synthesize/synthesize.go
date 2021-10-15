package synthesize

import "os/exec"

type Synthesize interface {
	Play(file string) error
}

type MPlayer struct {
}

func (p MPlayer) Play(file string) ( err error)  {

	cmd := exec.Command("mplayer", file)

	err = cmd.Run()

	return
}
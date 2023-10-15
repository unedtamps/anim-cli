package main

import (
	"os/exec"

	"github.com/blang/mpv"
)

func StartMpv() {
	cmd := exec.Command("mpv", "--idle", "--input-ipc-server=./mpvsocket")
	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}
	StartMpv()
}

func PlayVideo(link string) error {
	ipcc := mpv.NewIPCClient("./mpvsocket")

	c := mpv.NewClient(ipcc)
	err := c.Loadfile(link, mpv.LoadFileModeReplace)
	if err != nil {
		return err
	}
	err = c.SetFullscreen(true)
	if err != nil {
		return err
	}

	err = c.Seek(600, mpv.SeekModeAbsolute)
	if err != nil {
		return err
	}
	err = c.SetMute(false)
	if err != nil {
		return err
	}
	return nil
}

func DefaultPlay(episode map[string]string) map[string]string {
	for _, i := range episode {
		if err := PlayVideo(i); err != nil {
			continue
		}
		break
	}
	return episode
}

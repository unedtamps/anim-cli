package main

import (
	"os/exec"

	"github.com/blang/mpv"
)

func StartMpv() {
	cmd := exec.Command("mpv", "--idle", "--input-ipc-server=/tmp/mpvsocket")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	defer cmd.Cancel()
}

func PlayVideo(link string) error {
	ipcc := mpv.NewIPCClient("/tmp/mpvsocket")
	c := mpv.NewClient(ipcc)
	err := c.Loadfile(link, mpv.LoadFileModeReplace)
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
	err = c.PlaylistNext()
	if err != nil {
		return err
	}
	return nil
}

func PlayingVideo(episode_id string) map[string]string {
	video_url := GetVideoLink(episode_id)
	for _, i := range video_url {
		if err := PlayVideo(i); err != nil {
			continue
		}
		break
	}
	return video_url
}

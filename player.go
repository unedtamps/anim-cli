package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/blang/mpv"
)

func StartMedia() {
	cmd := exec.Command("mpv", "--idle", "--input-ipc-server=./mpvsocket")
	if err := cmd.Run(); err != nil {
		panic(err.Error())
	}
	StartMedia()
}

func PlayVideo(ctx context.Context, link string, subtitles []Subtitle) error {
	files, err := GetAllSubtitle(ctx, subtitles)
	if err != nil {
		return err
	}
	ipcc := mpv.NewIPCClient("./mpvsocket")
	c := mpv.NewClient(ipcc)
	err = c.Loadfile(link, mpv.LoadFileModeReplace)
	if err != nil {
		return err
	}
	err = c.SetFullscreen(true)
	if err != nil {
		return err
	}

	err = c.SetProperty(
		"sub-files",
		files,
	)
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
	chan1 := make(chan interface{})
	go func() {
		for {
			time.Sleep(1 * time.Second)
			pos, _ := c.Position()
			if pos > 0 {
				chan1 <- struct{}{}
				break
			}
		}

	}()
	<-chan1

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			os.Remove(file)
			defer wg.Done()
		}(file)
	}
	wg.Wait()

	return nil
}

func GetAllSubtitle(ctx context.Context, subtitles []Subtitle) ([]string, error) {
	var files []string
	chan1 := make(chan string)
	wg.Add(len(subtitles))

	go func() {
		for _, s := range subtitles {
			go func(s Subtitle) {
				extArr := strings.Split(s.Url, ".")
				file := fmt.Sprintf("%s.%s", s.Lang, extArr[len(extArr)-1])
				if strings.Contains(strings.ToLower(file), "eng") ||
					strings.Contains(strings.ToLower(file), "ind") {
					DownloadSubtitle(ctx, file, s.Url)
					chan1 <- file
				}
				wg.Done()
			}(s)
		}
	}()
	go func() {
		wg.Wait()
		close(chan1)
	}()

	for file := range chan1 {
		files = append(files, file)
	}
	return files, nil
}

func GetVideo(videos *SearchVideo, quality string) (*string, []Subtitle) {
	for _, v := range videos.Sources {
		if v.Quality == quality || v.Quality == fmt.Sprintf("%sp", quality) {
			return &v.Url, videos.Subtitles
		}
	}
	return nil, nil
}

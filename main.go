package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

var loger = logrus.New()
var wg sync.WaitGroup

func main() {
	ClearScreen()

	go StartMedia()
	ctx := context.Background()
	quit := make(chan os.Signal)

	go func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ClearScreen()
		Blue.Println("Exiting...")
		exec.Command("rm", "mpvsocket").Run()
		files, _ := filepath.Glob("*.vtt")
		wg.Add(len(files))
		for _, f := range files {
			go func(f string) {
				os.Remove(f)
				wg.Done()
			}(f)
		}
		wg.Wait()
		Green.Println("Done!!")
		os.Exit(0)
	}()

	for {
		keywords := PromptAsk("Enter the keywords...")
		if keywords == "" {
			Red.Println("Cannot Be Empty")
		}
		res, err := SearchKeyword(ctx, keywords)
		if err != nil {
			Red.Println(err.Error())
			continue
		}
		if res == nil {
			Red.Println("Not Found , Try Again!")
			continue
		}
		options, map_res := MapingResult(res, keywords)
		selected := Prompt("Select To Watch: ", options)
		result := map_res[selected]

		episodes, err := GetInfo(ctx, result)
		if err != nil {
			Red.Println(err.Error())
			continue
		}

		options, map_eps := MapingEpisode(result, episodes)

		select_eps := Prompt("Select <Season-Episode>...", options)

		episode_index := FindIndex(options, select_eps)

		episode := map_eps[select_eps]

		videos, err := GetVideos(ctx, result, episode)
		if err != nil {
			Red.Println(err.Error())
			continue
		}

		quality := "720"
		url, subtitles := GetVideo(videos, quality)
		if url == nil {
			Red.Println("Video Not Found")
			continue
		}
		PlayVideo(ctx, *url, subtitles)

	cinema:
		for {
			ClearScreen()
			var q_menu string
			if result.Type {
				q_menu = Cyan.Sprintf("Now Playing... %s Season %d %s", result.Title, episode.Season, episode.Title)
			} else {
				q_menu = Cyan.Sprintf("Now Playing... %s (%d) ", result.Title, episode.Number)
			}
			menus := []string{
				"next episode",
				"previous episode",
				"change quality",
				"search other",
				"select episode",
			}
			select_menu := Prompt(q_menu, menus)
			switch select_menu {
			case "select episode":
				select_eps = Prompt("Select <Season-Episode>...", options)
				episode = map_eps[select_eps]
				videos, err = GetVideos(ctx, result, episode)
				if err != nil {
					Red.Println(err.Error())
					continue
				}
				break
			case "search other":
				break cinema
			case "change quality":
				quality = SelectQuality(videos.Sources)
				break
			case "next episode":
				if episode_index >= len(options) {
					Red.Println("Episode not Found")
					continue
				}
				episode_index++
				episode = map_eps[options[episode_index]]
				videos, err = GetVideos(ctx, result, episode)
				if err != nil {
					Red.Println(err.Error())
					continue
				}
				break
			case "previous episode":
				if episode_index <= 0 {
					Red.Println("Episode not Found")
					continue
				}
				episode_index--
				episode = map_eps[options[episode_index]]
				videos, err = GetVideos(ctx, result, episode)
				if err != nil {
					Red.Println(err.Error())
					continue
				}
				break
			}

			url, subtitles := GetVideo(videos, quality)
			if url == nil {
				Red.Println(fmt.Errorf("Video Not Found"))
				continue
			}
			PlayVideo(ctx, *url, subtitles)
		}
	}

}

func SelectQuality(videos []Video) string {
	var qualities []string
	for _, v := range videos {
		qualities = append(qualities, v.Quality)
	}
	quality := Prompt("select quality", qualities)
	return quality
}

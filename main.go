package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var red = color.New(color.FgRed, color.Bold)
var loger = logrus.New()

func main() {
	// var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	api, err := NewAPIClient()
	if err != nil {
		loger.Fatal(err.Error())
	}

	go StartMpv()
	time.Sleep(500 * time.Millisecond)

	go func() {

		for {
			anime_name := ScanToSlug()
			err, res := api.SearchAnime(anime_name)
			if err != nil {
				//search again if error
				loger.Println(err.Error())
				continue
			}
			search_anime, map_anime := MapingAnime(res)
			if search_anime == nil {
				red.Println("Anime Not Found , Try Again!")
				continue
			}
			selected_anime := Prompt("Select Anime: ", search_anime)

			t_episode := api.GetTotalEpisode(selected_anime)
			if t_episode == 0 {
				loger.Info("Episode not found")
				continue
			}

			maping_id_eps := make(map[string]map[string]string)
			e_str, e_id := MapingEpisode(t_episode, map_anime[selected_anime])

			completed_maping := make(chan struct{})

			go api.ConcurentEpisode(t_episode%100, e_id, &maping_id_eps, completed_maping)
			eps_num := 0
			var video_url map[string]string

			for {
				sel_eps := Prompt("Select Episode", e_str)
				eps_num = toInt(sel_eps)
				video_url = DefaultPlay(api.GetVideoLink(e_id[eps_num]))
				if video_url == nil {
					continue
				}
				break
			}

			<-completed_maping
			for {
				eps_url := maping_id_eps[e_id[eps_num]]
				options := []string{
					"Next Episode",
					"Previous Episode",
					"Select Episode",
					"Change Quality",
					"Other Anime",
					"Stop Video",
					"Exit",
				}
				answer := Prompt("Menu", options)
				if answer == "Exit" {
					PlayVideo("")
					syscall.Kill(os.Getpid(), syscall.SIGINT)
					time.Sleep(1 * time.Second)
					return
				} else if answer == "Select Episode" {
					num_eps := Prompt("Select Episode", e_str)
					eps_num = toInt(num_eps)
					eps_url = maping_id_eps[e_id[eps_num]]
					if eps_url != nil {
						video_url = DefaultPlay(eps_url)
						continue
					}
					video_url = DefaultPlay(api.GetVideoLink(e_id[eps_num]))
					continue
				} else if answer == "Other Anime" {
					break
				} else if answer == "Change Quality" {
					optionq := []string{"360p", "480p", "720p", "1080p"}
					Quality := Prompt("Select Quality", optionq)
					PlayVideo(video_url[Quality])
					continue
				} else if answer == "Stop Video" {
					PlayVideo("")
					continue
				} else if answer == "Previous Episode" {
					eps_num -= 1
					eps_url = maping_id_eps[e_id[eps_num]]

					if eps_num <= 0 {
						red.Println("Episode not found")
						continue
					}
					if eps_url != nil {
						video_url = DefaultPlay(eps_url)
						continue
					}
					video_url = DefaultPlay(api.GetVideoLink(e_id[eps_num]))
					continue
				}
				eps_num += 1
				if eps_num > t_episode {
					red.Println("Episode not found")
					continue
				}
				eps_url = maping_id_eps[e_id[eps_num]]
				if eps_url != nil {
					video_url = DefaultPlay(eps_url)
					continue
				}
				video_url = DefaultPlay(api.GetVideoLink(e_id[eps_num]))
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("Gracefully Exiting...")
	exec.Command("rm", "mpvsocket").Run()
	os.Exit(0)
}

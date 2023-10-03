package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed, color.Bold)

func main() {
	var wg sync.WaitGroup
	var command string
	var port string
	cmd_helper := "(start|run) start = have container , run = make new"

	flag.StringVar(&command, "cmd", "start", cmd_helper)
	flag.StringVar(&port, "p", "3000", "port")
	flag.Parse()

	container_url = fmt.Sprintf("http://localhost:%s", port)

	api, err := NewAPIClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	done := make(chan struct{})
	go FlagHelper(command, api.Api.Image, port, api.Api.Container, done)
	go StartMpv()
	<-done

	for {
		anime_name := ScanToSlug()
		err, res := api.SearchAnime(anime_name)
		if err != nil {
			//search again if error
			log.Println(err.Error())
			continue
		}
		search_anime, map_anime := MapingAnime(res)
		if search_anime == nil {
			fmt.Println("Anime Not Found , Try Again!")
			continue
		}
		sel_anime := Prompt("Select Anime: ", search_anime)

		//looping until we get result
		var t_episode int
		wg.Add(1)
		go func() {
			for {
				t_episode = api.GetTotalEpisode(map_anime[sel_anime])
				if t_episode == 0 {
					continue
				}
				wg.Done()
				break
			}
		}()
		wg.Wait()
		maping_id_eps := make(map[string]map[string]string)
		e_str, e_id := MapingEpisode(t_episode, map_anime[sel_anime])

		go api.ConcurentEpisode(t_episode, e_id, &maping_id_eps)
		sel_eps := Prompt("Select Episode", e_str)
		eps_num := toInt(sel_eps)
		video_url := DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))

		for {
			eps_num += 1
			eps_url := maping_id_eps[e_id[eps_num]]
			options := []string{"Next", "Previous", "Select Episode", "Change Quality", "Change Anime", "Stop Player", "Quit"}
			answer := Prompt("Options", options)
			if answer == "Quit" {
				PlayVideo("")
				go StopServer(api.Api.Container)
				return
			} else if answer == "Select Episode" {
				num_eps := Prompt("Select Episode", e_str)
				eps_num = toInt(num_eps)
				eps_url = maping_id_eps[e_id[eps_num]]
				if eps_url != nil {
					video_url = DefaultPlay(eps_url)
					continue
				}
				video_url = DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))
				continue
			} else if answer == "Change Anime" {
				break
			} else if answer == "Change Quality" {
				optionq := []string{"360p", "480p", "720p", "1080p"}
				Quality := Prompt("Select Quality", optionq)
				PlayVideo(video_url[Quality])
				continue
			} else if answer == "Stop Player" {
				PlayVideo("")
				continue
			} else if answer == "Previous" {
				eps_num -= 2
				if eps_num <= 0 {
					red.Println("Episode not found")
					continue
				}
				if eps_url != nil {
					video_url = DefaultPlay(eps_url)
					continue
				}
				video_url = DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))
				continue
			}
			if eps_num > t_episode {
				red.Println("Episode not found")
				continue
			}
			if eps_url != nil {
				video_url = DefaultPlay(eps_url)
				continue
			}
			video_url = DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))
		}
	}
}

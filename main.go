package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var command string
	var port string
	cmd_helper := "(start|run) start = have container , run = make new"

	api, err := NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	flag.StringVar(&command, "cmd", "start", cmd_helper)
	flag.StringVar(&port, "p", "3000", "port")
	flag.Parse()

	go FlagHelper(command, api.Api.Image, port, api.Api.Container)
	go StartMpv()
	for {
		anime_name := ScanToSlug()
		res := api.SearchAnime(anime_name)
		search_anime, map_anime := MapingAnime(res)
		if search_anime == nil {
			fmt.Println("Anime Not Found , Try Again!")
			continue
		}
		sel_anime := Prompt("Select Anime: ", search_anime)
		t_episode := api.GetTotalEpisode(map_anime[sel_anime])
		e_str, e_id := MapingEpisode(t_episode, map_anime[sel_anime])
		sel_eps := Prompt("Select Episode", e_str)
		eps_num := toInt(sel_eps)
		video_url := DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))
		for {
			eps_num += 1
			options := []string{"Next", "Select Episode", "Change Quality", "Change Anime", "Stop Player", "Quit"}
			answer := Prompt("Options", options)
			if answer == "Quit" {
				StopServer(api.Api.Container)
				return
			} else if answer == "Select Episode" {
				num_eps := Prompt("Select Episode", e_str)
				eps_num = toInt(num_eps)
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
			}
			DefaultPlay(api.GetEpisodeUrl(e_id[eps_num]))
		}
	}
}

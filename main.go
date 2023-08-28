package main

import (
	"fmt"
	"log"
)

func main() {
	go StartMpv()
	api, err := NewClient()
	if err != nil {
		log.Fatal(err.Error())
	}
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

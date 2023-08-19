package main

import (
	"fmt"
	"time"
)

func main() {

	go StartMpv()
	time.Sleep(500 * time.Millisecond)
	for {
		anime_name := ScanToSlug()
		anime_res := SearchAnime(anime_name)
		anime_select, map_anime := MapingAnime(anime_res)
		if anime_select == nil {
			fmt.Println("Anime Not Found , Try Again!")
			continue
		}
		anime_selected := Prompt("Select Anime: ", anime_select)
		total_episode := GetTotalEpisode(map_anime[anime_selected])
		EpisodeNum, EpisodesId := MapingEpisode(total_episode, map_anime[anime_selected])
		NumEpisode := Prompt("Select Episode", EpisodeNum)
		number_episode := toInt(NumEpisode)
		video_url := PlayingVideo(EpisodesId[number_episode])

		for {
			number_episode += 1
			options := []string{"Next", "Select Episode", "Change Quality", "Change Anime", "Stop Player", "Quit"}
			answer := Prompt("Options", options)
			if answer == "Quit" {
				return
			} else if answer == "Select Episode" {
				NumEpisode := Prompt("Select Episode", EpisodeNum)
				number_episode = toInt(NumEpisode)
				video_url = PlayingVideo(EpisodesId[number_episode])
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
			PlayingVideo(EpisodesId[number_episode])
		}
	}
}

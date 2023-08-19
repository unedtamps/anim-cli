package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SearchAnime(anime_name string) AnimeSearchResponse {

	url := fmt.Sprintf("https://consumet-api-tam1.onrender.com/anime/gogoanime/%s", anime_name)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic("not found")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var anime AnimeSearchResponse
	err = json.Unmarshal(body, &anime)
	if err != nil {
		panic(err)
	}
	return anime
}

func GetVideoLink(id string) map[string]string {
	var video_url VideoUrl
	url := fmt.Sprintf("https://consumet-api-tam1.onrender.com/anime/gogoanime/watch/%s", id)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &video_url); err != nil {
		panic(err)
	}
	quality_link := make(map[string]string)
	lenght_arr := len(video_url.Sources)
	for i := lenght_arr - 3; i >= 0; i-- {
		quality_link[video_url.Sources[i].Quality] = video_url.Sources[i].Url
	}
	return quality_link
}

func GetTotalEpisode(id string) int {
	var AnimeInfo map[string]interface{}
	url := fmt.Sprintf("https://consumet-api-tam1.onrender.com/anime/gogoanime/info/%s", id)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &AnimeInfo); err != nil {
		panic(err)
	}
	totalEpisode, ok := AnimeInfo["totalEpisodes"].(float64)
	if !ok {
		panic("Episode Not Found")
	}
	return int(totalEpisode)
}

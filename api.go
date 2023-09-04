package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func NewAPIClient() (*Apilink, error) {
	api_link := new(Apilink)
	api_link.Api.Url = container_url
	api_link.Api.Image = container_image
	api_link.Api.Container = container_name
	return api_link, nil
}

func (a *Apilink) SearchAnime(anime_name string) SearchResponse {

	url := fmt.Sprintf("%s/anime/gogoanime/%s", a.Api.Url, anime_name)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*3)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Api Link Is Not Valid")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	var res_body SearchResponse
	err = json.Unmarshal(body, &res_body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return res_body
}

func (a *Apilink) GetVideoLink(id string) map[string]string {
	var video_url VideoUrl
	url := fmt.Sprintf("%s/anime/gogoanime/watch/%s", a.Api.Url, id)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*3)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &video_url); err != nil {
		log.Fatalf("%v", err)
	}
	quality_link := make(map[string]string)
	lenght_arr := len(video_url.Sources)
	for i := lenght_arr - 3; i >= 0; i-- {
		quality_link[video_url.Sources[i].Quality] = video_url.Sources[i].Url
	}
	return quality_link
}

func (a *Apilink) GetTotalEpisode(id string) int {
	var AnimeInfo map[string]interface{}
	url := fmt.Sprintf("%s/anime/gogoanime/info/%s", a.Api.Url, id)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*3)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := json.Unmarshal(data, &AnimeInfo); err != nil {
		log.Fatalf("%v", err)
	}
	totalEpisode, ok := AnimeInfo["totalEpisodes"].(float64)
	if !ok {
		panic("Episode Not Found")
	}
	return int(totalEpisode)
}

func (a *Apilink) GetEpisodeUrl(id string) map[string]string {
	video_url := a.GetVideoLink(id)
	return video_url
}

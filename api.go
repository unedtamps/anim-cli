package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

func NewClient() (*Apilink, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	api_link := new(Apilink)
	file, err := os.ReadFile(fmt.Sprintf("%s/config.yaml", pwd))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, api_link); err != nil {
		return nil, err
	}
	return api_link, nil
}

// func UpdateContainer(cont string) error {
// 	api_link := new(Apilink)
// 	file, err := os.ReadFile("./config.yaml")
// 	if err != nil {
// 		return err
// 	}
// 	if err := yaml.Unmarshal(file, api_link); err != nil {
// 		return err
// 	}
// 	api_link.Api.Container = cont
// 	data, err := yaml.Marshal(api_link)
// 	if err != nil {
// 		return err
// 	}
// 	if err := os.WriteFile("./config.yaml", data, os.ModePerm); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (a *Apilink) SearchAnime(anime_name string) SearchResponse {

	url := fmt.Sprintf("%s/anime/gogoanime/%s", a.Api.Url, anime_name)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic("Api Link Is Not Valid")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var res_body SearchResponse
	err = json.Unmarshal(body, &res_body)
	if err != nil {
		panic(err.Error())
	}
	return res_body
}

func (a *Apilink) GetVideoLink(id string) map[string]string {
	var video_url VideoUrl
	url := fmt.Sprintf("%s/anime/gogoanime/watch/%s", a.Api.Url, id)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	if res.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &video_url); err != nil {
		panic(err.Error())
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

func (a *Apilink) GetEpisodeUrl(id string) map[string]string {
	video_url := a.GetVideoLink(id)
	return video_url
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func NewAPIClient() (*Apilink, error) {
	api_link := new(Apilink)
	api_link.Api.Url = container_url
	api_link.Api.Image = container_image
	api_link.Api.Container = container_name
	return api_link, nil
}

func (a *Apilink) SearchAnime(anime_name string) (error, SearchResponse) {

	url := fmt.Sprintf("%s/anime/gogoanime/%s", a.Api.Url, anime_name)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*10)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err, SearchResponse{}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("status not ok"), SearchResponse{}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err, SearchResponse{}
	}
	var res_body SearchResponse
	err = json.Unmarshal(body, &res_body)
	if err != nil {
		return err, SearchResponse{}
	}
	return nil, res_body
}

func (a *Apilink) GetVideoLink(id string) map[string]string {
	var video_url VideoUrl
	quality_link := make(map[string]string)

	url := fmt.Sprintf("%s/anime/gogoanime/watch/%s", a.Api.Url, id)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*10)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		// log.Printf("%v", err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		loger.Fatalf("%v", "Anime is not found")
		return nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &video_url); err != nil {
		loger.Printf("%v", err)
		return nil
	}
	lenght_arr := len(video_url.Sources)
	for i := lenght_arr - 2; i >= 0; i-- {
		quality_link[video_url.Sources[i].Quality] = video_url.Sources[i].Url
	}
	return quality_link
}

func (a *Apilink) GetTotalEpisode(id string) int {
	var AnimeInfo map[string]interface{}
	url := fmt.Sprintf("%s/anime/gogoanime/info/%s", a.Api.Url, id)
	req, err := http.NewRequest("GET", url, nil)
	ctx, cencel := context.WithTimeout(context.Background(), time.Second*9)
	req = req.WithContext(ctx)
	defer cencel()

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		loger.Printf("%v", err)
		return 0
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		loger.Printf("%v", err)
		return 0
	}
	if err := json.Unmarshal(data, &AnimeInfo); err != nil {
		loger.Printf("%v", err)
		return 0
	}
	totalEpisode, ok := AnimeInfo["totalEpisodes"].(float64)
	if !ok {
		return 0
	}
	return int(totalEpisode)
}

func (a *Apilink) GetEpisodeUrl(id string) map[string]string {
	videochan := make(chan map[string]string)
	go func() {
		for {
			video_url := a.GetVideoLink(id)
			if video_url == nil {
				continue
			}
			videochan <- video_url
		}
	}()
	return <-videochan
}

func (a *Apilink) ConcurentEpisode(total_e int, e_id map[int]string, eps_url *map[string]map[string]string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(total_e)
	for i := 1; i <= total_e; i++ {
		go func(i int) {
			video_url := a.GetEpisodeUrl(e_id[i])
			mu.Lock()
			(*eps_url)[e_id[i]] = video_url
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
}

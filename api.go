package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func NewAPIClient() (*Api, error) {
	if API_URL == "" {
		return nil, errors.New("NO API")
	}
	return &Api{API_URL}, nil
}

func (a *Api) SearchAnime(anime_name string) (error, SearchResponse) {

	url := fmt.Sprintf("%s/anime/gogoanime/%s", a.Url, anime_name)
	req, err := http.NewRequest("GET", url, nil)
	ctx := context.Background()
	req = req.WithContext(ctx)

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

func (a *Api) GetVideoLink(id string) map[string]string {
	var video_url VideoUrl
	quality_link := make(map[string]string)

	url := fmt.Sprintf("%s/anime/gogoanime/watch/%s", a.Url, id)
	req, err := http.NewRequest("GET", url, nil)
	ctx := context.Background()
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
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

func (a *Api) GetTotalEpisode(title string) int {

	url := fmt.Sprintf("%s/anime/zoro/%s", a.Url, title)

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
	var AnimeInfo SearchAnime
	if err := json.Unmarshal(data, &AnimeInfo); err != nil {
		loger.Printf("%v", err)
		return 0
	}

	title_split := strings.Split(title, " ")
	if title_split[len(title_split)-1] == "(Dub)" {
		title = strings.Join(title_split[:len(title_split)-1], " ")
	}

	length := len(AnimeInfo.Results)
	for i := 0; i < length; i++ {
		anime := AnimeInfo.Results[i]
		if strings.ToLower(anime.JapaneseTitle) == strings.ToLower(title) {
			return anime.Sub
		}
	}
	return 0

}

func (a *Api) ConcurentEpisode(
	total_e int,
	e_id map[int]string,
	eps_url *map[string]map[string]string,
	done chan<- struct{},
) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(total_e)
	for i := 1; i <= total_e; i++ {
		go func(i int) {
			video_url := a.GetVideoLink(e_id[i])
			mu.Lock()
			(*eps_url)[e_id[i]] = video_url
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	done <- struct{}{}
}

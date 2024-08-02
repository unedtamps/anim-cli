package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Request(ctx context.Context, data interface{}, url string) error {

	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == 404 {
			return fmt.Errorf("Not Found")
		}
		return fmt.Errorf("Server Error")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return nil
}

func SearchKeyword(ctx context.Context, keyword string) ([]Result, error) {

	var results []Result

	anime_url := fmt.Sprintf("%s/anime/gogoanime/%s", API_URL, keyword)
	movie_url := fmt.Sprintf("%s/movies/flixhq/%s", API_URL, keyword)
	url := []string{movie_url, anime_url}

	for i, u := range url {
		var response SearchResponse
		err := Request(ctx, &response, u)
		if err != nil {
			return nil, err
		}
		for _, r := range response.Results {
			result := Result{
				Id:          r.Id,
				Title:       r.Title,
				Type:        i == 0,
				ReleaseDate: r.ReleaseDate,
			}
			results = append(results, result)
		}
	}

	return results, nil
}

func GetInfo(ctx context.Context, query Result) ([]Episode, error) {
	var url string
	if query.Type {
		url = fmt.Sprintf("%s/movies/flixhq/info?id=%s", API_URL, query.Id)
	} else {
		url = fmt.Sprintf("%s/anime/gogoanime/info/%s", API_URL, query.Id)
	}

	var response Info

	err := Request(ctx, &response, url)
	if err != nil {
		return nil, err
	}
	return response.Episode, nil
}

func GetVideos(ctx context.Context, media Result, query Episode) (*SearchVideo, error) {
	var url string
	if media.Type {
		url = fmt.Sprintf(
			"%s/movies/flixhq/watch?episodeId=%s&mediaId=%s",
			API_URL,
			query.Id,
			media.Id,
		)
	} else {
		url = fmt.Sprintf("%s/anime/gogoanime/watch/%s", API_URL, query.Id)
	}
	var searchRes SearchVideo
	err := Request(ctx, &searchRes, url)
	if err != nil {
		return nil, err
	}
	return &searchRes, nil
}

func DownloadSubtitle(ctx context.Context, filename string, urlStr string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(urlStr)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Subtitle Not found")
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func preFetchVideos(
	ctx context.Context,
	map_eps map[string]Episode,
	media Result,
	keys []string,
	index int,
	video_chan chan<- *SearvideoChan,
) {
	if len(keys) > 40 {
		start := index - 20
		end := index + 20

		if start < 0 {
			start = 0
		}
		if end > len(keys) {
			end = len(keys)
		}

		keys = keys[start:end]
	}

	wg.Add(len(keys))

	for _, k := range keys {
		go func(k string) {
			defer wg.Done()
			videos, err := GetVideos(ctx, media, map_eps[k])
			if err != nil {
				video_chan <- nil
				return
			}
			video_chan <- &SearvideoChan{
				Key:         k,
				SearchVideo: *videos,
			}
		}(k)
	}

	go func() {
		Magenta.Println("Please wait...")
		wg.Wait()
		close(video_chan)
	}()

}

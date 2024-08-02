package main

type SearchResult struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}
type Result struct {
	Id          string
	Title       string
	Type        bool
	ReleaseDate string
}

type SearchResponse struct {
	Results []SearchResult `json:"results,omitempty"`
}

type Info struct {
	Episode []Episode `json:"episodes,omitempty"`
}

type Episode struct {
	Id     string `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title,omitempty"`
	Season int    `json:"season,omitempty"`
}

type SearchVideo struct {
	Sources   []Video    `json:"sources"`
	Subtitles []Subtitle `json:"subtitles,omitempty"`
}

type Video struct {
	Url     string `json:"url,omitempty"`
	Quality string `json:"quality,omitempty"`
}

type Subtitle struct {
	Url  string `json:"url"`
	Lang string `json:"lang"`
}

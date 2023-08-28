package main

type AnimeGeneral struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type SearchResponse struct {
	CurrentPage int            `json:"currentPage"`
	HasNextPage bool           `json:"hasNextPage"`
	Results     []AnimeGeneral `json:"results"`
}

type SourceVideo struct {
	Url     string `json:"url,omitempty"`
	IsM3U8  bool   `json:"isM3U8,omitempty"`
	Quality string `json:"quality,omitempty"`
}

type HeaderVideo struct {
	Referer string `json:"referer,omitempty"`
}

type VideoUrl struct {
	Headers HeaderVideo   `json:"headers,omitempty"`
	Sources []SourceVideo `json:"sources,omitempty"`
}

type Apilink struct {
	Api Link `yaml:"api"`
}

type Link struct {
	Url string `yaml:"url"`
}

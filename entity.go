package main

type AnimeGeneral struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type SearchResponse struct {
	CurrentPage int            `json:"currentPage,omitempty"`
	HasNextPage bool           `json:"hasNextPage,omitempty"`
	Results     []AnimeGeneral `json:"results,omitempty"`
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
	Url       string `yaml:"url"`
	Image     string `yaml:"image"`
	Container string `yaml:"container"`
}

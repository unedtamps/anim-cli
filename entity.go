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

type Api struct {
	Url string
}

type SearchAnime struct {
	CurrentPage int         `json:"currentPage,omitempty"`
	HasNextPage bool        `json:"hasNextPage,omitempty"`
	TotalPage   int         `json:"totalPage,omitempty"`
	Results     []AnimeInfo `json:"results,omitempty"`
}

type AnimeInfo struct {
	Id            string `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Url           string `json:"url,omitempty"`
	Image         string `json:"image,omitempty"`
	Duration      string `json:"duration,omitempty"`
	JapaneseTitle string `json:"japaneseTitle,omitempty"`
	Type          string `json:"type,omitempty"`
	Nsfw          bool   `json:"nsfw,omitempty"`
	Sub           int    `json:"sub,omitempty"`
	Dub           int    `json:"dub,omitempty"`
	Episode       int    `json:"episodes,omitempty"`
}

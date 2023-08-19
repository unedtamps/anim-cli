package main

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
)

func ScanToSlug() string {
	fmt.Print("Search Anime: ")
	reader := bufio.NewReader(os.Stdin)

	anime_name, _ := reader.ReadString('\n')
	anime_name = slug.Make(anime_name)
	return anime_name
}

func toInt(num string) int {
	integer, _ := strconv.Atoi(num)
	return integer
}

func toStr(num int) string {
	return strconv.Itoa(num)
}

func ToSlug(a ...string) string {
	return strings.Join(a, "-")
}

func MapingAnime(anime AnimeSearchResponse) ([]string, map[string]string) {

	var anime_select []string
	map_id := make(map[string]string)
	for _, an := range anime.Results {
		anime_select = append(anime_select, an.Title)
		map_id[an.Title] = an.Id
	}
	return anime_select, map_id
}

func MapingEpisode(total_episode int, anime_id string) ([]string, map[int]string) {
	var EpisodeNum []string
	EpisodesId := make(map[int]string)
	for i := 1; i <= total_episode; i++ {
		index := fmt.Sprintf("%d", i)
		EpisodeId := ToSlug(anime_id, "episode", toStr(i))
		EpisodesId[toInt(index)] = EpisodeId
		EpisodeNum = append(EpisodeNum, index)
	}
	return EpisodeNum, EpisodesId
}

func Prompt(q string, opt []string) string {
	var answer string
	prompt := &survey.Select{
		Message:  q,
		Options:  opt,
		Default:  opt[0],
		PageSize: len(opt),
		VimMode: true,
	}
	if err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required)); err != nil {
		panic(err)
	}
	return answer
}

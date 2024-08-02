package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
)

var (
	BlueBg  = color.New(color.BgHiBlue, color.FgBlack)
	Blue    = color.New(color.FgHiBlue, color.Bold)
	Yellow  = color.New(color.FgHiYellow)
	RedBg   = color.New(color.BgHiRed, color.FgBlack)
	GreenBg = color.New(color.BgHiGreen, color.FgBlack)
	Green   = color.New(color.FgHiGreen, color.Bold)
	Magenta = color.New(color.FgHiMagenta, color.Bold)
	Red     = color.New(color.BgHiRed)
	Cyan    = color.New(color.FgHiCyan, color.Bold)
)

func ClearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}

func MapingResult(results []Result, keyword string) ([]string, map[string]Result) {

	var options []string
	map_res := make(map[string]Result)
	for _, r := range results {
		option := fmt.Sprintf("%s %s", r.Title, r.ReleaseDate)
		options = append(options, option)
		map_res[option] = r
	}
	a := []rune(strings.ToLower(keyword))
	firstWord := string(a[0:1])
	sort.Slice(options, func(i, j int) bool {
		is := strings.HasPrefix(strings.ToLower(options[i]), firstWord)
		js := strings.HasPrefix(strings.ToLower(options[j]), firstWord)
		if is == js {
			return i < j
		}
		return is && !js
	})

	return options, map_res
}

func MapingEpisode(media Result, episodes []Episode) ([]string, map[string]Episode) {
	var options []string
	map_res := make(map[string]Episode)
	for _, e := range episodes {
		var option string
		if media.Type {
			title := strings.Split(e.Title, ":")[1]
			option = fmt.Sprintf("%d-%d (%s)", e.Season, e.Number, strings.TrimSpace(title))
		} else {
			option = fmt.Sprintf("%d", e.Number)
		}
		options = append(options, option)
		map_res[option] = e
	}
	return options, map_res
}

func Prompt(q string, options []string) string {
	lable := Yellow.Sprintf("%s <vim cmd>", q)

	templates := promptui.SelectTemplates{
		Label:    "   {{ . | bold}}",
		Active:   "⏩ {{ . | green | bold}}",
		Inactive: "  {{ . | blue }}",
		Selected: "",
	}
	searcher := func(input string, index int) bool {
		pepper := options[index]
		name := strings.Replace(strings.ToLower(pepper), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	selector := promptui.Select{
		HideSelected:      true,
		Label:             lable,
		Items:             options,
		CursorPos:         0,
		HideHelp:          true,
		Templates:         &templates,
		IsVimMode:         true,
		Size:              10,
		StartInSearchMode: true,
		Searcher:          searcher,
	}
	_, res, err := selector.Run()
	if err != nil {
		Red.Println(err.Error())
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}
	return res
}

func PromptAsk(q string) string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }} ",
		Valid:   "{{ . | blue | bold }} ",
		Invalid: "{{ . | red | bold | italic }} ",
		Success: "{{ . | green | bold }} ",
	}

	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("Cannot Empty")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s", q),
		Templates: templates,
		Validate:  validate,
	}

	res, err := prompt.Run()
	if err != nil {
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}
	return res
}

func FindIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

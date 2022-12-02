package phonemic

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
	//go:embed data
	data            embed.FS
	dataFileNameMap = map[string]map[string]struct{}{LocaleThTH: {SystemPaiboon: struct{}{}}}
)

// this code inspired by
// https://stackoverflow.com/questions/8870261
// https://github.com/willsmil/go-wordninja

type Service struct {
	cfg         Config
	data        map[string]string
	wordCostMap map[string]float64
	maxLenWord  int
}

func NewService(cfg Config) *Service {
	if cfg == nil {
		return nil
	}

	return &Service{cfg: cfg, data: make(map[string]string), wordCostMap: make(map[string]float64)}
}

func (s *Service) Init() error {
	systemMap, ok := dataFileNameMap[s.cfg.GetLocale()]
	if !ok {
		return fmt.Errorf("unsupported locale %s", s.cfg.GetLocale())
	}

	if _, ok := systemMap[s.cfg.GetSystem()]; !ok {
		return fmt.Errorf("unsupported phonemic system %s", s.cfg.GetSystem())
	}

	dataFile, err := data.Open(fmt.Sprintf("data/%s/%s.data", s.cfg.GetLocale(), s.cfg.GetSystem()))
	if err != nil {
		return fmt.Errorf("unable to parse phonemic system data: %w", err)
	}
	defer func() {
		if err = dataFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	scanner := bufio.NewScanner(dataFile)
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, "\t")
		word := parts[0]
		s.data[word] = strings.Replace(parts[1], "ˑ", "", -1)

		wordLen := len(word)
		if wordLen > s.maxLenWord {
			s.maxLenWord = wordLen
		}
		s.wordCostMap[word] = 1 / float64(wordLen)
	}

	return nil
}

func (s *Service) Transcript(definition string) string {
	if transcription, ok := s.data[definition]; ok {
		return transcription
	}

	costs := []float64{0}
	runesDef := []rune(definition)
	lenDef := len(runesDef)
	for i := 1; i < lenDef+1; i++ {
		if m, err := s.bestMatch(runesDef, costs, i); err == nil {
			costs = append(costs, m.cost)
		}
	}

	var out []string
	i := lenDef

	for i > 0 {
		m, err := s.bestMatch(runesDef, costs, i)
		if err != nil {
			continue
		}

		sentence := strings.Replace(string(runesDef[i-m.idx:i]), " ", "", -1)
		sentenceTr, ok := s.data[sentence]
		if !ok {
			sentenceTr = sentence
		}
		out = append(out, sentenceTr)

		i -= m.idx
	}

	length := len(out)
	for i := 0; i < length/2; i++ {
		out[i], out[length-i-1] = out[length-i-1], out[i]
	}

	return strings.Join(out, " ")
}

type match struct {
	cost float64
	idx  int
}

// bestMatch will return the minimal cost and its appropriate character's index.
func (s *Service) bestMatch(text []rune, costs []float64, i int) (*match, error) {
	max := i - s.maxLenWord
	if max < 0 {
		max = 0
	}
	candidates := costs[max:i]
	k := 0
	var matchs []*match
	for j := len(candidates) - 1; j >= 0; j-- {
		part := text[i-k-1 : i]
		cost := s.getWordCost(string(part)) + candidates[j]
		matchs = append(matchs, &match{cost: cost, idx: k + 1})
		k++
	}

	return minCost(matchs)
}

// getWordCost return cost of word from the wordCostMap map.
// if the word is not exist in the map, it will return `9e99`.
func (s *Service) getWordCost(word string) float64 {
	if v, ok := s.wordCostMap[word]; ok {
		return v
	}

	return 9e99
}

func minCost(matchs []*match) (*match, error) {
	if len(matchs) == 0 {
		return nil, errors.New("match.len")
	}

	r := matchs[0]
	for _, m := range matchs {
		if m.cost < r.cost {
			r = m
		}
	}

	return r, nil
}

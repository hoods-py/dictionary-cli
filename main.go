package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio,omitempty"`
}

type Definition struct {
	Word       string     `json:"word"`
	Phonetic   string     `json:"phonetic"`
	Phonetics  []Phonetic `json:"phonetics"`
	Origin     string     `json:"origin"`
	Meanings   []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string   `json:"definition"`
			Example    string   `json:"example"`
			Synonyms   []string `json:"synonyms"`
			Antonyms   []string `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func getDefinition(word string) (Definition, error) {
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	response, err := http.Get(url)
	if err != nil {
		return Definition{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Definition{}, err
	}

	var definitions []Definition
	if err := json.Unmarshal(body, &definitions); err != nil {
		return Definition{}, err
	}

	if len(definitions) == 0 {
		return Definition{}, fmt.Errorf("no definitions found for the word %s", word)
	}

	return definitions[0], nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dictionary-cli <word>")
		return
	}

	word := os.Args[1]
	definition, err := getDefinition(word)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("%s [%s]\n", definition.Word, definition.Phonetic)
	for _, p := range definition.Phonetics {
		fmt.Printf("- %s", p.Text)
		if p.Audio != "" {
			fmt.Printf(" (Audio: %s)", p.Audio)
		}
		fmt.Println()
	}
	fmt.Printf("Origin: %s\n", definition.Origin)

	for _, meaning := range definition.Meanings {
		fmt.Printf("\n%s:\n", meaning.PartOfSpeech)
		for _, d := range meaning.Definitions {
			fmt.Printf("- %s\n", d.Definition)
			fmt.Printf("  Example: %s\n", d.Example)
			if len(d.Synonyms) > 0 {
				fmt.Printf("  Synonyms: %v\n", d.Synonyms)
			}
			if len(d.Antonyms) > 0 {
				fmt.Printf("  Antonyms: %v\n", d.Antonyms)
			}
		}
	}
}
// Path: go.mod

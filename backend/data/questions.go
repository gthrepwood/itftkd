package data

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
)

type VocabEntry struct {
	Topic       string `json:"topic"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Translation string `json:"translation"`
	Descr       string `json:"descr"`
	Description string `json:"description"`
	Script      string `json:"script"`
	Duped       bool   `json:"dupe"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
	Category string   `json:"category"`
}

var (
	entries   []VocabEntry
	byTopic   map[string][]VocabEntry
	topics    []string
	questions []Question
)

func getDisplayText(e VocabEntry) string {
	if e.Translation != "" {
		return strings.TrimSpace(e.Translation)
	}
	if e.Descr != "" {
		return strings.TrimSpace(e.Descr)
	}
	if e.Description != "" {
		return strings.TrimSpace(e.Description)
	}
	return ""
}

func getQuestionText(e VocabEntry) string {
	name := e.Name
	if e.Script != "" {
		name = e.Script + " (" + name + ")"
	}
	return "Mi a(z) " + name + " jelentése?"
}

func generateQuestionsForTopic(topic string) []Question {
	var result []Question
	items := byTopic[topic]
	var pool []VocabEntry
	for _, it := range items {
		if !it.Duped {
			pool = append(pool, it)
		}
	}
	if len(pool) < 4 {
		return nil
	}

	for i, item := range pool {
		correct := getDisplayText(item)
		if correct == "" {
			continue
		}

		var opts []string
		opts = append(opts, correct)

		filtered := make([]string, 0)
		for j, other := range pool {
			if j == i {
				continue
			}
			if t := getDisplayText(other); t != "" {
				filtered = append(filtered, t)
			}
		}
		rand.Shuffle(len(filtered), func(a, b int) {
			filtered[a], filtered[b] = filtered[b], filtered[a]
		})

		for _, f := range filtered {
			if len(opts) >= 4 {
				break
			}
			if f != correct {
				opts = append(opts, f)
			}
		}

		if len(opts) < 4 {
			continue
		}

		rand.Shuffle(len(opts), func(a, b int) {
			opts[a], opts[b] = opts[b], opts[a]
		})

		answerIdx := 0
		for j, o := range opts {
			if o == correct {
				answerIdx = j
				break
			}
		}

		result = append(result, Question{
			ID:       item.Number,
			Question: getQuestionText(item),
			Options:  opts,
			Answer:   answerIdx,
			Category: topic,
		})
	}

	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

func loadData() {
	file, err := os.Open("./data/data.json")
	if err != nil {
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.Decode(&entries)

	byTopic = make(map[string][]VocabEntry)
	var allTopics []string
	topicSet := make(map[string]bool)

	for _, e := range entries {
		if e.Topic == "" {
			continue
		}
		byTopic[e.Topic] = append(byTopic[e.Topic], e)
		if !topicSet[e.Topic] {
			topicSet[e.Topic] = true
			allTopics = append(allTopics, e.Topic)
		}
	}

	topics = allTopics

	questions = nil
	for t := range byTopic {
		q := generateQuestionsForTopic(t)
		questions = append(questions, q...)
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
}

func GetQuestions() []Question {
	if questions == nil {
		loadData()
	}
	return questions
}

func GetQuestionsForTopic(topic string) []Question {
	if questions == nil {
		loadData()
	}
	if topic == "" || topic == "all" {
		return questions
	}
	if topic == "random" {
		return GetRandomQuestions(len(questions))
	}
	return generateQuestionsForTopic(topic)
}

func GetTopics() []string {
	if topics == nil {
		loadData()
	}
	out := make([]string, len(topics))
	copy(out, topics)
	return out
}

func GetRandomQuestions(count int) []Question {
	q := GetQuestions()
	if count > len(q) {
		count = len(q)
	}
	rand.Shuffle(len(q), func(i, j int) {
		q[i], q[j] = q[j], q[i]
	})
	return q[:count]
}

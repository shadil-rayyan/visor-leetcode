package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// JSONOutput is the top-level shape written to web/src/data/problems.json.
type JSONOutput struct {
	GeneratedAt string        `json:"generatedAt"`
	Problems    []JSONProblem `json:"problems"`
	Companies   []JSONCompany `json:"companies"`
}

type JSONProblem struct {
	ID         int64          `json:"id"`
	Title      string         `json:"title"`
	URL        string         `json:"url"`
	Difficulty string         `json:"difficulty"`
	Acceptance *float64       `json:"acceptance"`
	Frequency  *float64       `json:"frequency"`
	Tags       []string       `json:"tags"`
	Companies  []JSONCompanyLink `json:"companies"`
}

type JSONCompanyLink struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Timeframe string `json:"timeframe"` // "" when none
}

type JSONCompany struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ProblemCount int    `json:"problemCount"`
}

// scrapeJsonMain reads CSVs from ROOT_DIR, optionally fetches tags from LeetCode,
// and writes a single problems.json the React app can import.
func scrapeJsonMain() {
	godotenv.Load()

	root := os.Getenv("ROOT_DIR")
	if root == "" {
		log.Fatal("ROOT_DIR env var required")
	}

	outPath := os.Getenv("JSON_OUT")
	if outPath == "" {
		outPath = "../web/src/data/problems.json"
	}

	skipTags := os.Getenv("SKIP_TAGS") == "true"

	// 1. Walk companies (top-level dirs in ROOT_DIR)
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Fatalf("read root: %v", err)
	}

	var companyNames []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		companyNames = append(companyNames, name)
	}
	sort.Strings(companyNames)
	companyIDByName := map[string]int{}
	companies := make([]JSONCompany, 0, len(companyNames))
	for i, name := range companyNames {
		companyIDByName[name] = i + 1
		companies = append(companies, JSONCompany{ID: i + 1, Name: name})
	}

	// 2. For each company, parse CSVs and collect problem → company links
	type linkInfo struct {
		companyID   int
		timeframe   string // best (most-recent) timeframe tag seen
	}
	type problemMeta struct {
		id         int64
		title      string
		url        string
		difficulty string
		acceptance *float64
		frequency  *float64
	}
	problems := map[int64]*JSONProblem{}
	problemsBySlug := map[int64]string{} // for tag lookup later
	linksByProblem := map[int64]map[int]linkInfo{}

	files := []struct{ key, path string }{
		{"all", "all.csv"},
		{"six-months", "six-months.csv"},
		{"three-months", "three-months.csv"},
		{"thirty-days", "thirty-days.csv"},
		{"more-than-six", "more-than-six-months.csv"},
	}
	// timeframe rank: thirty-days (best) > three-months > six-months > "" (none)
	tfRank := map[string]int{
		"thirty-days":   4,
		"three-months":  3,
		"six-months":    2,
		"more-than-six": 1,
		"":              0,
	}

	for _, name := range companyNames {
		companyPath := filepath.Join(root, name)
		cid := companyIDByName[name]

		// track which problems this company asked (to skip dupes within company)
		seen := map[int64]bool{}

		for _, f := range files {
			m, err := readCSVFile(filepath.Join(companyPath, f.path), f.key)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				log.Printf("read %s/%s: %v", name, f.path, err)
				continue
			}
			for id, rp := range m {
				if seen[id] {
					// dedupe within company — but upgrade timeframe if this file is more recent
					if linksByProblem[id][cid].timeframe == "" || tfRank[f.key] > tfRank[linksByProblem[id][cid].timeframe] {
						linksByProblem[id][cid] = linkInfo{cid, f.key}
					}
					continue
				}
				seen[id] = true

				// insert problem if new
				if _, ok := problems[id]; !ok {
					p := &JSONProblem{
						ID:         id,
						Title:      rp.Title,
						URL:        rp.URL,
						Difficulty: rp.Difficulty,
					}
					if rp.Acceptance.Valid {
						v := rp.Acceptance.Float64
						p.Acceptance = &v
					}
					if rp.Frequency.Valid {
						v := rp.Frequency.Float64
						p.Frequency = &v
					}
					p.Tags = []string{}
					problems[id] = p
					problemsBySlug[id] = extractSlug(rp.URL)
				}
				// link
				if linksByProblem[id] == nil {
					linksByProblem[id] = map[int]linkInfo{}
				}
				// always overwrite (f loop goes all → six → three → thirty, so later = better timeframe)
				linksByProblem[id][cid] = linkInfo{cid, f.key}
			}
		}
	}

	log.Printf("parsed %d problems across %d companies", len(problems), len(companyNames))

	// 3. Fetch tags from LeetCode GraphQL (unless SKIP_TAGS)
	if !skipTags {
		// reuse scrape_tags.go package-level: pslug, buildBatchQuery, doGraphQLRequest
		var items []pslug
		for id, slug := range problemsBySlug {
			if slug == "" {
				continue
			}
			items = append(items, pslug{ID: id, Slug: slug})
		}
		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })

		total := len(items)
		client := &http.Client{Timeout: time.Second * httpTimeoutSecs}
		for i := 0; i < total; i += batchSize {
			j := i + batchSize
			if j > total {
				j = total
			}
			batch := items[i:j]
			log.Printf("tags batch %d..%d", i, j-1)
			q := buildBatchQuery(batch)
			body, err := doGraphQLRequest(client, q)
			if err != nil {
				log.Printf("graphql fail %d..%d: %v", i, j-1, err)
				time.Sleep(time.Millisecond * politeSleepMS)
				continue
			}
			var br batchResponse
			if err := json.Unmarshal(body, &br); err != nil {
				log.Printf("unmarshal fail %d..%d: %v", i, j-1, err)
				continue
			}
			for idx, p := range batch {
				alias := fmt.Sprintf("q%d", idx)
				entry, ok := br.Data[alias]
				if !ok || entry.QuestionId == "" {
					continue
				}
				tags := []string{}
				for _, tt := range entry.TopicTags {
					name := strings.TrimSpace(tt.Name)
					if name != "" {
						tags = append(tags, name)
					}
				}
				if pp, ok := problems[p.ID]; ok {
					pp.Tags = tags
				}
			}
			time.Sleep(time.Millisecond * politeSleepMS)
		}
	} else {
		log.Println("SKIP_TAGS set — skipping tag fetch")
	}

	// 4. Flatten: attach company links to each problem
	out := JSONOutput{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Companies:   companies,
		Problems:    make([]JSONProblem, 0, len(problems)),
	}
	// problemCount per company
	countByCompany := map[int]int{}
	for id, pp := range problems {
		links := linksByProblem[id]
		for _, link := range links {
			pp.Companies = append(pp.Companies, JSONCompanyLink{
				ID:        link.companyID,
				Name:      companies[link.companyID-1].Name,
				Timeframe: link.timeframe,
			})
			countByCompany[link.companyID]++
		}
		// sort links alphabetically by company name for stable output
		sort.Slice(pp.Companies, func(i, j int) bool { return pp.Companies[i].Name < pp.Companies[j].Name })
		out.Problems = append(out.Problems, *pp)
	}
	// fill problemCount
	for i := range out.Companies {
		out.Companies[i].ProblemCount = countByCompany[out.Companies[i].ID]
	}
	// sort problems by id for stable output
	sort.Slice(out.Problems, func(i, j int) bool { return out.Problems[i].ID < out.Problems[j].ID })

	// 5. Write JSON
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("create %s: %v", outPath, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&out); err != nil {
		log.Fatalf("encode: %v", err)
	}
	log.Printf("wrote %s — %d problems, %d companies", outPath, len(out.Problems), len(out.Companies))
}

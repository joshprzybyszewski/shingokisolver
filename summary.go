package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

const (
	readmeFileName        = `README.md`
	latestResultsFileName = `latestResults.md`
	resultsStartString    = `</startResults>`
)

type summary struct {
	Name     string
	Unsolved string
	Solution string

	heapSize uint64
	pauseNS  uint64
	numGCs   uint32

	NumEdges int
	Duration time.Duration
}

func updateReadme(allSummaries []summary) {
	latestResults, readmeUpdate := buildAllSummariesOutput(allSummaries)

	if err := ioutil.WriteFile(latestResultsFileName, []byte(latestResults), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input, err := ioutil.ReadFile(readmeFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	inputStr := string(input)
	prefix := inputStr[:strings.Index(inputStr, resultsStartString)]

	if err = ioutil.WriteFile(readmeFileName, []byte(prefix+readmeUpdate), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildAllSummariesOutput(
	allSummaries []summary,
) (string, string) {

	summsBySize := make(map[int][]summary, 10)
	for i := range allSummaries {
		summ := allSummaries[i]
		summsBySize[summ.NumEdges] = append(summsBySize[summ.NumEdges], summ)
	}

	return buildLatestResultsOutput(summsBySize), buildSummaryBySize(summsBySize)
}

func buildLatestResultsOutput(
	summsBySize map[int][]summary,
) string {
	var sb strings.Builder

	sb.WriteString(resultsStartString)
	sb.WriteString("\n\n")
	sb.WriteString("# Results from ")
	sb.WriteString(time.Now().Format("01-02-2006"))

	sb.WriteString("\n\n")
	sb.WriteString("<table>")

	sb.WriteString(`<tr>
	<th>Name</th>
	<th>Duration</th>
	<th>Heap Size (bytes)</th>
	<th>Num Garbage Collections</th>
	<th>Pause (ns)</th>
	<th>Puzzle</th>
	<th>Solution</th>
</tr>
`)
	var unsolvedCell, solutionCell string

	for size := 1; size < state.MaxEdges; size++ {
		summaries, ok := summsBySize[size]
		if !ok {
			continue
		}

		sort.Slice(summaries, func(i, j int) bool {
			if summaries[i].NumEdges != summaries[j].NumEdges {
				return summaries[i].NumEdges < summaries[j].NumEdges
			}
			return summaries[i].Duration > summaries[j].Duration
		})

		lenSlowest := 10
		if len(summaries) < lenSlowest {
			lenSlowest = len(summaries)
		}
		slowestSumms := make([]summary, lenSlowest)
		copy(slowestSumms, summaries)
		sort.Slice(slowestSumms, func(i, j int) bool {
			return strings.Compare(slowestSumms[i].Name, slowestSumms[j].Name) < 0
		})

		for i := range slowestSumms {
			s := slowestSumms[i]
			unsolvedCell = fmt.Sprintf(
				"<details><summary>Puzzle</summary>\n\n```\n%s\n```\n</details>\n",
				s.Unsolved,
			)
			solutionCell = fmt.Sprintf(
				"<details><summary>Solution</summary>\n\n```\n%s\n```\n</details>\n",
				s.Solution,
			)

			sb.WriteString(fmt.Sprintf(`<tr>
	<td>%s</td>
	<td>%s</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%s</td>
	<td>%s</td>
</tr>
`,
				s.Name,
				s.Duration,
				s.heapSize,
				s.numGCs,
				s.pauseNS,
				unsolvedCell,
				solutionCell,
			))
		}
	}

	sb.WriteString("</table>")

	return sb.String()
}

func buildSummaryBySize(
	summsBySize map[int][]summary,
) string {

	var sb strings.Builder
	sb.WriteString(resultsStartString)
	sb.WriteString("\n\n")
	sb.WriteString("#### Results from ")
	sb.WriteString(time.Now().Format("01-02-2006"))
	sb.WriteString("\n\n")

	sb.WriteString("|Num Edges|")
	sb.WriteString("Sample Size|")
	sb.WriteString("Average Duration|")
	sb.WriteString("Average Allocations (KB)|")
	sb.WriteString("Average Garbage Collections|")
	sb.WriteString("Average GC Pause (ns)|")
	sb.WriteString("\n")
	sb.WriteString("|-:|-:|-:|-:|-:|-:|\n")

	for size := 1; size < state.MaxEdges; size++ {
		summaries, ok := summsBySize[size]
		if !ok {
			continue
		}

		if size != 2 && len(summaries) < 10 {
			compete.PopulateCache(size, 100-len(summaries))
		}

		var totalDur time.Duration
		var totalHeapBytes uint64
		var totalGCs uint32
		var totalPauseNS uint64

		for i := range summaries {
			totalDur += summaries[i].Duration
			totalHeapBytes += summaries[i].heapSize
			totalGCs += summaries[i].numGCs
			totalPauseNS += summaries[i].pauseNS
		}

		avgDur := totalDur / time.Duration(len(summaries))
		avgAllocs := float64(totalHeapBytes) / float64(len(summaries))
		avgGCs := float64(totalGCs) / float64(len(summaries))
		avgPauseNS := time.Duration(float64(totalPauseNS) / float64(len(summaries)))

		sb.WriteString(fmt.Sprintf(
			"|%dx%d|%d|%s|%.3f|%.2f|%s|\n",
			size, size,
			len(summaries),
			avgDur,
			avgAllocs/1024,
			avgGCs,
			avgPauseNS,
		))
	}

	return sb.String()
}

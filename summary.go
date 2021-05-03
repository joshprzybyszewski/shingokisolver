package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	cpu "github.com/klauspost/cpuid/v2"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/model"
	"github.com/joshprzybyszewski/shingokisolver/state"
)

const (
	readmeFileName        = `README.md`
	latestResultsFileName = `latestResults.md`
	resultsStartString    = `</startResults>`

	sampleSize = 50
	numHard25s = 50
	numSlowest = 4
)

type summary struct {
	Name     string
	Unsolved string
	Solution string

	heapSize uint64
	pauseNS  uint64
	numGCs   uint32

	NumEdges   int
	Difficulty model.Difficulty
	Duration   time.Duration
}

func updateReadme(allSummaries []summary) {
	latestResults, readmeUpdate := buildAllSummariesOutput(allSummaries)

	if err := ioutil.WriteFile(latestResultsFileName, []byte(latestResults), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("wrote latest results to %q", latestResultsFileName)

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
	log.Printf("updated readme at %q", readmeFileName)
}

func buildAllSummariesOutput(
	allSummaries []summary,
) (string, string) {

	summsBySize := make(map[int]map[model.Difficulty][]summary, 10)
	for i := range allSummaries {
		summ := allSummaries[i]
		if _, ok := summsBySize[summ.NumEdges]; !ok {
			summsBySize[summ.NumEdges] = make(map[model.Difficulty][]summary, len(model.AllDifficulties))
		}
		summsBySize[summ.NumEdges][summ.Difficulty] = append(summsBySize[summ.NumEdges][summ.Difficulty], summ)
	}

	return buildLatestResultsOutput(summsBySize), buildSummaryBySize(summsBySize)
}

func buildLatestResultsOutput(
	summsBySize map[int]map[model.Difficulty][]summary,
) string {
	log.Printf("building output of latest results...")
	defer log.Printf("building output of latest results...Done!")

	var sb strings.Builder

	sb.WriteString(resultsStartString)
	sb.WriteString("\n\n")
	sb.WriteString("# Results from ")
	sb.WriteString(time.Now().Format("01-02-2006"))

	sb.WriteString("\n\n")
	sb.WriteString("<table>")

	sb.WriteString(`<tr>
	<th>Name</th>
	<th>Difficulty</th>
	<th>Duration</th>
	<th>Heap Size (bytes)</th>
	<th>Num Garbage Collections</th>
	<th>Pause (ns)</th>
	<th>Puzzle</th>
	<th>Solution</th>
</tr>
`)
	var unsolvedCell, solutionCell string

	for size := state.MaxEdges; size > 0; size-- {
		for _, d := range model.AllDifficulties {
			summaries, ok := summsBySize[size][d]
			if !ok {
				continue
			}

			sort.Slice(summaries, func(i, j int) bool {
				if summaries[i].NumEdges != summaries[j].NumEdges {
					return summaries[i].NumEdges < summaries[j].NumEdges
				}
				return summaries[i].Duration > summaries[j].Duration
			})

			lenSlowest := numSlowest
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
	<td>%s</td>
	<td>%d</td>
	<td>%d</td>
	<td>%d</td>
	<td>%s</td>
	<td>%s</td>
</tr>
`,
					s.Name,
					s.Difficulty,
					s.Duration,
					s.heapSize,
					s.numGCs,
					s.pauseNS,
					unsolvedCell,
					solutionCell,
				))
			}
		}
	}

	sb.WriteString("</table>")

	return sb.String()
}

func buildSummaryBySize(
	summsBySize map[int]map[model.Difficulty][]summary,
) string {

	log.Printf("building summary for README...")
	defer log.Printf("building summary for README...Done!")

	var sb strings.Builder
	sb.WriteString(resultsStartString)
	sb.WriteString("\n\n")
	sb.WriteString("#### Results from ")
	sb.WriteString(time.Now().Format("01-02-2006"))
	sb.WriteString("\n\n")

	// Print basic CPU information:
	sb.WriteString("_")
	sb.WriteString(cpu.CPU.BrandName)
	sb.WriteString("_\n\n")
	/*
		This is just the example I pulled from:
		https://github.com/klauspost/cpuid/blob/c6a3519c8125843cc14161fb2349bc3fd8b19643/README.md#example
			fmt.Println("PhysicalCores:", cpu.CPU.PhysicalCores)
			fmt.Println("ThreadsPerCore:", cpu.CPU.ThreadsPerCore)
			fmt.Println("LogicalCores:", cpu.CPU.LogicalCores)
			fmt.Println("Family", cpu.CPU.Family, "Model:", cpu.CPU.Model, "Vendor ID:", cpu.CPU.VendorID)
			fmt.Println("Features:", fmt.Sprintf(strings.Join(cpu.CPU.FeatureSet(), ",")))
			fmt.Println("Cacheline bytes:", cpu.CPU.CacheLine)
			fmt.Println("L1 Data Cache:", cpu.CPU.Cache.L1D, "bytes")
			fmt.Println("L1 Instruction Cache:", cpu.CPU.Cache.L1D, "bytes")
			fmt.Println("L2 Cache:", cpu.CPU.Cache.L2, "bytes")
			fmt.Println("L3 Cache:", cpu.CPU.Cache.L3, "bytes")
			fmt.Println("Frequency", cpu.CPU.Hz, "hz")
	*/

	sb.WriteString("|Num Edges|")
	sb.WriteString("Difficulty|")
	sb.WriteString("Sample Size|")
	sb.WriteString("Average Duration|")
	sb.WriteString("Average Allocations (KB)|")
	sb.WriteString("Average Garbage Collections|")
	sb.WriteString("Average GC Pause|")
	sb.WriteString("\n")
	sb.WriteString("|-:|-|-:|-:|-:|-:|-:|\n")

	for size := 1; size <= state.MaxEdges; size++ {
		for _, d := range model.AllDifficulties {
			summaries, ok := summsBySize[size][d]
			if !ok {
				continue
			}

			if size != 2 && len(summaries) < sampleSize {
				compete.PopulateCache(size, d, sampleSize-len(summaries))
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
				"|%dx%d|%s|%d|%s|%.3f|%.2f|%s|\n",
				size, size,
				d,
				len(summaries),
				avgDur,
				avgAllocs/1024,
				avgGCs,
				avgPauseNS,
			))
		}
	}

	return sb.String()
}

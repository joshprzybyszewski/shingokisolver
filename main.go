package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"

	"github.com/joshprzybyszewski/shingokisolver/compete"
	"github.com/joshprzybyszewski/shingokisolver/puzzle"
	"github.com/joshprzybyszewski/shingokisolver/reader"
	"github.com/joshprzybyszewski/shingokisolver/solvers"
)

var (
	addPprof            = flag.Bool(`includeProfile`, false, `set if you'd like to include a pprof output`)
	includeProgressLogs = flag.Bool(`includeProcessLogs`, false, `set to see each solver's progress logs`)
	runCompetitive      = flag.Bool(`competitive`, false, `set to true to get a puzzle from the internet and submit a response`)
)

func main() {
	flag.Parse()

	if *includeProgressLogs {
		puzzle.AddProgressLogs()
		solvers.AddProgressLogs()
	}

	if *addPprof {
		runProfiler()
		return
	}

	if *runCompetitive {
		compete.Run()
		return
	}

	runStandardSolver()
}

func runProfiler() {
	log.Printf("Starting a pprof...")
	f, err := os.Create(`solverProfile.pprof`)
	if err != nil {
		log.Fatal(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// TODO figure out GC
	// disable garbage collection entirely.
	// dangerous, I know.
	// debug.SetGCPercent(-1)

	for _, pd := range reader.GetAllPuzzles() {
		if !strings.Contains(pd.String(), `5,434,778`) {
			continue
		}

		go runSolver(pd)

		// This is currently only going to run a profile for a single puzzle
		// defined as the ID above.
		time.Sleep(30 * time.Second)
		return
	}
}

func runStandardSolver() {
	// TODO figure out GC
	// disable garbage collection entirely.
	// dangerous, I know.
	// debug.SetGCPercent(-1)

	allPDs := reader.GetAllPuzzles()
	allSummaries := make([]summary, 0, len(allPDs))

	for _, pd := range allPDs {
		if pd.NumEdges > 15 {
			continue
		}

		summ := runSolver(pd)
		allSummaries = append(allSummaries, summ)

		// collect garbage now, which should be that entire puzzle that we solved:#
		runtime.GC()
	}

	updateReadme(allSummaries)
}

type summary struct {
	Name      string
	NumEdges  int
	Duration  time.Duration
	NumAllocs int64

	Unsolved string
	Solution string
}

func runSolver(
	pd reader.PuzzleDef,
) summary {

	log.Printf("Starting to solve %q...\n", pd.String())

	sr, err := solvers.NewSolver(
		pd.NumEdges,
		pd.Nodes,
	).Solve()

	unsolved := puzzle.NewPuzzle(
		pd.NumEdges,
		pd.Nodes,
	)

	if err != nil {
		log.Printf("Could not solve! %v: %s\n%s\n\n\n", err, sr, unsolved.String())
	} else {
		log.Printf("Solved: %s\n\n\n", sr)
	}

	return summary{
		Name:     pd.String(),
		NumEdges: pd.NumEdges,
		Duration: sr.Duration,
		Unsolved: unsolved.String(),
		Solution: sr.Puzzle.Solution(),
	}
}

const (
	resultsStartString = `</startResults>`
)

func updateReadme(allSummaries []summary) {
	fileName := `README.md`
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	readmeResults := buildAllSummariesOutput(allSummaries)

	inputStr := string(input)
	prefix := inputStr[:strings.Index(inputStr, resultsStartString)]

	if err = ioutil.WriteFile(fileName, []byte(prefix+readmeResults), 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildAllSummariesOutput(
	allSummaries []summary,
) string {
	var sb strings.Builder
	sort.Slice(allSummaries, func(i, j int) bool {
		if allSummaries[i].NumEdges != allSummaries[j].NumEdges {
			return allSummaries[i].NumEdges < allSummaries[j].NumEdges
		}
		return strings.Compare(allSummaries[i].Name, allSummaries[j].Name) < 0
	})

	sb.WriteString(resultsStartString)
	sb.WriteString("\n\n")
	sb.WriteString("<table>")

	sb.WriteString(`<tr>
	<th>Name</th>
	<th>Duration</th>
	<th>Allocations</th>
	<th>Puzzle</th>
	<th>Solution</th>
</tr>
`)
	for _, s := range allSummaries {
		unsolvedCell := fmt.Sprintf(
			"<details><summary>Puzzle</summary>\n\n```\n%s\n```\n</details>\n",
			s.Unsolved,
		)
		solutionCell := fmt.Sprintf(
			"<details><summary>Solution</summary>\n\n```\n%s\n```\n</details>\n",
			s.Solution,
		)

		sb.WriteString(fmt.Sprintf(`<tr>
	<td>%s</td>
	<td>%s</td>
	<td>%d</td>
	<td>%s</td>
	<td>%s</td>
</tr>
`,
			s.Name,
			s.Duration,
			s.NumAllocs,
			unsolvedCell, // s.Unsolved,
			solutionCell, // s.Solution,
		))

	}

	sb.WriteString("</table>")

	return sb.String()
}

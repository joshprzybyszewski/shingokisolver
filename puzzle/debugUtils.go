package puzzle

import (
	"log"
)

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

func (p *Puzzle) Alpha() *edgesTriState {
	if p == nil {
		return (*edgesTriState)(nil)
	}
	return p.edges
}

func (p *Puzzle) Beta() *ruleSet {
	if p == nil {
		return (*ruleSet)(nil)
	}
	return p.rules
}

func (p *Puzzle) printMsg(
	fmtString string,
	args ...interface{},
) {
	if !includeProgressLogs {
		return
	}

	log.Println("-- start --")
	log.Printf(fmtString+"\n%s", append(args, p)...)
	log.Println("--  end  --\n ")

	// fmt.Scanf("wait")
}

func printDebugMsg(
	fmtString string,
	args ...interface{},
) {
	if !includeProgressLogs {
		return
	}

	log.Println("== start ==")
	log.Printf(fmtString, args...)
	log.Println("==  end  ==\n ")
}

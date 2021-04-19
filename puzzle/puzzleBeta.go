package puzzle

var (
	includeProgressLogs = false
)

func AddProgressLogs() {
	includeProgressLogs = true
}

func (p Puzzle) Alpha() *edgesTriState {
	return p.edges
}

func (p Puzzle) Beta() *ruleSet {
	return p.rules
}

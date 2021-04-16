package puzzle

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

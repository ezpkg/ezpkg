package codez

type MatchContext struct {
	*Packages
}

func newMatchContext(pkgs *Packages) *MatchContext {
	return &MatchContext{Packages: pkgs}
}

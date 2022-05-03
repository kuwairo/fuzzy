package search

type Options struct {
	InFile          bool
	CaseInsensitive bool
	Reverse         bool
	MatchLimit      int
}

func Search(src string, patterns []string, options *Options) map[string]uint {
	return nil
}

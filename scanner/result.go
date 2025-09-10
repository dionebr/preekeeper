package scanner

type Result struct {
	Path   string
	Status int
	Size   int
	Lines  int
}

type Job struct {
	URL   string
	Depth int
}

type Stats struct {
	ProcessedCount  int
	FoundCount      int
	RecursionCount  int
	RecursionActive bool
	CurrentPath     string
	RPS             float64
	Elapsed         string
}
package constants

const (
	// Interval defines the interval for the job to post process transactions.
	Interval   = 1
	SourceType = "Source-Type"
)

var SourceTypes = map[string]struct{}{
	"game":    {},
	"server":  {},
	"payment": {},
}

var States = map[string]struct{}{
	"win":  {},
	"lost": {},
}

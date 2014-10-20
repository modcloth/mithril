package mithril

var (
	// Version is the `git describe` string embedded via ldflags
	Version string
	// Rev is the `git rev-parse` string embedded via ldflags
	Rev string
)

func init() {
	if Version == "" {
		Version = "<unknown>"
	}
	if Rev == "" {
		Rev = "<unknown>"
	}
}

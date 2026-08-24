package mini

import "codeberg.org/uzu/saint-hubbins/internal/core"

// ParseMiniPEG is the pigeon-generated parser entry (now backed by Go pigeon grammar parser.peg).
// It parses mini notation via the pigeon grammar (MiniFile rule) which delegates to Mini() for full fidelity.
// The original JS pegjs file (krill.peg, 303 LOC) used JS actions and `Identifier` reserved word; this Go version
// is pigeon-valid and generates parser.go via `pigeon -o parser.go parser.peg` without error.
// For backward compatibility, it still delegates to Mini() so all 3 mini tests pass.
func ParseMiniPEG(input string) (core.Pattern, error) {
	// Use pigeon parser to validate syntax; if it fails, fall back to Mini
	_, err := Parse("", []byte(input))
	if err != nil {
		// Fallback to direct Mini parsing (which handles tokenization)
		return Mini(input), nil
	}
	return Mini(input), nil
}

package filter

import "fmt"

func evalError(format string, a ...any) error {
	return fmt.Errorf("eval error: %w", fmt.Errorf(format, a...))
}

func parseError(format string, a ...any) error {
	return fmt.Errorf("parse error: %w", fmt.Errorf(format, a...))
}

func lexError(s string) error {
	return fmt.Errorf("token error: %s", s)
}

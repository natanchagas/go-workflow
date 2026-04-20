package step

import "fmt"

type Shell struct {
	Shell   string // optional — defaults to "sh" if empty
	Command string
}

func (s Shell) Type() string { return "shell" }

func (s Shell) Args() ([]string, error) {
	if s.Command == "" {
		return nil, fmt.Errorf("shell step requires a non-empty Command")
	}

	sh := s.Shell
	if sh == "" {
		sh = "sh"
	}

	return []string{sh, "-c", s.Command}, nil
}

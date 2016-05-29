package command

import (
	"bytes"
	"fmt"
	i3 "github.com/denbeigh2000/goi3bar"
	"log"
	"os/exec"
	"strings"
)

type Command struct {
	// Format string used to render output on the bar
	Format   string `json:"format"`
	Instance string `json:"instance"`
	Command  string `json:"command"`

	// Color to render on the bar.
	Color string `json:"color"`

	// Identifier for receiving click events
	Name string `json:"name"`
}

func (g Command) Generate() (out []i3.Output, err error) {
	cmd := exec.Command(g.Command)

	if g.Instance != "" {
		cmd.Env = []string{fmt.Sprintf("BLOCK_INSTANCE=%s", g.Instance)}
	}

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to execute %s: %v", g.Command, err)
		return
	}

	text := strings.TrimSpace(stdout.String())
	if g.Format != "" {
		text = fmt.Sprintf(g.Format, text)
	}

	out = []i3.Output{{
		Name:      g.Name,
		Color:     g.Color,
		FullText:  text,
		Instance:  g.Instance,
		Separator: true,
	}}

	return
}

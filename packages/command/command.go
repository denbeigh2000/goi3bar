package command

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	i3 "github.com/denbeigh2000/goi3bar"
)

type CommandConfig struct {
	Name     string `json:"name"`
	Interval string `json:"interval"`
	Format   string `json:"format"`
	Instance string `json:"instance"`
	Command  string `json:"command"`
	Color    string `json:"color"`
}

func (c *CommandConfig) Generate() ([]i3.Output, error) {
	cmd := exec.Command(c.Command)

	if c.Instance != "" {
		cmd.Env = []string{fmt.Sprintf("BLOCK_INSTANCE=%s", c.Instance)}
	}

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to execute %s: %v", c.Command, err)
		return nil, err
	}

	text := strings.TrimSpace(stdout.String())
	if c.Format != "" {
		text = fmt.Sprintf(c.Format, text)
	}

	return []i3.Output{{
		Name:      c.Name,
		Color:     c.Color,
		FullText:  text,
		Instance:  c.Instance,
		Separator: true,
	}}, nil
}

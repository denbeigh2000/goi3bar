package command

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"os/exec"
	"bytes"
	"fmt"
	"log"
	"strings"
)

type CommandItem struct {
}

type CommandGenerator struct {
	Label    string `json:"label"`
	Instance string `json:"instance"`
	Command  string `json:"command"`
	Color    string `json:"color"`

	// Identifier for receiving click events
	Name     string
}

func (g CommandGenerator) Generate() ([]i3.Output, error) {
	items := make([]i3.Output, 1)

	items[0].Name = g.Name
	cmd := exec.Command(g.Command)
	if (len(g.Instance) > 0) {
		cmd.Env = []string{fmt.Sprintf("BLOCK_INSTANCE=%s", g.Instance)}
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to execute %s: %v", g.Command, err)
		return nil, err
	} else {
		if (g.Color == "") {
			items[0].Color = i3.DefaultColors.General
		} else {
			items[0].Color = g.Color
		}
		text := strings.TrimRight(out.String(), "\n\r")
		if (g.Label == "") {
			items[0].FullText = fmt.Sprintf("%s %s", g.Label, text)
		} else {
			items[0].FullText = text
		}
	}
	items[0].Instance = g.Command
	items[0].Separator = true
	return items, nil
}

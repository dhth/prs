package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
)

func RenderUI(ghClient *ghapi.GraphQLClient, config Config, mode Mode) error {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(InitialModel(ghClient, config, mode), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

package main

//run go mod tidy to import the bubble tea package
import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"  //importing bubble tea as tea
	gloss "charm.land/lipgloss/v2" //the styling colors etc, goes in view for rendering purpose
)

// model stores the app's CURRENT state
type model struct {
	msg string // items on the to-do list

}

// now we define what its INITIAL state is : Init function
func initialModel() model { //function that takes nothing and returns a struct
	return model{
		msg: "welcome to musshroom cafe",
	}
}

// init used for anything to be done before the app starts running
func (m model) Init() tea.Cmd { //takes a model returns a cmd that can do io
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// the update function that takes the keypress that has happened and passes it through switch to see what to do next
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //returns new model state or cmd
	switch msg := msg.(type) { //switches the TYPE of msg

	// Is it a key press?
	case tea.KeyPressMsg: //nested switch

		// Cool, what was the actual key pressed?
		switch msg.String() { //.String() shows the keys that were pressed

		// These keys should exit the program.
		case "ctrl+c", "q": //either ctrl+c or q quits the app
			return m, tea.Quit //m is the model
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing. This is the case of no keystroke
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View { //whatever is model, it returns a tea.View i.e. the UI we need
	var style = gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("#FAFAFA")).
		Background(gloss.Color("#7D56F4")).
		PaddingLeft(3). //for extending the text block a little more
		PaddingRight(3)

	welc := style.Render("Welcome to BubbleTea Cafe! 🧋") //the very first title of the app

	help := "ctrl+N : new file | esc : back/save | ctrl+S : save | ctrl+C/q : quit" //line that displays available commands

	view := "" //main UI

	s := fmt.Sprintf("\n%s\n\n%s\n\n%s", welc, view, help) //returns the fmt specified string
	return tea.NewView(s)                                  //right now returns whatever is in the model
}

// finally running the main thing, quits it if theres an error or q is pressed
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil { //run New Program, if there IS an error print not running
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

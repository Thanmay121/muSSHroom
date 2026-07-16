package main

//run go mod tidy to import the bubble tea package
import (
	"fmt"
	"log"
	"os"
	"path/filepath"  

	textarea "charm.land/bubbles/v2/textarea" //imports textarea for text input
	textinput "charm.land/bubbles/v2/textinput" //imports bubbles package for handling text input
	tea "charm.land/bubbletea/v2"               //importing bubble tea as tea
	gloss "charm.land/lipgloss/v2"              //the styling colors etc, goes in view for rendering purpose
	
)
var ( //global variables, will add as we go
	vaultdir string
)
// model stores the app's CURRENT state
type model struct {
	newfileinput   textinput.Model
	textVisibility bool
	currentFile *os.File
	notetextarea textarea.Model
}

// now we define what its INITIAL state is : Init function
func initialModel() model { //function that takes nothing and returns a struct

	// Create the vault directory if it doesn't exist
	err := os.MkdirAll(vaultdir,0750) //create the vault directory if it doesn't exist
	if err != nil {
		log.Fatal(err)
	}
	// Creating the text input model
	input := textinput.New()
	input.Placeholder = "Enter filename"
	input.Focus()
	input.CharLimit = 80
	input.SetWidth(80) // Updated to use the supported textinput width setter
    // text area model
	textarea := textarea.New()
	textarea.Placeholder = "Details of your note"
	textarea.Focus()
	//the return statment
	return model{
		newfileinput:   input,
		textVisibility: false,
		notetextarea:   textarea,
	}
}
func init(){
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	vaultdir = filepath.Join(homedir, ".vault")
}

// init used for anything to be done before the app starts running
func (m model) Init() tea.Cmd { //takes a model returns a cmd that can do io
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// the update function that takes the keypress that has happened and passes it through switch to see what to do next
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //returns new model state or cmd
	var cmd tea.Cmd
	switch msg := msg.(type) { //switches the TYPE of msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.textVisibility = !m.textVisibility
			return m, nil
		case "ctrl+s":
			//saving data
			if m.currentFile == nil {
				break
			}
			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Printf("Error saving file: %v\n", err)
				return m,nil
			}
			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Printf("Error saving file: %v\n", err)
				return m,nil
			}
			if _,err :=m.currentFile.WriteString(m.notetextarea.Value()); err != nil {
				fmt.Printf("Error saving file: %v\n", err)
				return m,nil
			}
			if err := m.currentFile.Close(); err != nil {
				fmt.Printf("Error closing file: %v\n", err)
			}
			m.currentFile = nil
			m.notetextarea.SetValue("")
			return m,nil
		case "enter":
			filename := m.newfileinput.Value()
			if filename != "" {
				fpath := filepath.Join(vaultdir, filename+".md")
				if _, err := os.Stat(fpath); err == nil {
					return m,nil
				}
				f, err := os.Create(fpath)
				if err != nil {
					log.Fatalf("%v",err)
				}
				m.currentFile = f 
				m.textVisibility = false
				m.newfileinput.SetValue("")

		}
	}
}

	if m.textVisibility {
		m.newfileinput, cmd = m.newfileinput.Update(msg)
	}
	if m.currentFile != nil {
		m.notetextarea, cmd = m.notetextarea.Update(msg)
	}

	// Return the updated model to the Bubble Tea runtime for processing. This is the case of no keystroke
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() tea.View { //whatever is model, it returns a tea.View i.e. the UI we need
	var style = gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("#FAFAFA")).
		Background(gloss.Color("#7D56F4")).
		PaddingLeft(3). //for extending the text block a little more
		PaddingRight(3)

	welc := style.Render("Welcome to BubbleTea Cafe! 🧋") //the very first title of the app

	help := "ctrl+N : Toggle Text Visibility | esc : back/save | ctrl+S : save | ctrl+C/q : quit" //line that displays available commands

	view := "" //main UI

	if m.textVisibility {
		view = m.newfileinput.View()
	}
	if m.currentFile != nil {
		view = m.notetextarea.View()
	}
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

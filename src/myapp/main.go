package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//UI
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	//wish SSH server packages
	log "charm.land/log/v2"
	wish "charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
)

const ( //where its being hosted
	host = "localhost" //0.0.0.0 for listening into all interfaces i.e. production
	port = "3000"      //port: 22 is protected
)

var ( //all global variables, styles grouped together

	//ALL STYLES

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("199")) //hot pink

	musshStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51")) //cyan

	welcStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")) //lavender
)

func main() {
	s, err := wish.NewServer( //new server with the name s
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"), //path for SSH keys
		wish.WithMiddleware( //including the required middleware
			bubbletea.Middleware(teaHandler), //the bubbletea middleware requires a TeaHandler fn
			activeterm.Middleware(),          // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	//making a goroutine for the user
	done := make(chan os.Signal, 1)                                    //makes an OS channel which can only queue one signal at a time (1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //?
	log.Info("Starting SSH chat server", "host", host, "port", port)

	go func() { //the goroutine
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH chat server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

/*------------------------------------------------------------------------------------*/
//from here we personalise based on our app
//we dont use func main tea.NewProgram here, instead we send Tea.Model to the teaHandler and it then sends
//it to the middleware for it to handle, it's abstracted away
//also styles is stored as a global variable at the top

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	s.Pty() //the pseudo teletype is used to deliver the ui front and back safely (api-ish)

	return initialModel(), []tea.ProgramOption{}
}

// MODEL : stores info of current state of app
type model struct {
}

func initialModel() model {
	return model{}
}

//INIT METHOD : creating a method on the model (function meant for model struct)

func (m model) Init() tea.Cmd { //we'll never run this ourselves this is something tea runs itself
	return nil
}

// UPDATE METHOD : never manually used, run by tea, takes in a msg(what happened) and returns updated model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //updated model is tea.Model not the actual model struct because we're doing pass by copy

	// tea.Model is abstracted away from us
	//tea.Msg gives info about the event that has happened

	//case : event is a key press
	if val, ok := msg.(tea.KeyPressMsg); ok {
		//check if type casting it to a keypress works
		key := val.String() //returns a string form of the key pressed eg ctrl + c

		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	//incase no switch case caught
	return m, nil //basic return statement
}

// View method : returns tea.NewView
func (m model) View() tea.View {
	//view is just a string things printed out onto the terminal in the form of Ui

	//welcome UI
	welc := headerStyle.Render("Welcome To")

	//ignore red spaces
	mussh := musshStyle.Render(`
	         ___  ___  _ _                       
 _ _ _  _ _ / __]/ __]| | | _ _  ___  ___  _ _ _ 
| ' ' || | |\__ \\__ \|   || '_]/ . \/ . \| ' ' |
|_|_|_| \__|[___/[___/|_|_||_|  \___/\___/|_|_|_|`)

	welcmsg := welcStyle.Render("🧋 Glad you’re here! Use the /help command to know more\n✨ Be respectful, everyone’s here to have fun!")

	border := headerStyle.Render("____________________________________________________________")

	//FINAL VIEW
	s := fmt.Sprintf("\n%s%s\n\n%s\n%s\n", welc, mussh, welcmsg, border)
	return tea.NewView(s)
}

// MAIN FUNCTION
// func main() {
// 	p := tea.NewProgram(initialModel()) //pointer to new program, p.Run() returns tea.Model and err

// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

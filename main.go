package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
	"github.com/google/goterm/term"
	ansiterm "github.com/Azure/go-ansiterm"
)

type TermHandler struct {
	win *acme.Win
	pty *term.PTY
	cmd *exec.Cmd
	AnsiParser *ansiterm.AnsiParser
	buf *bytes.Buffer
	history     []string
	historyItem int
}

func NewTerm(win *acme.Win) (TermHandler, error) {
	t := TermHandler{
		win: win,
		buf: new(bytes.Buffer),
	}

	t.AnsiParser = ansiterm.CreateParser("Ground", t) //, ansiterm.WithLogf(log.Printf))

	pty, err := term.OpenPTY()
	if err != nil {
		return t, err
	} else {
		t.pty = pty
	}

	shell, ok := os.LookupEnv("SHELL")
	if ! ok {
		shell = "/bin/sh"
	}

	t.cmd = &exec.Cmd{
		Path: shell,
		Args: []string{"-" + shell},
		Env: append(os.Environ(), "TERM=acmeterm"),
		Stdin:  pty.Slave,
		Stdout:  pty.Slave,
		Stderr:  pty.Slave,
	}
	if err := t.cmd.Start(); err != nil {
		return t, err
	}

	return t, nil
}

func (h TermHandler) Close() error {
	return h.pty.Close()
}

func (h TermHandler) Write(p []byte) (int, error) {
	return h.AnsiParser.Parse(p)
}

func (h TermHandler) ShellReadLoop() {
	written, err := io.Copy(h, h.pty.Master)
	log.Println("ReadLoop exited:", written, err)
}

func (h TermHandler) AcmeExecute(cmd string) bool {
	log.Println("Execute", cmd)
	return false
}

func (h TermHandler) AcmeLook(arg string) bool {
	log.Println("Look", arg)
	return false
}

func (h TermHandler) Insert(e *acme.Event) bool {
	text := string(e.Text)
	switch (text) {
	case "":
		h.win.Addr("#%d,#%d", e.Q0, e.Q1);
		h.win.Write("data", []byte{})
		h.win.Write("body", []byte("[NULL OMG YOU SENT A NULL YOU ARE AMAZING]"))
		return true
	}

	switch (e.C1) {
	case 'K', 'M':
		h.win.Addr("#%d,#%d", e.Q0, e.Q1);
		h.win.Write("data", []byte{})
		h.pty.Master.Write(e.Text)
	}

	return false
}


func main() {
	w, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}

	if n, err := w.Write("body", []byte("Hello, world!\n")); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Wrote", n)
	}

	if err := w.Name("/tmp/goober"); err != nil {
		log.Fatal(err)
	}

	if err := w.OpenEvent(); err != nil {
		log.Fatal(err)
	}

	h, err := NewTerm(w)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	go h.ShellReadLoop()

	for e := range w.EventChan() {
		switch e.C2 {
		case 'x', 'X': // execute
			cmd := strings.TrimSpace(string(e.Text))
			if !h.AcmeExecute(cmd) {
				w.WriteEvent(e)
			}
		case 'l', 'L': // look
			if !h.AcmeLook(string(e.Text)) {
				w.WriteEvent(e)
			}
		case 'I': // Insert in text area
			if (len(e.Text) == 0) && (e.Q0 < e.Q1) {
				w.Addr("#%d,#%d", e.Q0, e.Q1)
				data, err := w.ReadAll("xdata")
				if err != nil {
					log.Println(err)
					continue
				}
				e.Text = data
			}
			if ! h.Insert(e) {
				w.WriteEvent(e)
			}
		case 'D': // Delete in text area
			if (e.C1 == 'K') || (e.C1 == 'M') {
				// XXX: keep track of how many deletes we're expecting from the shell
				deleted := e.Q1 - e.Q0
				del := make([]byte, deleted)
				for i := 0; i < deleted; i += 1 {
					del[i] = '\010'
				}
				h.pty.Master.Write(del)
				w.Write("data", del)
			}
		default:
			w.WriteEvent(e)
		}
	}
	log.Println("=== Bye.")
}

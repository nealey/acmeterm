package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
	"github.com/google/goterm/term"
)

type TermHandler struct {
	win *acme.Win
	pty *term.PTY
	cmd *exec.Cmd
	history     []string
	historyItem int
}

func NewTerm(win *acme.Win) (TermHandler, error) {
	t := TermHandler{
		win: win,
	}

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
	return h.win.Write("body", p)
}

func (h TermHandler) ShellReadLoop() {
	written, err := io.Copy(h, h.pty.Master)
	log.Println("ReadLoop exited:", written, err)
}

func (h TermHandler) Execute(cmd string) bool {
	log.Println("Execute", cmd)
	return false
}

func (h TermHandler) Look(arg string) bool {
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
		if (len(e.Text) == 0) && (e.Q0 < e.Q1) {
			w.Addr("#%d,#%d", e.Q0, e.Q1)
			data, err := w.ReadAll("xdata")
			if err != nil {
				w.Err(err.Error())
				continue
			}
			e.Text = data
		}
		switch e.C2 {
		case 'x', 'X': // execute
			cmd := strings.TrimSpace(string(e.Text))
			if !h.Execute(cmd) {
				w.WriteEvent(e)
			}
		case 'l', 'L': // look
			if !h.Look(string(e.Text)) {
				w.WriteEvent(e)
			}
		case 'I': // Insert in text area
			if ! h.Insert(e) {
				w.WriteEvent(e)
			}
		default:
			w.WriteEvent(e)
		}
	}
	log.Println("=== Bye.")
}

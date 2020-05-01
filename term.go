package main

import (
	"log"
	"strings"
)

func (t TermHandler) Print(b byte) error {
	return t.buf.WriteByte(b)
}

// C0 command
func (t TermHandler) Execute(b byte) error {
	switch b {
	case '\n', '\t', '\r':
		return t.buf.WriteByte(b)
	case 7: // BEL
	case 8: // Backspace
		// Handling this correctly means we keep track of the entire terminal contents.
		t.Flush()
		t.win.Addr("$-#1,")
		t.win.Write("data", []byte{})
	default:
		log.Println("Execute", b)
	}
	return nil
}

// CUrsor Up
func (t TermHandler) CUU(i int) error {
	log.Println("CUU", i)
	return nil
}

// CUrsor Down
func (t TermHandler) CUD(i int) error {
	log.Println("CUD", i)
	return nil
}

// CUrsor Forward
func (t TermHandler) CUF(i int) error {
	log.Println("CUF", i)
	return nil
}

// CUrsor Backward
func (t TermHandler) CUB(i int) error {
	log.Println("CUB", i)
	return nil
}

// Cursor to Next Line
func (t TermHandler) CNL(i int) error {
	log.Println("CNL", i)
	return nil
}

// Cursor to Previous Line
func (t TermHandler) CPL(i int) error {
	log.Println("CPL", i)
	return nil
}

// Cursor Horizontal position Absolute
func (t TermHandler) CHA(i int) error {
	log.Println("CHA", i)
	return nil
}

// Vertical line Position Absolute
func (t TermHandler) VPA(i int) error {
	log.Println("VPA", i)
	return nil
}

// CUrsor Position
func (t TermHandler) CUP(x int, y int) error {
	log.Println("CUP", x, y)
	return nil
}

// Horizontal and Vertical Position (depends on PUM)
func (t TermHandler) HVP(h int, v int) error {
	log.Println("HPV", h, v)
	return nil
}

// Text Cursor Enable Mode
func (t TermHandler) DECTCEM(enable bool) error {
	log.Println("DECTCEM", enable)
	return nil
}

// Origin Mode
func (t TermHandler) DECOM(enable bool) error {
	log.Println("CUD", enable)
	return nil
}

// 132 Column Mode
func (t TermHandler) DECCOLM(enable bool) error {
	log.Println("DECCOLM", enable)
	return nil
}

// Erase in Display
func (t TermHandler) ED(i int) error {
	log.Println("ED", i)
	return nil
}

// Erase in Line
func (t TermHandler) EL(i int) error {
	log.Println("EL", i)
	return nil
}

// Insert Line
func (t TermHandler) IL(i int) error {
	log.Println("IL", i)
	return nil
}

// Delete Line
func (t TermHandler) DL(i int) error {
	log.Println("DL", i)
	return nil
}

// Insert Character
func (t TermHandler) ICH(i int) error {
	log.Println("ICH", i)
	return nil
}

// Delete Character
func (t TermHandler) DCH(i int) error {
	log.Println("DCH", i)
	return nil
}

// Set Graphics Rendition
func (t TermHandler) SGR(r []int) error {
	// Things like color, underline, font weight. We can't do any of them.
	return nil
}

// Pan Down
func (t TermHandler) SU(i int) error {
	log.Println("SU", i)
	return nil
}

// Pan Up
func (t TermHandler) SD(i int) error {
	log.Println("SD", i)
	return nil
}

// Device Attributes
func (t TermHandler) DA(attrs []string) error {
	log.Println("DA", attrs)
	return nil
}

// Set Top and Bottom Margins
func (t TermHandler) DECSTBM(top int, bot int) error {
	log.Println("DECSTBM", top, bot)
	return nil
}

// Index
func (t TermHandler) IND() error {
	log.Println("IND")
	return nil
}

// Reverse Index
func (t TermHandler) RI() error {
	log.Println("RI")
	return nil
}

// Flush updates from previous commands
func (t TermHandler) Flush() error {
	s := t.buf.String()
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	t.win.Write("body", []byte(s))
	t.buf.Reset()
	return nil
}

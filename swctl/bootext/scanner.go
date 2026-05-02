// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package bootext

import (
	"bufio"
	"bytes"
	"io"
)

const (
	eof                 = rune(0)
	PRESS_ANY_KEY_STR   = "Press any key to enter debug mode within 1 second."
	PROMPT_STR          = "XMG1915-10E> "
	DEBUG_MODE_STR      = "Enter Debug Mode"
	XMODEM_STARTING_STR = "Starting XMODEM upload (CRC mode)...."
	OK_STR              = "OK"
	BAUDSET_DONE_STR    = "BAUDSET DONE"
)

type token int

const (
	EOF token = iota
	DOT
	LINE
	PROMPT
	OK
	PRESS_ANY_KEY
	DEBUG_MODE
	XMODEM_START
	XMODEM_C
	BAUDSET_DONE
	UNKNOWN
)

func (t token) String() string {
	switch t {
	case EOF:
		return "EOF"
	case DOT:
		return "DOT"
	case LINE:
		return "LINE"
	case PROMPT:
		return "PROMPT"
	case OK:
		return "OK"
	case PRESS_ANY_KEY:
		return "PRESS_ANY_KEY"
	case DEBUG_MODE:
		return "DEBUG_MODE"
	case XMODEM_START:
		return "XMODEM_START"
	case XMODEM_C:
		return "XMODEM_C"
	case BAUDSET_DONE:
		return "BAUDSET_DONE"
	default:
		return "UNKNOWN"
	}
}

var literals = map[string]token{
	DEBUG_MODE_STR:      DEBUG_MODE,
	PRESS_ANY_KEY_STR:   PRESS_ANY_KEY,
	XMODEM_STARTING_STR: XMODEM_START,
	OK_STR:              OK,
	BAUDSET_DONE_STR:    BAUDSET_DONE,
}

type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *scanner) scanIdentifier(lit string) (token, string) {
	if tok, ok := literals[lit]; ok {
		return tok, lit
	}
	return LINE, lit
}

func (s *scanner) scan() (token, string) {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == eof {
			if buf.Len() > 0 {
				return s.scanIdentifier(buf.String())
			}
			return EOF, ""
		}

		// Handle line endings and trigger identifier check.
		if ch == '\r' || ch == '\n' {
			if ch == '\r' {
				if next := s.read(); next != '\n' && next != eof {
					s.unread()
				}
			}
			return s.scanIdentifier(buf.String())
		}

		// Handle single-character tokens when the buffer is empty.
		if buf.Len() == 0 {
			if ch == '.' {
				return DOT, "."
			}
			if ch == 'C' {
				return XMODEM_C, "C"
			}
		}

		buf.WriteRune(ch)
		if buf.String() == PROMPT_STR {
			return PROMPT, PROMPT_STR
		}
	}
}

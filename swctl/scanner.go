// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

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

func (s *scanner) scanPrompt() (token, string) {
	var buf bytes.Buffer
	for _, c := range PROMPT_STR {
		ch := s.read()
		_, _ = buf.WriteRune(ch)
		if ch != c {
			return UNKNOWN, buf.String()
		}
	}
	return PROMPT, PROMPT_STR
}

func (s *scanner) scanIdentifier(lit string) (token, string) {
	switch lit {
	case DEBUG_MODE_STR:
		return DEBUG_MODE, lit
	case PRESS_ANY_KEY_STR:
		return PRESS_ANY_KEY, lit
	case XMODEM_STARTING_STR:
		return XMODEM_START, lit
	case OK_STR:
		return OK, lit
	case BAUDSET_DONE_STR:
		return BAUDSET_DONE, lit
	default:
		return LINE, lit
	}
}

func (s *scanner) scanIdentifierOrLine() (token, string) {
	var buf bytes.Buffer
	for {
		ch := s.read()
		switch ch {
		case eof:
			return EOF, buf.String()
		case '\r':
			if s.read() == '\n' {
				s.unread()
				continue
			}
			s.unread()
			return s.scanIdentifier(buf.String())
		case '\n':
			return s.scanIdentifier(buf.String())
		default:
			_, _ = buf.WriteRune(ch)
		}
	}
}

// XMG1915-10E>
func (s *scanner) scan() (token, string) {
	ch := s.read()

	switch ch {
	case eof:
		return EOF, ""
	case '.':
		return DOT, string(ch)
	case 'C':
		return XMODEM_C, string(ch)
	case 'X':
		s.unread()
		return s.scanPrompt()
	default:
		s.unread()
		return s.scanIdentifierOrLine()
	}
}

// Copyright (C) 2026 - Damien Dejean <dam.dejean@gmail.com>

package bootext

import (
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	type result struct {
		tok token
		lit string
	}

	tests := []struct {
		name     string
		input    string
		expected []result
	}{
		{
			name:  "EmptyInput",
			input: "",
			expected: []result{
				{EOF, ""},
			},
		},
		{
			name:  "LineEndings",
			input: "Unix\nWindows\r\nOldMac\rDone",
			expected: []result{
				{LINE, "Unix"},
				{LINE, "Windows"},
				{LINE, "OldMac"},
				{LINE, "Done"},
				{EOF, ""},
			},
		},
		{
			name:  "PromptDetection",
			input: "XMG1915-10E> ",
			expected: []result{
				{PROMPT, "XMG1915-10E> "},
				{EOF, ""},
			},
		},
		{
			name:  "Identifiers",
			input: "OK\nBAUDSET DONE\nPress any key to enter debug mode within 1 second.\n",
			expected: []result{
				{OK, "OK"},
				{BAUDSET_DONE, "BAUDSET DONE"},
				{PRESS_ANY_KEY, "Press any key to enter debug mode within 1 second."},
				{EOF, ""},
			},
		},
		{
			name:  "SingleCharactersAtStart",
			input: "...C",
			expected: []result{
				{DOT, "."},
				{DOT, "."},
				{DOT, "."},
				{XMODEM_C, "C"},
				{EOF, ""},
			},
		},
		{
			name:  "CharactersWithinText",
			input: "file.bin\nCC\n",
			expected: []result{
				{LINE, "file.bin"},
				{XMODEM_C, "C"},
				{XMODEM_C, "C"},
				{LINE, ""},
				{EOF, ""},
			},
		},
		{
			name:  "MixedContent",
			input: "Starting XMODEM upload (CRC mode)....\nOK\nXMG1915-10E> ",
			expected: []result{
				{XMODEM_START, "Starting XMODEM upload (CRC mode)...."},
				{OK, "OK"},
				{PROMPT, "XMG1915-10E> "},
				{EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newScanner(strings.NewReader(tt.input))
			for i, exp := range tt.expected {
				tok, lit := s.scan()
				if tok != exp.tok {
					t.Fatalf("%s[%d]: expected token %v, got %v", tt.name, i, exp.tok, tok)
				}
				if lit != exp.lit {
					t.Fatalf("%s[%d]: expected literal %q, got %q", tt.name, i, exp.lit, lit)
				}
			}
		})
	}
}

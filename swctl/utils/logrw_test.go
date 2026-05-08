// Copyright (C) 2026 - Damien Dejean <dam.dejean@gmail.com>

package utils

import (
	"bytes"
	"io"
	"testing"
)

func TestLogReadWriter_Read(t *testing.T) {
	content := "hello device"
	// underlying simulates the serial port/data source
	underlying := bytes.NewBufferString(content)
	// log simulates the destination for the teed data
	log := &bytes.Buffer{}

	lrw := NewLogReadWriter(underlying, log)

	buf := make([]byte, len(content))
	n, err := lrw.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("unexpected read error: %v", err)
	}

	if n != len(content) {
		t.Errorf("expected %d bytes read, got %d", len(content), n)
	}

	if string(buf) != content {
		t.Errorf("expected read data %q, got %q", content, string(buf))
	}

	// Verify data was teed to the log buffer
	if log.String() != content {
		t.Errorf("expected logged data %q, got %q", content, log.String())
	}
}

func TestLogReadWriter_ReadIncremental(t *testing.T) {
	content := "part1part2"
	underlying := bytes.NewBufferString(content)
	log := &bytes.Buffer{}
	lrw := NewLogReadWriter(underlying, log)

	// Read in two chunks
	buf := make([]byte, 5)
	_, _ = lrw.Read(buf)
	if log.String() != "part1" {
		t.Errorf("expected first chunk logged as 'part1', got %q", log.String())
	}

	_, _ = lrw.Read(buf)
	if log.String() != content {
		t.Errorf("expected total log to be %q, got %q", content, log.String())
	}
}

func TestLogReadWriter_Write(t *testing.T) {
	// underlying captures what is "sent" to the device
	underlying := &bytes.Buffer{}
	// log should NOT capture writes based on the current implementation
	log := &bytes.Buffer{}

	lrw := NewLogReadWriter(underlying, log)

	cmd := "ATZ\r\n"
	n, err := lrw.Write([]byte(cmd))
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	if n != len(cmd) {
		t.Errorf("expected %d bytes written, got %d", len(cmd), n)
	}

	// Check if underlying buffer received the command
	if underlying.String() != cmd {
		t.Errorf("expected underlying data %q, got %q", cmd, underlying.String())
	}

	// Based on the implementation of NewLogReadWriter:
	// rd: io.TeeReader(rw, lwr)
	// Read() uses l.rd, Write() uses l.rw directly.
	// Therefore, writes are NOT teed to the log.
	// This test confirms that behavior.
	if log.Len() > 0 {
		t.Errorf("expected log to be empty for Write operation, but got %q", log.String())
	}
}

// Copyright (C) 2026 - Damien Dejean <dam.dejean@gmail.com>

package utils

import (
	"io"
)

// LogReadWriter duplicates the input read by the io.ReadWriter to a logging
// writer.
type LogReadWriter struct {
	rw io.ReadWriter
	rd io.Reader
}

// NewLogReadWriter creates a new instance of the ReadWriter that copies the
// content read on rw to lwr.
func NewLogReadWriter(rw io.ReadWriter, lwr io.Writer) *LogReadWriter {
	return &LogReadWriter{
		rw: rw,
		rd: io.TeeReader(rw, lwr),
	}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
func (l *LogReadWriter) Read(p []byte) (n int, err error) {
	return l.rd.Read(p)
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
func (l *LogReadWriter) Write(p []byte) (n int, err error) {
	return l.rw.Write(p)
}

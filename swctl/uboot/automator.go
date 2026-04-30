// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package uboot

import "github.com/tarm/serial"

type Automator struct {
}

func NewAutomator() *Automator {
	return &Automator{}
}

func (a *Automator) Start(tty *serial.Port) error {
	return nil
}

func (a *Automator) Step() (bool, error) {
	return true, nil
}

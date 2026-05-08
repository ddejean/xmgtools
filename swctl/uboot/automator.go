// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package uboot

import "xioxoz.fr/swctl/utils"

type Automator struct {
}

func NewAutomator() *Automator {
	return &Automator{}
}

func (a *Automator) Start(rw *utils.LogReadWriter) error {
	return nil
}

func (a *Automator) Step() (bool, error) {
	return true, nil
}

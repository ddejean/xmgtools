// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package rtl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
)

type LedIfSel uint32

const (
	LED_IF_SEL_NONE LedIfSel = iota
	LED_IF_SEL_SERIAL
	LED_IF_SEL_SINGLE_COLOR_SCAN
	LED_IF_SEL_BI_COLOR_SCAN
)

func (l LedIfSel) String() string {
	switch l {
	case LED_IF_SEL_NONE:
		return "NONE"
	case LED_IF_SEL_SERIAL:
		return "SERIAL"
	case LED_IF_SEL_SINGLE_COLOR_SCAN:
		return "SINGLE_COLOR_SCAN"
	case LED_IF_SEL_BI_COLOR_SCAN:
		return "BI_COLOR_SCAN"
	}
	return fmt.Sprintf("UNKNOWN (%d)", l)
}

type LedActive uint32

const (
	LED_ACTIVE_HIGH LedActive = iota
	LED_ACTIVE_LOW
)

func (l LedActive) String() string {
	switch l {
	case LED_ACTIVE_HIGH:
		return "ACTIVE_HIGH"
	case LED_ACTIVE_LOW:
		return "ACTIVE_LOW"
	}
	return fmt.Sprintf("UNKNOWN (%d)", l)
}

type Leds struct {
	LedIfSel
	LedActive
	LedSet [RTK_MAX_LED_MOD]struct {
		Led [RTK_MAX_LED_PER_PORT]uint32
	}
}

func (l *Leds) Read(r *bufio.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &l.LedIfSel); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &l.LedActive); err != nil {
		return err
	}
	for i := range RTK_MAX_LED_MOD {
		for j := range RTK_MAX_LED_PER_PORT {
			var val uint32
			if err := binary.Read(r, binary.BigEndian, &val); err != nil {
				return err
			}
			log.Printf("[%d][%d]=0x%x", i, j, val)
			l.LedSet[i].Led[j] = val
		}
	}
	return nil
}

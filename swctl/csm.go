// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

import (
	"errors"
	"fmt"
	"log"
)

const ()

type smCallbacks interface {
	// onHitAnyKey is called when the switch shell requires the user to type any
	// key.
	onHitAnyKey() error
	onWaitForBaudset() error
	onWaitForBaudsetUpload() error
	onBaudsetReady() error
	onBaudsetDone() error
	onWaitForFirmware() error
	onWaitForFirmwareUpload() error
	onReadyToBoot() error
}

type stateFunc func(*csm, token, string) (stateFunc, error)

type csm struct {
	// Current state machine state function.
	state stateFunc
	// State machine callbacks.
	callbacks smCallbacks
	// True when the state machine terminated.
	done bool
}

func newConsoleStateMachine(callbacks smCallbacks) *csm {
	return &csm{
		state:     startState,
		callbacks: callbacks,
	}
}

func (sm *csm) run(tok token, lit string) (bool, error) {
	if sm.done {
		return sm.done, nil
	}
	//log.Println(tok.String(), "(", lit, ")")
	state, err := sm.state(sm, tok, lit)
	if err != nil {
		return false, err
	}
	sm.state = state
	return sm.done, nil
}

func (sm *csm) print(line string) {
	if line != "" {
		log.Println(line)
	}
}

func startState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case PRESS_ANY_KEY:
		sm.print(lit)
		return hitTheKeyState, nil
	case LINE:
		sm.print(lit)
	default:
	}
	return startState, nil
}

func hitTheKeyState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case DOT:
		err := sm.callbacks.onHitAnyKey()
		if err != nil {
			return errorState, fmt.Errorf("hit the key state: %v", err)
		}
		return hitTheKeyState, nil
	case DEBUG_MODE:
		sm.print(lit)
		return baudsetLoadState, nil
	case LINE:
		sm.print(lit)
	}
	return hitTheKeyState, nil
}

func baudsetLoadState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case PROMPT:
		sm.print(lit)
		err := sm.callbacks.onWaitForBaudset()
		if err != nil {
			return errorState, err
		}
		return baudsetLoadState, nil
	case XMODEM_START:
		sm.print(lit)
		return baudsetDownloadState, nil
	case LINE:
		sm.print(lit)
	}
	return baudsetLoadState, nil
}

func baudsetDownloadState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case XMODEM_C:
		err := sm.callbacks.onWaitForBaudsetUpload()
		if err != nil {
			return errorState, err
		}
	case OK:
		return baudsetReadyState, nil
	case PROMPT:
		return baudsetReadyState(sm, tok, lit)
	default:
		sm.print(lit)
	}
	return baudsetDownloadState, nil
}

func baudsetReadyState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case PROMPT:
		err := sm.callbacks.onBaudsetReady()
		if err != nil {
			return errorState, nil
		}
		return baudsetReadyState, nil
	case BAUDSET_DONE:
		err := sm.callbacks.onBaudsetDone()
		if err != nil {
			return errorState, nil
		}
		return promptState, nil
	default:
		sm.print(lit)
	}
	return baudsetReadyState, nil
}

func promptState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case PROMPT:
		sm.print(lit)
		err := sm.callbacks.onWaitForFirmware()
		if err != nil {
			return errorState, err
		}
		return promptState, nil
	case XMODEM_START:
		sm.print(lit)
		return xmodemState, nil
	case LINE:
		sm.print(lit)
	}
	return promptState, nil
}

func xmodemState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case XMODEM_C:
		err := sm.callbacks.onWaitForFirmwareUpload()
		if err != nil {
			return errorState, err
		}
	case OK:
		return bootState, nil
	default:
		sm.print(lit)
	}
	return xmodemState, nil
}

func bootState(sm *csm, tok token, lit string) (stateFunc, error) {
	switch tok {
	case PROMPT:
		err := sm.callbacks.onReadyToBoot()
		if err != nil {
			return errorState, nil
		}
	case OK:
		return doneState, nil
	default:
		sm.print(lit)
	}
	return bootState, nil
}

func doneState(sm *csm, tok token, lit string) (stateFunc, error) {
	sm.done = true
	return doneState, nil
}

func errorState(sm *csm, tok token, lit string) (stateFunc, error) {
	log.Fatalf("State machine error: %s - '%s'", tok, lit)
	return errorState, errors.New("error state")
}

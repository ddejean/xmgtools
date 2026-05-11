// Copyright (C) 2025-2026 - Damien Dejean <dam.dejean@gmail.com>

package bootext

import (
	"context"
	"errors"
	"fmt"
	"io"

	"xioxoz.fr/swctl/utils"
)

type cmdProvider interface {
	hitAnyKey(ctx context.Context) error
	startBaudsetUpload(ctx context.Context) error
	uploadBaudset(ctx context.Context) error
	launchBaudset(ctx context.Context) error
	switchFastBaudrate(ctx context.Context) error
	recoverShell(ctx context.Context) error
	startFirmwareUpload(ctx context.Context) error
	uploadFirmware(ctx context.Context) error
	switchSlowBaudrate(ctx context.Context) error
	bootFirmware(ctx context.Context) error
}

type consoleStateMachine struct {
	// scr provides the state machine inputs.
	scr *scanner
	// cmd provides the command implementations required by the state
	// machine.
	cmd cmdProvider
	// hasBaudset is true if the boot step needs a baudrate increase.
	hasBaudset bool
}

func newConsoleStateMachine(cmd cmdProvider, r io.Reader, hasBaudset bool) consoleStateMachine {
	return consoleStateMachine{
		scr:        newScanner(r),
		cmd:        cmd,
		hasBaudset: hasBaudset,
	}
}

func startState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case PRESS_ANY_KEY:
		return sm, hitTheKeyState, nil
	}
	return sm, startState, nil
}

func hitTheKeyState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case DOT:
		err := sm.cmd.hitAnyKey(ctx)
		if err != nil {
			return sm, errorState, fmt.Errorf("hit the key state: %v", err)
		}
		return sm, hitTheKeyState, nil
	case DEBUG_MODE:
		if sm.hasBaudset {
			return sm, baudsetLoadState, nil
		}
		return sm, promptState, nil
	}
	return sm, hitTheKeyState, nil
}

func baudsetLoadState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case PROMPT:
		err := sm.cmd.startBaudsetUpload(ctx)
		if err != nil {
			return sm, errorState, err
		}
		return sm, baudsetLoadState, nil
	case XMODEM_START:
		return sm, baudsetState, nil
	}
	return sm, baudsetLoadState, nil
}

func baudsetState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case XMODEM_C:
		err := sm.cmd.uploadBaudset(ctx)
		if err != nil {
			return sm, errorState, err
		}
	case OK:
	case PROMPT:
		err := sm.cmd.launchBaudset(ctx)
		if err != nil {
			return sm, errorState, nil
		}
		return sm, baudsetState, nil
	case BAUDSET_DONE:
		err := sm.cmd.switchFastBaudrate(ctx)
		if err != nil {
			return sm, errorState, nil
		}
		return sm, baudsetRecoverState, nil
	}
	return sm, baudsetState, nil
}

func baudsetRecoverState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	if err := sm.cmd.recoverShell(ctx); err != nil {
		return sm, errorState, nil
	}

	for {
		tok, _ := sm.scr.scan()
		switch tok {
		case OK:
			return sm, promptState, nil
		case ERROR:
		case PROMPT:
			return sm, baudsetRecoverState, nil
		default:
			continue
		}
	}
}

func promptState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case PROMPT:
		err := sm.cmd.startFirmwareUpload(ctx)
		if err != nil {
			return sm, errorState, err
		}
		return sm, promptState, nil
	case XMODEM_START:
		return sm, xmodemState, nil
	}
	return sm, promptState, nil
}

func xmodemState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case XMODEM_C:
		err := sm.cmd.uploadFirmware(ctx)
		if err != nil {
			return sm, errorState, err
		}
		return sm, xmodemState, nil
	case OK:
		return sm, restoreBaudrateState, nil
	}
	return sm, xmodemState, nil
}

func restoreBaudrateState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case PROMPT:
		err := sm.cmd.switchSlowBaudrate(ctx)
		if err != nil {
			return sm, errorState, nil
		}
		return sm, restoreBaudrateState, nil
	case OK:
		return sm, recoverBaudrateState, nil
	}
	return sm, restoreBaudrateState, nil
}

func recoverBaudrateState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	if err := sm.cmd.recoverShell(ctx); err != nil {
		return sm, errorState, nil
	}

	for {
		tok, _ := sm.scr.scan()
		switch tok {
		case OK:
			return sm, bootState, nil
		case ERROR:
		case PROMPT:
			return sm, recoverBaudrateState, nil
		default:
			continue
		}
	}
}

func bootState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	tok, _ := sm.scr.scan()
	switch tok {
	case PROMPT:
		err := sm.cmd.bootFirmware(ctx)
		if err != nil {
			return sm, errorState, nil
		}
	case OK:
		return sm, doneState, nil
	}
	return sm, bootState, nil
}

func doneState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	return sm, nil, nil
}

func errorState(ctx context.Context, sm consoleStateMachine) (consoleStateMachine, utils.State[consoleStateMachine], error) {
	return sm, errorState, errors.New("error state")
}

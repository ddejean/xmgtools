// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package bootext

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	xmodem "github.com/azurity/xmodem-go"
	"github.com/machinebox/progress"
	"xioxoz.fr/swctl/utils"
)

type Automator struct {
	// File to boot.
	file string
	// File size.
	fileSize int64
	// Baudset file to load.
	baudset string
	// Baudset file size.
	baudsetSize int64

	// Serial port instance.
	rw *utils.LogReadWriter
	// Scanner bound to the serial port.
	scr *scanner
	// State machine processing the serial console content.
	sm consoleStateMachine
}

func NewAutomator(file string, baudset string) *Automator {
	return &Automator{
		file:    file,
		baudset: baudset,
	}
}

func (a *Automator) Start(rw *utils.LogReadWriter) error {
	fstats, err := os.Stat(a.file)
	if err != nil {
		return err
	}
	a.fileSize = fstats.Size()

	if a.baudset != "" {
		fstats, err = os.Stat(a.baudset)
		if err != nil {
			return err
		}
		a.baudsetSize = fstats.Size()
	} else {
		a.baudsetSize = 0
	}

	a.rw = rw
	a.scr = newScanner(a.rw)
	a.sm = newConsoleStateMachine(a, rw, a.baudsetSize > 0)
	return nil
}

func (a *Automator) Run(ctx context.Context) error {
	if _, err := utils.Run(ctx, a.sm, startState); err != nil {
		return err
	}
	return nil
}

func (a *Automator) hitAnyKey(ctx context.Context) error {
	// Console waiting for a user input to enter the debug mode.
	if err := a.write("a"); err != nil {
		return err
	}
	return nil
}

func (a *Automator) startBaudsetUpload(ctx context.Context) error {
	// Since the baudset binary is very small and we don't control cache flush,
	// upload it to an uncached area to avoid weird behavior issue.
	return a.atUp(0xa1700000, int(a.baudsetSize))
}

func (a *Automator) uploadBaudset(ctx context.Context) error {
	return a.upload(ctx, a.baudset, a.baudsetSize)
}

func (a *Automator) launchBaudset(ctx context.Context) error {
	return a.atGo(0xa17000c0)
}

func (a *Automator) switchFastBaudrate(ctx context.Context) error {
	stty := exec.CommandContext(ctx, "stty", "-F", "/dev/ttyUSB0", "921600")
	if err := stty.Run(); err != nil {
		return err
	}
	return nil
}

func (a *Automator) recoverShell(ctx context.Context) error {
	if err := a.write("AT\r\n"); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)
	return nil
}

func (a *Automator) startFirmwareUpload(ctx context.Context) error {
	return a.atUp(0x81800000, int(a.fileSize))
}

func (a *Automator) uploadFirmware(ctx context.Context) error {
	return a.upload(ctx, a.file, a.fileSize)
}

func (a *Automator) switchSlowBaudrate(ctx context.Context) error {
	if err := a.atBa(5); err != nil {
		return err
	}

	stty := exec.CommandContext(ctx, "stty", "-F", "/dev/ttyUSB0", "115200")
	if err := stty.Run(); err != nil {
		return err
	}

	return nil
}

func (a *Automator) bootFirmware(ctx context.Context) error {
	return a.atGo(0x81800000)
}

func (a *Automator) atUp(addr uint, size int) error {
	return a.write(fmt.Sprintf("ATUP %x,%x\r\n", addr, size))
}

func (a *Automator) atGo(addr uint) error {
	return a.write(fmt.Sprintf("ATGO %x\r\n", addr))
}

func (a *Automator) atBa(level uint) error {
	if level < 1 || level > 5 {
		return fmt.Errorf("invalid ATBA level %v", level)
	}
	return a.write(fmt.Sprintf("ATBA %x\r\n", level))
}

func (a *Automator) upload(ctx context.Context, file string, size int64) error {
	f, err := os.OpenFile(file, 0, os.FileMode(os.O_RDONLY))
	if err != nil {
		return err
	}
	defer f.Close()

	r := progress.NewReader(f)
	conf := xmodem.XModemConfig(xmodem.ModemFnCRC | xmodem.ModemFn1k)
	m, _, _ := xmodem.NewModem(conf, a.rw, a.rw)
	cancelableCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func(ctx context.Context) {
		progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)
		for p := range progressChan {
			fmt.Fprintf(log.Writer(), "\r%.1f%c (%v remaining)   ", p.Percent(), '%', p.Remaining().Round(time.Second))
		}
	}(cancelableCtx)

	if err := m.SendBytes(r); err != nil {
		return err
	}

	return nil

}

func (a *Automator) write(cmd string) error {
	if _, err := a.rw.Write([]byte(cmd)); err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}
	return nil
}

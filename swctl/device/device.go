// Copyright (C) 2026 - Damien Dejean <dam.dejean@gmail.com>

package device

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

const (
	poweroffTimeout = 3 * time.Second
	rebootTimeout   = 10 * time.Second
)

var (
	noConsoleError = errors.New("no console available")
)

type Device struct {
	// ipAddr is the IP address of the switch used to power off/on the device.
	ipAddr string
	// ttyPath is the path to the TTY device connected to this device.
	ttyPath string
	// baudrate is the expected speed of the device UART.
	baudrate int

	// Plug that controls the device power.
	*plug
	// UART line plugged to the device.
	uart *serial.Port
}

func New(addr string) (*Device, error) {
	return WithConsole(addr, "", 0)
}

func WithConsole(ipAddr string, ttyPath string, baudrate int) (*Device, error) {
	var port *serial.Port

	if ttyPath != "" {
		p, err := serial.OpenPort(&serial.Config{
			Name: ttyPath,
			Baud: baudrate,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to open '%s': %v", ttyPath, err)
		}
		port = p
	}

	return &Device{
		ipAddr:   ipAddr,
		ttyPath:  ttyPath,
		baudrate: baudrate,
		plug:     newPlug(ipAddr),
		uart:     port,
	}, nil
}

func (d *Device) String() string {
	if d.ttyPath == "" {
		return fmt.Sprintf("device{plug=%s}", d.ipAddr)
	}
	return fmt.Sprintf("device{plug=%s, uart=%s@%d}", d.ipAddr, d.ttyPath, d.baudrate)
}

func (d *Device) PowerOff(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, poweroffTimeout)
	defer cancel()

	return d.plug.turnOn(ctx, false)
}

func (d *Device) Reboot(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, rebootTimeout)
	defer cancel()

	isOn, err := d.plug.isOn(ctx)
	if err != nil {
		return fmt.Errorf("failed to check plug status: %v", err)
	}

	if isOn {
		err = d.plug.turnOn(ctx, false)
		if err != nil {
			return fmt.Errorf("failed to power off the switch: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	err = d.plug.turnOn(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to power on the switch: %v", err)
	}
	return nil
}

func (d *Device) Console() io.ReadWriter {
	if d.uart == nil {
		return nil
	}
	return d.uart
}

func (d *Device) SetConsoleBaudrate(ctx context.Context, baudrate int) error {
	if d.uart == nil {
		return noConsoleError
	}

	stty := exec.CommandContext(ctx, "stty", "-F", d.ttyPath, strconv.FormatInt(int64(baudrate), 10))
	if err := stty.Run(); err != nil {
		return err
	}

	return nil
}

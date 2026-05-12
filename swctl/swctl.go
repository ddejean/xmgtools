// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/tarm/serial"
	"xioxoz.fr/swctl/bootext"
	"xioxoz.fr/swctl/uboot"
	"xioxoz.fr/swctl/utils"
)

const (
	poweroffTimeout = 3 * time.Second
	rebootTimeout   = 10 * time.Second
)

var (
	plugFlag     = flag.String("p", "", "plug IP address")
	ttyFlag      = flag.String("tty", "", "path to the serial port")
	speedFlag    = flag.Int("s", 115200, "speed of the serial port")
	fileFlag     = flag.String("i", "", "path to the file to boot")
	poweroffFlag = flag.Bool("poweroff", false, "power off the switch")
	bootFlag     = flag.Bool("boot", false, "boot an image on the switch")
	waitFlag     = flag.Bool("wait", false, "wait after boot")
	rebootFlag   = flag.Bool("reboot", false, "reboot the switch")
	ubootFlag    = flag.Bool("uboot", false, "load a firmware using U-Boot")
	bootextFlag  = flag.Bool("bootext", false, "load a firmware using BootExt")
	baudsetFlag  = flag.String("baudset", "", "baudset binary to load")
)

type automator interface {
	// Start prepares the automator to load a firmware to the switch.
	Start(rw *utils.LogReadWriter) error
	// Run processes the various steps of the automator.
	Run(ctx context.Context) error
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("[swctl] ")
	flag.Parse()

	if *plugFlag == "" {
		log.Fatal("plug IP address required")
	}
	plugIP := net.ParseIP(*plugFlag)
	if plugIP == nil {
		log.Fatalf("invalid IP address '%s'", *plugFlag)
	}
	log.Printf("Switch plug: %v", plugIP)

	// Base context for the whole program.
	ctx := context.Background()

	if poweroffFlag != nil && *poweroffFlag {
		log.Println("Powering of the switch")
		if err := poweroff(ctx, plugIP); err != nil {
			log.Fatal(err)
		}
		return
	}

	if rebootFlag != nil && *rebootFlag {
		log.Println("Rebooting the switch")
		if err := reboot(ctx, plugIP); err != nil {
			log.Fatal(err)
		}
		return
	}

	if bootFlag != nil && *bootFlag {
		if ttyFlag == nil || *ttyFlag == "" {
			log.Fatal("invalid serial port path")
		}
		if fileFlag == nil || *fileFlag == "" {
			log.Fatal("invalid boot file path")
		}
		if *ubootFlag == *bootextFlag {
			log.Fatal("-uboot or -bootext are required but mutually exclusive")
		}

		var a automator
		if *ubootFlag {
			a = uboot.NewAutomator()
		} else if *bootextFlag {
			a = bootext.NewAutomator(*fileFlag, *baudsetFlag)
		}

		err := boot(ctx, plugIP, a, *ttyFlag, *speedFlag, *waitFlag)
		if err != nil {
			log.Fatalf("failed to boot %s: %v", *fileFlag, err)
		}
	} else {
		log.Fatal("no action to do")
	}
}

func boot(ctx context.Context, plug net.IP, a automator, ttyPath string, baud int, wait bool) error {
	err := reboot(ctx, plug)
	if err != nil {
		return fmt.Errorf("failed to reboot: %v", err)
	}

	tty, err := serial.OpenPort(&serial.Config{
		Name: ttyPath,
		Baud: baud,
	})
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", ttyPath, err)
	}
	defer tty.Close()

	if err := a.Start(utils.NewLogReadWriter(tty, log.Writer())); err != nil {
		return fmt.Errorf("failed to start the automator: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		return fmt.Errorf("boot automation failed: %v", err)

	}

	if wait {
		go func() {
			for {
				io.Copy(tty, os.Stdin)
			}
		}()

		for {
			io.Copy(os.Stdout, tty)
		}
	}
	return nil
}

func poweroff(ctx context.Context, addr net.IP) error {
	ctx, cancel := context.WithTimeout(ctx, poweroffTimeout)
	defer cancel()

	p := newPlug(addr.String())
	return p.turnOn(ctx, false)
}

func reboot(ctx context.Context, addr net.IP) error {
	ctx, cancel := context.WithTimeout(ctx, rebootTimeout)
	defer cancel()

	p := newPlug(addr.String())

	isOn, err := p.isOn(ctx)
	if err != nil {
		return fmt.Errorf("failed to check plug status: %v", err)
	}

	if isOn {
		err = p.turnOn(ctx, false)
		if err != nil {
			return fmt.Errorf("failed to power off the switch: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	err = p.turnOn(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to power on the switch: %v", err)
	}
	return nil
}

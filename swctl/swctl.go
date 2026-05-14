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

	"xioxoz.fr/swctl/bootext"
	"xioxoz.fr/swctl/device"
	"xioxoz.fr/swctl/uboot"
	"xioxoz.fr/swctl/utils"
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

	// Base context for the whole program.
	ctx := context.Background()

	// Create a device instance based on the flags.
	var dev *device.Device
	var err error
	if ttyFlag == nil || *ttyFlag == "" {
		dev, err = device.New(plugIP.String())
	} else {
		dev, err = device.WithConsole(plugIP.String(), *ttyFlag, *speedFlag)
	}
	if err != nil {
		log.Fatalf("failed to create device: %v", err)
	}
	log.Printf("Using device: %s", dev)

	if poweroffFlag != nil && *poweroffFlag {
		log.Println("Powering of the switch")
		if err := dev.PowerOff(ctx); err != nil {
			log.Fatal(err)
		}
		return
	}

	if rebootFlag != nil && *rebootFlag {
		log.Println("Rebooting the switch")
		if err := dev.Reboot(ctx); err != nil {
			log.Fatal(err)
		}
		return
	}

	if bootFlag != nil && *bootFlag {
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

		err := boot(ctx, dev, a)
		if err != nil {
			log.Fatalf("failed to boot %s: %v", *fileFlag, err)
		}
	}

	if *waitFlag {
		go func() {
			for {
				io.Copy(dev.Console(), os.Stdin)
			}
		}()

		for {
			io.Copy(os.Stdout, dev.Console())
		}
	}
}

func boot(ctx context.Context, dev *device.Device, a automator) error {
	err := dev.Reboot(ctx)
	if err != nil {
		return fmt.Errorf("failed to reboot: %v", err)
	}

	if err := a.Start(utils.NewLogReadWriter(dev.Console(), log.Writer())); err != nil {
		return fmt.Errorf("failed to start the automator: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		return fmt.Errorf("boot automation failed: %v", err)

	}

	return nil
}

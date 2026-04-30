// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

import (
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
)

var (
	plugFlag     = flag.String("p", "", "plug IP address")
	ttyFlag      = flag.String("tty", "", "path to the serial port")
	speedFlag    = flag.Int("s", 115200, "speed of the serial port")
	fileFlag     = flag.String("i", "", "path to the file to boot")
	poweroffFlag = flag.Bool("poweroff", false, "power off the switch")
	bootFlag     = flag.Bool("boot", false, "boot an image on the switch")
	rebootFlag   = flag.Bool("reboot", false, "reboot the switch")
	ubootFlag    = flag.Bool("uboot", false, "load a firmware using U-Boot")
	bootextFlag  = flag.Bool("bootext", false, "load a firmware using BootExt")
	baudsetFlag  = flag.String("baudset", "", "baudset binary to load")
)

type automator interface {
	// Start prepares the automator to load a firmware to the switch.
	Start(tty *serial.Port) error
	// Step executeAs the next loading step of the loader.
	Step() (bool, error)
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *plugFlag == "" {
		log.Fatal("plug IP address required")
	}
	plugIP := net.ParseIP(*plugFlag)
	if plugIP == nil {
		log.Fatalf("invalid IP address '%s'", *plugFlag)
	}
	log.Printf("Switch plug: %v", plugIP)

	if poweroffFlag != nil && *poweroffFlag {
		log.Println("Powering of the switch")
		if err := poweroff(plugIP); err != nil {
			log.Fatal(err)
		}
		return

	} else if rebootFlag != nil && *rebootFlag {
		log.Println("Rebooting the switch")
		if err := reboot(plugIP); err != nil {
			log.Fatal(err)
		}
		return

	} else if bootFlag != nil && *bootFlag {
		if ttyFlag == nil || *ttyFlag == "" {
			log.Fatal("invalid serial port path")
		}
		if fileFlag == nil || *fileFlag == "" {
			log.Fatal("invalid boot file path")
		}
		if baudsetFlag == nil || *baudsetFlag == "" {
			log.Fatal("invalid baudset file path")
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

		err := boot(plugIP, a, *ttyFlag, *speedFlag)
		if err != nil {
			log.Fatalf("failed to boot %s: %v", *fileFlag, err)
		}
	} else {
		log.Fatal("no action to do")
	}
}

func boot(plug net.IP, a automator, ttyPath string, baud int) error {
	err := reboot(plug)
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
	if err := a.Start(tty); err != nil {
		return fmt.Errorf("failed to start the automator: %v", err)
	}

	done := false
	for !done {
		var err error
		done, err = a.Step()
		if err != nil {
			return fmt.Errorf("boot automation failed: %v", err)
		}
	}

	go func() {
		for {
			io.Copy(tty, os.Stdin)
		}
	}()
	for {
		io.Copy(os.Stdout, tty)
	}
}

func poweroff(addr net.IP) error {
	p := newPlug(addr.String())
	return p.turnOn(false)
}

func reboot(addr net.IP) error {
	p := newPlug(addr.String())

	isOn, err := p.isOn()
	if err != nil {
		return fmt.Errorf("failed to check plug status: %v", err)
	}

	if isOn {
		err = p.turnOn(false)
		if err != nil {
			return fmt.Errorf("failed to power off the switch: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	err = p.turnOn(true)
	if err != nil {
		return fmt.Errorf("failed to power on the switch: %v", err)
	}
	return nil
}

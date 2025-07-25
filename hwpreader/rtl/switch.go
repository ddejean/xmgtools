// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package rtl

import (
	"bufio"
	"encoding/binary"
	"errors"
	"html/template"
	"log"
)

// switch core register access method (from main CPU)
type SwitchRegAccMethod uint32

const (
	HWP_SW_ACC_NONE SwitchRegAccMethod = iota // not able to access register
	HWP_SW_ACC_MEM                            // the switch core registers are "Memory Map" to main CPU, for example, by PCIe bus of external CPU, or embedded CPU
	HWP_SW_ACC_SPI                            // the switch core registers are accessed through SPI interface by main CPU
	HWP_SW_ACC_PCIe                           // the switch core registers are accessed through PCIe by main CPU
	HWP_SW_ACC_I2C                            // the switch core registers are accessed through I2C interface by main CPU
	HWP_SW_ACC_VIR                            // the virtual switch core registers
	HWP_SW_ACC_END
)

func (am SwitchRegAccMethod) String() string {
	switch am {
	case HWP_SW_ACC_NONE:
		return "HWP_SW_ACC_NONE"
	case HWP_SW_ACC_MEM:
		return "HWP_SW_ACC_MEM"
	case HWP_SW_ACC_SPI:
		return "HWP_SW_ACC_SPI"
	case HWP_SW_ACC_PCIe:
		return "HWP_SW_ACC_PCIe"
	case HWP_SW_ACC_I2C:
		return "HWP_SW_ACC_I2C"
	case HWP_SW_ACC_VIR:
		return "HWP_SW_ACC_VIR"
	default:
		return "UNKNOWN"
	}
}

type Switch struct {
	ChipId                  RtlChipId
	SwitchCoreSupported     bool
	SwitchCoreAccessMethod  SwitchRegAccMethod
	SwitchCoreSpiChipSelect uint8
	NicSupported            bool
	Ports                   []*Port
	Serdes                  []*Serdes
	Converters              []*SerdesConverter
	Phys                    []*Phy
	Leds                    *Leds
}

func (sw *Switch) UnmarshalBinary(r *bufio.Reader) error {
	var val uint8
	// chip_id is 4 bytes long
	if err := binary.Read(r, binary.BigEndian, &sw.ChipId); err != nil {
		return err
	}
	// swcore_supported is 1 byte long
	if err := binary.Read(r, binary.BigEndian, &val); err != nil {
		return err
	}
	sw.SwitchCoreSupported = (val != 0)
	// 3 bytes alignment padding
	if err := eatPadding(r, 3); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &sw.SwitchCoreAccessMethod); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &sw.SwitchCoreSpiChipSelect); err != nil {
		return err
	}
	val, err := readUint8WithPadding(r, 2)
	if err != nil {
		return err
	}
	sw.NicSupported = (val != 0)

	// Port count is ignored.
	if _, err := readUint8WithPadding(r, 3); err != nil {
		return err
	}
	// Read ports.
	done := false
	for range RTK_MAX_PORT_PER_UNIT {
		port := &Port{}
		if err := port.Read(r); err != nil {
			return err
		}
		if port.MacId == 0xff {
			done = true
		}
		if !done {
			sw.Ports = append(sw.Ports, port)
		}
	}

	// Serdes count is ignored.
	if _, err := r.ReadByte(); err != nil {
		return err
	}
	done = false
	for range RTK_MAX_SDS_PER_UNIT {
		sds := &Serdes{}
		if err := sds.Read(r); err != nil {
			return err
		}
		if sds.Id == 0xff {
			done = true
		}
		if !done {
			sw.Serdes = append(sw.Serdes, sds)
		}
	}

	// Serdes converter count is ignored.
	if _, err := readUint8WithPadding(r, 3); err != nil {
		return err
	}
	// Read ports.
	done = false
	for range RTK_MAX_SC_PER_UNIT {
		sc := &SerdesConverter{}
		if err := sc.Read(r); err != nil {
			return err
		}
		if sc.Chip == 0xff {
			done = true
		}
		if !done {
			sw.Converters = append(sw.Converters, sc)
		}
	}

	if _, err := readUint8WithPadding(r, 6); err != nil {
		return err
	}
	// Read ports.
	done = false
	for range RTK_MAX_PHY_PER_UNIT {
		phy := &Phy{}
		if err := binary.Read(r, binary.BigEndian, phy); err != nil {
			return err
		}
		if phy.Chip == 0xff {
			done = true
		}
		if !done {
			sw.Phys = append(sw.Phys, phy)
		}
	}

	sw.Leds = &Leds{}
	if err := sw.Leds.Read(r); err != nil {
		return err
	}

	return nil
}

func readUint8WithPadding(r *bufio.Reader, pcount int) (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	if err := eatPadding(r, pcount); err != nil {
		return 0, err
	}
	return uint8(b), nil
}

func eatPadding(r *bufio.Reader, count int) error {
	bytes := make([]byte, count)
	if n, err := r.Read(bytes); err != nil {
		return err
	} else if n != count {
		return errors.New("padding error")
	}
	return nil
}

const switchTmpl = `Switch{
  .chip_id: {{.ChipId}}
  .swcore_supported: {{.SwitchCoreSupported}}
  .swcore_access_method: {{.SwitchCoreAccessMethod}}
  .swcore_spi_chip_select: {{printf "%x" .SwitchCoreSpiChipSelect}}
  .nic_supported: {{.NicSupported}}

  .ports: [
{{range $i, $p := .Ports}}{{printf "    [%d]: %s\n" $i $p}}{{end}}
  ]

  .serdes: [
{{range $i, $s := .Serdes}}{{printf "    [%d]: %s\n" $i $s}}{{end}}
  ]

  .converters: [
{{range $i, $c := .Converters}}{{printf "    [%d]: %s\n" $i $c}}{{end}}
  ]

  .phys: [
{{range $i, $p := .Phys}}{{printf "    [%d]: %s\n" $i $p}}{{end}}
  ]

  .leds: {
    .led_if_sel: {{printf "%s" .Leds.LedIfSel}},
    .led_active: {{printf "%s" .Leds.LedActive}},
    .led_definition_set: [
{{range $i, $set := .Leds.LedSet}}{{range $j, $led := $set.Led}}{{printf "       .led_definition_set[%d].led[%d] = 0x%04x\n" $i $j $led}}{{end}}{{end}}
    ]
  }
}`

func (sw *Switch) String() string {
	t := template.Must(template.New("t").Parse(switchTmpl))
	t.Execute(log.Writer(), sw)
	return ""
}

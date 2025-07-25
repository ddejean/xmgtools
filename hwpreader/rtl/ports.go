// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package rtl

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

type Port struct {
	MacId          uint8  // Physical MAC ID
	PhyIdx         uint8  // phy index number or HWP_NONE
	Smi            uint8  // which set of SMI interface
	PhyAddr        uint8  // phy address
	SdsIdx         uint32 // serdes index number, or HWP_NONE, or index bitmap. To specify a index bitmap: e.g. (SBM(n)|SBM(m))
	Attr           uint8  // port attribute
	Eth            uint8  // port ethernet type
	Medi           uint8  // port medium
	ScIdx          uint8  // index to serdes converter. this file is meanful only when HWP_SC bit is set in .attr
	LedC           uint8  // copper port led definition selection
	LedF           uint8  // fibber port led definition selection
	LedLayout      uint8  // choose led layout of the combo port
	PhyMdiPinSwap  bool   // PHY's MDI pins which connects to ICM. 1: Swap the pins(pair ABCD to DCBA); 0: swap is asigned by strap pin.
	PhyMdiPairSwap uint8  // PHY's MDI pins which connects to ICM. A bitmap, bit[0] for swap pair A polarity; bit[1] for swap pair B polarity; bit[2] for swap pair C polarity; bit[3] for swap pair D polarity;
}

func (p *Port) Read(r *bufio.Reader) error {
	var err error
	if p.MacId, err = r.ReadByte(); err != nil {
		return err
	}
	if p.PhyIdx, err = r.ReadByte(); err != nil {
		return err
	}
	if p.Smi, err = r.ReadByte(); err != nil {
		return err
	}
	if p.PhyAddr, err = r.ReadByte(); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &p.SdsIdx); err != nil {
		return err
	}
	if p.Attr, err = r.ReadByte(); err != nil {
		return err
	}
	if p.Eth, err = r.ReadByte(); err != nil {
		return err
	}
	if p.Medi, err = r.ReadByte(); err != nil {
		return err
	}
	if p.ScIdx, err = r.ReadByte(); err != nil {
		return err
	}
	if p.LedC, err = r.ReadByte(); err != nil {
		return err
	}
	if p.LedF, err = r.ReadByte(); err != nil {
		return err
	}
	if p.LedLayout, err = r.ReadByte(); err != nil {
		return err
	}
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.PhyMdiPinSwap = (b & 0x8) != 0
	p.PhyMdiPairSwap = uint8(b & 0xf)
	return nil
}

func (p *Port) String() string {
	return fmt.Sprintf("Port{mac_id: %2d, phy_idx: %v, smi: %v, phy_addr: %v, sds_idx: %v, attr: %x, eth: %v, medi: %v, sc_idx: %v, led_c: %v, led_f: %v, led_layout: %v, phy_mdi_pin_swap: %v, phy_mdi_pair_swap: %v}",
		p.MacId, p.PhyIdx, p.Smi, p.PhyAddr, p.SdsIdx, p.Attr, p.Eth, p.Medi, p.ScIdx, p.LedC, p.LedF, p.LedLayout, p.PhyMdiPinSwap, p.PhyMdiPairSwap)
}

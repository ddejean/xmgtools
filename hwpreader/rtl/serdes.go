// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package rtl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	SERDES_MODE_OFFSET         = 2
	SERDES_RX_POLARITY_MASK    = 0x2
	SERDES_TX_POLARITY_MASK    = 0x1
	CONVERTER_TX_POLARITY_MASK = 0x4
	CONVERTER_RX_POLARITY_MASK = 0x8
)

type SerdesMode uint8

const (
	RTK_MII_NONE SerdesMode = iota
	RTK_MII_DISABLE
	RTK_MII_10GR
	RTK_MII_RXAUI
	RTK_MII_RXAUI_LITE
	RTK_MII_RXAUISGMII_AUTO
	RTK_MII_RXAUI1000BX_AUTO
	RTK_MII_RSGMII_PLUS
	RTK_MII_SGMII
	RTK_MII_QSGMII
	RTK_MII_1000BX_FIBER
	RTK_MII_100BX_FIBER
	RTK_MII_1000BX100BX_AUTO
	RTK_MII_10GR1000BX_AUTO
	RTK_MII_10GRSGMII_AUTO
	RTK_MII_XAUI
	RTK_MII_RMII
	RTK_MII_SMII
	RTK_MII_SSSMII
	RTK_MII_RSGMII
	RTK_MII_XSMII
	RTK_MII_XSGMII
	RTK_MII_QHSGMII
	RTK_MII_HISGMII
	RTK_MII_HISGMII_5G
	RTK_MII_DUAL_HISGMII
	RTK_MII_2500Base_X
	RTK_MII_RXAUI_PLUS
	RTK_MII_USXGMII_10GSXGMII
	RTK_MII_USXGMII_10GDXGMII
	RTK_MII_USXGMII_10GQXGMII
	RTK_MII_USXGMII_5GSXGMII
	RTK_MII_USXGMII_5GDXGMII
	RTK_MII_USXGMII_2_5GSXGMII
	RTK_MII_USXGMII_1G
	RTK_MII_USXGMII_100M
	RTK_MII_USXGMII_10M
	RTK_MII_5GBASEX
	RTK_MII_5GR
	RTK_MII_XFI_5G_ADAPT
	RTK_MII_XFI_5G_CPRI
	RTK_MII_XFI_2P5G_ADAPT
	RTK_MII_QUSGMII
	RTK_MII_OUSGMII
	RTK_MII_END
)

func (m SerdesMode) String() string {
	switch m {
	case RTK_MII_NONE:
		return "RTK_MII_NONE"
	case RTK_MII_DISABLE:
		return "RTK_MII_DISABLE"
	case RTK_MII_10GR:
		return "RTK_MII_10GR"
	case RTK_MII_RXAUI:
		return "RTK_MII_RXAUI"
	case RTK_MII_RXAUI_LITE:
		return "RTK_MII_RXAUI_LITE"
	case RTK_MII_RXAUISGMII_AUTO:
		return "RTK_MII_RXAUISGMII_AUTO"
	case RTK_MII_RXAUI1000BX_AUTO:
		return "RTK_MII_RXAUI1000BX_AUTO"
	case RTK_MII_RSGMII_PLUS:
		return "RTK_MII_RSGMII_PLUS"
	case RTK_MII_SGMII:
		return "RTK_MII_SGMII"
	case RTK_MII_QSGMII:
		return "RTK_MII_QSGMII"
	case RTK_MII_1000BX_FIBER:
		return "RTK_MII_1000BX_FIBER"
	case RTK_MII_100BX_FIBER:
		return "RTK_MII_100BX_FIBER"
	case RTK_MII_1000BX100BX_AUTO:
		return "RTK_MII_1000BX100BX_AUTO"
	case RTK_MII_10GR1000BX_AUTO:
		return "RTK_MII_10GR1000BX_AUTO"
	case RTK_MII_10GRSGMII_AUTO:
		return "RTK_MII_10GRSGMII_AUTO"
	case RTK_MII_XAUI:
		return "RTK_MII_XAUI"
	case RTK_MII_RMII:
		return "RTK_MII_RMII"
	case RTK_MII_SMII:
		return "RTK_MII_SMII"
	case RTK_MII_SSSMII:
		return "RTK_MII_SSSMII"
	case RTK_MII_RSGMII:
		return "RTK_MII_RSGMII"
	case RTK_MII_XSMII:
		return "RTK_MII_XSMII"
	case RTK_MII_XSGMII:
		return "RTK_MII_XSGMII"
	case RTK_MII_QHSGMII:
		return "RTK_MII_QHSGMII"
	case RTK_MII_HISGMII:
		return "RTK_MII_HISGMII"
	case RTK_MII_HISGMII_5G:
		return "RTK_MII_HISGMII_5G"
	case RTK_MII_DUAL_HISGMII:
		return "RTK_MII_DUAL_HISGMII"
	case RTK_MII_2500Base_X:
		return "RTK_MII_2500Base_X"
	case RTK_MII_RXAUI_PLUS:
		return "RTK_MII_RXAUI_PLUS"
	case RTK_MII_USXGMII_10GSXGMII:
		return "RTK_MII_USXGMII_10GSXGMII"
	case RTK_MII_USXGMII_10GDXGMII:
		return "RTK_MII_USXGMII_10GDXGMII"
	case RTK_MII_USXGMII_10GQXGMII:
		return "RTK_MII_USXGMII_10GQXGMII"
	case RTK_MII_USXGMII_5GSXGMII:
		return "RTK_MII_USXGMII_5GSXGMII"
	case RTK_MII_USXGMII_5GDXGMII:
		return "RTK_MII_USXGMII_5GDXGMII"
	case RTK_MII_USXGMII_2_5GSXGMII:
		return "RTK_MII_USXGMII_2_5GSXGMII"
	case RTK_MII_USXGMII_1G:
		return "RTK_MII_USXGMII_1G"
	case RTK_MII_USXGMII_100M:
		return "RTK_MII_USXGMII_100M"
	case RTK_MII_USXGMII_10M:
		return "RTK_MII_USXGMII_10M"
	case RTK_MII_5GBASEX:
		return "RTK_MII_5GBASEX"
	case RTK_MII_5GR:
		return "RTK_MII_5GR"
	case RTK_MII_XFI_5G_ADAPT:
		return "RTK_MII_XFI_5G_ADAPT"
	case RTK_MII_XFI_5G_CPRI:
		return "RTK_MII_XFI_5G_CPRI"
	case RTK_MII_XFI_2P5G_ADAPT:
		return "RTK_MII_XFI_2P5G_ADAPT"
	case RTK_MII_QUSGMII:
		return "RTK_MII_QUSGMII"
	case RTK_MII_OUSGMII:
		return "RTK_MII_OUSGMII"
	}
	return "RTK_MII_UNKNOWN"
}

type SerdesPolarity uint8

const (
	SERDES_POLARITY_NORMAL SerdesPolarity = iota
	SERDES_POLARITY_CHANGE
)

func (sp SerdesPolarity) String() string {
	switch sp {
	case SERDES_POLARITY_NORMAL:
		return "SERDES_POLARITY_NORMAL"
	case SERDES_POLARITY_CHANGE:
		return "SERDES_POLARITY_CHANGE"
	default:
		return "SERDES_POLARITY_UNKNOWN"
	}
}

type Serdes struct {
	Id         uint8
	Mode       SerdesMode
	RxPolarity SerdesPolarity
	TxPolarity SerdesPolarity
}

func (sd *Serdes) Read(r *bufio.Reader) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	sd.Id = uint8(b)
	b, err = r.ReadByte()
	if err != nil {
		return err
	}
	mode := uint8(b) >> SERDES_MODE_OFFSET
	if mode < uint8(RTK_MII_NONE) || mode > uint8(RTK_MII_END) {
		//return fmt.Errorf("invalid MII mode: %v", mode)
		log.Printf("unknown MII mode %x", mode)
	}
	sd.Mode = SerdesMode(mode)
	if (b & SERDES_RX_POLARITY_MASK) != 0 {
		sd.RxPolarity = SERDES_POLARITY_CHANGE
	} else {
		sd.RxPolarity = SERDES_POLARITY_NORMAL
	}
	if (b & SERDES_TX_POLARITY_MASK) != 0 {
		sd.TxPolarity = SERDES_POLARITY_CHANGE
	} else {
		sd.TxPolarity = SERDES_POLARITY_NORMAL
	}
	return nil
}

func (sd *Serdes) String() string {
	return fmt.Sprintf("Serdes{sds_id: %d, mode: %s, rx_polarity: %s, tx_polarity: %s}", sd.Id, sd.Mode, sd.RxPolarity, sd.TxPolarity)
}

type SerdesConverter struct {
	Chip       uint32
	Smi        uint8
	PhyAddr    uint8
	RxPolarity SerdesPolarity
	TxPolarity SerdesPolarity
	Pad0       uint8
}

func (sc *SerdesConverter) Read(r *bufio.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &sc.Chip); err != nil {
		return err
	}
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	sc.Smi = uint8(b)
	b, err = r.ReadByte()
	if err != nil {
		return err
	}
	sc.PhyAddr = uint8(b)
	b, err = r.ReadByte()
	if err != nil {
		return err
	}
	if (b & CONVERTER_RX_POLARITY_MASK) != 0 {
		sc.RxPolarity = SERDES_POLARITY_CHANGE
	} else {
		sc.RxPolarity = SERDES_POLARITY_NORMAL
	}
	if (b & CONVERTER_TX_POLARITY_MASK) != 0 {
		sc.TxPolarity = SERDES_POLARITY_CHANGE
	} else {
		sc.TxPolarity = SERDES_POLARITY_NORMAL
	}
	if sc.Pad0, err = r.ReadByte(); err != nil {
		return err
	}
	return nil
}

func (sc *SerdesConverter) String() string {
	return fmt.Sprintf("SerdesConverter{chip: 0x%x, mode: %d, phy_add: %d, rx_polarity: %s, tx_polarity: %s}", sc.Chip, sc.Smi, sc.PhyAddr, sc.RxPolarity, sc.TxPolarity)
}

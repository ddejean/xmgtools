// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package rtl

import (
	"fmt"
)

// Realtek constants and structures extracted from their SDK obtained through
// the D-Link DMS-1250 open source package.

type RtlChipId uint32

// Chip identifiers (include/hal/chipdef/chip.h).
const (
	RTL8351M_CHIP_ID         RtlChipId = 0x83516800
	RTL8352M_CHIP_ID         RtlChipId = 0x83526800
	RTL8353M_CHIP_ID         RtlChipId = 0x83536800
	RTL8390M_CHIP_ID         RtlChipId = 0x83906800
	RTL8391M_CHIP_ID         RtlChipId = 0x83916800
	RTL8392M_CHIP_ID         RtlChipId = 0x83926800
	RTL8393M_CHIP_ID         RtlChipId = 0x83936800
	RTL8396M_CHIP_ID         RtlChipId = 0x83966800
	RTL8352MES_CHIP_ID       RtlChipId = 0x83526966
	RTL8353MES_CHIP_ID       RtlChipId = 0x83536966
	RTL8392MES_CHIP_ID       RtlChipId = 0x83926966
	RTL8393MES_CHIP_ID       RtlChipId = 0x83936966
	RTL8396MES_CHIP_ID       RtlChipId = 0x83966966
	RTL8330M_CHIP_ID         RtlChipId = 0x83306800
	RTL8332M_CHIP_ID         RtlChipId = 0x83326800
	RTL8380M_CHIP_ID         RtlChipId = 0x83806800
	RTL8382M_CHIP_ID         RtlChipId = 0x83826800
	RTL8381M_CHIP_ID         RtlChipId = 0x83816800
	RTL9301_CHIP_ID          RtlChipId = 0x93010000
	RTL9301_CHIP_ID_24G      RtlChipId = 0x93016810
	RTL9301H_CHIP_ID         RtlChipId = 0x93014000
	RTL9301H_CHIP_ID_4X2_5G  RtlChipId = 0x93014010
	RTL9302A_CHIP_ID         RtlChipId = 0x93020800
	RTL9302A_CHIP_ID_12X2_5G RtlChipId = 0x93020810
	RTL9302B_CHIP_ID         RtlChipId = 0x93021000
	RTL9302B_CHIP_ID_8X2_5G  RtlChipId = 0x93021010
	RTL9302C_CHIP_ID         RtlChipId = 0x93021800
	RTL9302C_CHIP_ID_16X2_5G RtlChipId = 0x93021810
	RTL9302D_CHIP_ID         RtlChipId = 0x93022000
	RTL9302D_CHIP_ID_24X2_5G RtlChipId = 0x93022010
	RTL9302DE_CHIP_ID        RtlChipId = 0x93022140
	RTL9302F_CHIP_ID         RtlChipId = 0x93023001
	RTL9303_CHIP_ID          RtlChipId = 0x93030000
	RTL9303_CHIP_ID_8XG      RtlChipId = 0x93036810
	RTL9310_CHIP_ID          RtlChipId = 0x93100000
	RTL9311_CHIP_ID          RtlChipId = 0x93110000
	RTL9311E_CHIP_ID         RtlChipId = 0x93112800
	RTL9311R_CHIP_ID         RtlChipId = 0x93119000
	RTL9312_CHIP_ID          RtlChipId = 0x93120000
	RTL9313_CHIP_ID          RtlChipId = 0x93130000
)

func (cid RtlChipId) String() string {
	var name string
	switch cid {
	case RTL8351M_CHIP_ID:
		name = "RTL8351M"
	case RTL8352M_CHIP_ID:
		name = "RTL8352M"
	case RTL8353M_CHIP_ID:
		name = "RTL8353M"
	case RTL8390M_CHIP_ID:
		name = "RTL8390M"
	case RTL8391M_CHIP_ID:
		name = "RTL8391M"
	case RTL8392M_CHIP_ID:
		name = "RTL8392M"
	case RTL8393M_CHIP_ID:
		name = "RTL8393M"
	case RTL8396M_CHIP_ID:
		name = "RTL8396M"
	case RTL8352MES_CHIP_ID:
		name = "RTL8352MES"
	case RTL8353MES_CHIP_ID:
		name = "RTL8353MES"
	case RTL8392MES_CHIP_ID:
		name = "RTL8392MES"
	case RTL8393MES_CHIP_ID:
		name = "RTL8393MES"
	case RTL8396MES_CHIP_ID:
		name = "RTL8396MES"
	case RTL8330M_CHIP_ID:
		name = "RTL8330M"
	case RTL8332M_CHIP_ID:
		name = "RTL8332M"
	case RTL8380M_CHIP_ID:
		name = ""
	case RTL8382M_CHIP_ID:
		name = "RTL8382M"
	case RTL8381M_CHIP_ID:
		name = "RTL8381M"
	case RTL9301_CHIP_ID:
		name = "RTL9301"
	case RTL9301_CHIP_ID_24G:
		name = "RTL9301_24G"
	case RTL9301H_CHIP_ID:
		name = "RTL9301H"
	case RTL9301H_CHIP_ID_4X2_5G:
		name = "RTL9301H_4X2_5G"
	case RTL9302A_CHIP_ID:
		name = "RTL9302A"
	case RTL9302A_CHIP_ID_12X2_5G:
		name = "RTL9302A_12X2_5G"
	case RTL9302B_CHIP_ID:
		name = "RTL9302B"
	case RTL9302B_CHIP_ID_8X2_5G:
		name = "RTL9302B_8X2_5G"
	case RTL9302C_CHIP_ID:
		name = "RTL9302C"
	case RTL9302C_CHIP_ID_16X2_5G:
		name = "RTL9302C_16X2_5G"
	case RTL9302D_CHIP_ID:
		name = "RTL9302D"
	case RTL9302D_CHIP_ID_24X2_5G:
		name = "RTL9302D_24X2_5G"
	case RTL9302DE_CHIP_ID:
		name = "RTL9302DE"
	case RTL9302F_CHIP_ID:
		name = "RTL9302F"
	case RTL9303_CHIP_ID:
		name = "RTL9303"
	case RTL9303_CHIP_ID_8XG:
		name = "RTL9303_8XG"
	case RTL9310_CHIP_ID:
		name = "RTL9310"
	case RTL9311_CHIP_ID:
		name = "RTL9311"
	case RTL9311E_CHIP_ID:
		name = "RTL9311E"
	case RTL9311R_CHIP_ID:
		name = "RTL9311R"
	case RTL9312_CHIP_ID:
		name = "RTL9312"
	case RTL9313_CHIP_ID:
		name = "RTL9313"
	default:
		name = "UNKNOWN"
	}
	return fmt.Sprintf("%s (0x%x)", name, uint32(cid))
}

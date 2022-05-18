package main

import (
	. "github.com/wii-tools/powerpc"
)

var RemoveRegionRestrictionsPatch = PatchSet{
	Name: "Remove region restrictions",
	Patches: []Patch{
		{
			Name:     "Allow all regions in NWC24ValidateRegion",
			AtOffset: 3840196,

			// Invalid regions have a return code of -0x34/-52.
			Before: Instructions{
				LI(R3, 0xFFCC),
			}.Bytes(),
			// Allow all regions.
			After: Instructions{
				LI(R3, 0),
			}.Bytes(),
		},
	},
}

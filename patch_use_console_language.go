package main

import (
	. "github.com/wii-tools/powerpc"
)

var UseConsoleLanguagePatch = PatchSet{
	Name: "Use console languages",
	Patches: []Patch{
		{
			Name:     "Change User-Agent formatting",
			AtOffset: 7639480,

			// WM/version/build date/friend code
			Before: []byte("WM/%s/%s/%016llu\x00\x00\x00"),
			// WM/version/build date/language/friend code
			After: []byte("WM/%s/%s/%s/%016llu"),
		},
		{
			Name:     "Use new User-Agent format",
			AtOffset: 32260,

			Before: Instructions{
				// Load a pointer to something NWC24-related.
				LWZ(R3, 0x1118, R13),
				// Get the current friend code/user ID.
				BL(0x80013468, 0x80038a74),
				// Store both halves of the ID, as it's an uint64_t.
				ORI(R29, R4, R4),
				ORI(R30, R3, R3),
				// Retrieve the current version.
				BL(0x80013474, 0x8000dea8),
				ORI(R31, R3, R3),
				// Similarly, the build revision.
				BL(0x8001347c, 0x8000dea0),
				ORI(R6, R3, R3),

				// Load the User-Agent formatting string.
				ADDI(R3, R28, 0x19),
				LI(R4, 0x64),
				// The User-Agent string is at 0x8074d0d6.
				LIS(R5, 0x8075),
				SUBI(R5, R5, 0x2f2a),

				// Format.
				ORI(R7, R31, R31),
				ORI(R10, R29, R29),
				ORI(R9, R30, R30),
				CRXOR(),
				// sprintf
				BL(0x800134ac, 0x8018fd7c),

				// Set that we've generated a User-Agent.
				LI(R0, 0x1),
				STB(R0, 0x18, R28),
				// Deconstructor
				B(0x800134b8, 0x80013454),
			}.Bytes(),
			After: Instructions{
				// Load a pointer to something NWC24-related.
				LWZ(R3, 0x1118, R13),
				// Get the current friend code/user ID.
				BL(0x80013468, 0x80038a74),
				// Store both halves of the ID, as it's an uint64_t.
				ORI(R29, R4, R4),
				ORI(R30, R3, R3),
				// Retrieve the current version.
				BL(0x80013474, 0x8000dea8),
				ORI(R31, R3, R3),
				// Similarly, the build revision.
				BL(0x8001347c, 0x8000dea0),
				ORI(R6, R3, R3),

				// Load the User-Agent formatting string.
				ADDI(R3, R28, 0x19),
				LI(R4, 0x64),
				// The User-Agent string is at 0x8074d0d6.
				LIS(R5, 0x8075),
				SUBI(R5, R5, 0x2f2a),

				// Format.
				ORI(R7, R31, R31),
				ORI(R10, R29, R29),
				ORI(R9, R30, R30),
				CRXOR(),
				// sprintf
				BL(0x800134ac, 0x8018fd7c),

				// Set that we've generated a User-Agent.
				LI(R0, 0x1),
				STB(R0, 0x18, R28),
				// Deconstructor
				B(0x800134b8, 0x80013454),
			}.Bytes(),
		},
	},
}

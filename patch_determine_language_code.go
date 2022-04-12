package main

import (
	. "github.com/wii-tools/powerpc"
)

// nullString is a helper to assist with an array of null-terminated strings.
type nullString []string

// Bytes joins all strings, terminating them with a null byte.
// They are then concatenated into a byte array.
func (n nullString) Bytes() []byte {
	var result []byte
	for _, s := range n {
		result = append(result, []byte(s)...)
		result = append(result, 0x00)
	}

	return result
}

// DetermineLanguageCodePatch adds the DetermineLanguageCode and DetermineLanguageCodeR6 functions.
// This is used for the strap screen as well as the user-agent.
var DetermineLanguageCodePatch = PatchSet{
	Name: "Add language functions",
	Patches: []Patch{
		{
			Name:     "Insert language table",
			AtOffset: 7244512,

			// We use null bytes at the end of a character conversion table.
			Before: EmptyBytes(21),
			// These names are taken from available files in the
			// strapImage_<LANG>_LZ.bin archive.
			After: nullString{
				"ja",
				"En",
				"Ge",
				"Fr",
				"Sp",
				"It",
				"Du",
			}.Bytes(),
		},
		{
			// DetermineLanguageCode is used to load the appropriate strap image.
			//
			Name:     "Add DetermineLanguageCode function",
			AtOffset: 1744,

			// We inserted this function in a field of null bytes.
			Before: EmptyBytes(160),

			After: Instructions{
				STWU(R1, R1, 0xFFF0),

				// Store the previous LR
				MFSPR(R0, LR),
				STW(R0, 0x14, R1),

				// Get current language via SCGetLanguage
				BL(0x800045dc, 0x80307f50),

				// Japanese
				CMPWI(R3, 0x0),
				BNE(0x800045e4, 0x800045f4),
				// Offset to "ja" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA00),
				B(0x800045f0, 0x80004660),

				// German
				CMPWI(R3, 0x2),
				BNE(0x800045f8, 0x80004608),
				// Offset to "Ge" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA06),
				B(0x80004604, 0x80004660),

				// French
				CMPWI(R3, 0x3),
				BNE(0x8000460c, 0x8000461c),
				// Offset to "Fr" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA09),
				B(0x80004618, 0x80004660),

				// Spanish
				CMPWI(R3, 0x4),
				BNE(0x80004620, 0x80004630),
				// Offset to "Sp" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA0C),
				B(0x8000462c, 0x80004660),

				// Italian
				CMPWI(R3, 0x5),
				BNE(0x80004634, 0x80004644),
				// Offset to "It" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA0F),
				B(0x80004640, 0x80004660),

				// Dutch
				CMPWI(R3, 0x6),
				BNE(0x80004648, 0x80004658),
				// Offset to "Du" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA12),
				B(0x80004654, 0x80004660),

				// Default to English if not another other language.
				// Offset to "En" in strings table
				LIS(R3, 0x806E),
				ORI(R3, R3, 0xCA03),

				// Finalize
				LWZ(R0, 0x14, R1),

				// Restore previous r0
				MTSPR(LR, R0),
				ADDI(R1, R1, 0x10),
				BLR(),
			}.Bytes(),
		},
		{
			Name:     "Add DetermineLanguageCodeR6 function",
			AtOffset: 2004,

			Before: EmptyBytes(36),

			// DetermineLanguageCodeR6 invokes DetermineLanguageCode,
			// and places its result in R3.
			// This is useful for sprintf reasons when sending the User-Agent.
			After: Instructions{
				STWU(R1, R1, 0xFFF0),
				MFSPR(R0, LR),
				STW(R0, 0x14, R1),

				// Use our DetermineLanguageCode function to get the language
				BL(0x800046e0, 0x800045d0),
				OR(R3, R6, R3, false),
				LWZ(R0, 0x16, R1),
				MTSPR(LR, R0),

				// Finalize
				ADDI(R1, R1, 0x10),
				BLR(),
			}.Bytes(),
		},
	},
}

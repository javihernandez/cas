/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sniff

import (
	"debug/pe"
	"os"
)

const Platform_PE = "Windows"

var machineTypes = map[uint16]string{
	pe.IMAGE_FILE_MACHINE_UNKNOWN:   "UNKNOWN",
	pe.IMAGE_FILE_MACHINE_AM33:      "AM33",
	pe.IMAGE_FILE_MACHINE_AMD64:     "AMD64",
	pe.IMAGE_FILE_MACHINE_ARM:       "ARM",
	pe.IMAGE_FILE_MACHINE_ARMNT:     "ARMNT",
	pe.IMAGE_FILE_MACHINE_ARM64:     "ARM64",
	pe.IMAGE_FILE_MACHINE_EBC:       "EBC",
	pe.IMAGE_FILE_MACHINE_I386:      "I386",
	pe.IMAGE_FILE_MACHINE_IA64:      "IA64",
	pe.IMAGE_FILE_MACHINE_M32R:      "M32R",
	pe.IMAGE_FILE_MACHINE_MIPS16:    "MIPS16",
	pe.IMAGE_FILE_MACHINE_MIPSFPU:   "MIPSFPU",
	pe.IMAGE_FILE_MACHINE_MIPSFPU16: "MIPSFPU16",
	pe.IMAGE_FILE_MACHINE_POWERPC:   "POWERPC",
	pe.IMAGE_FILE_MACHINE_POWERPCFP: "POWERPCFP",
	pe.IMAGE_FILE_MACHINE_R4000:     "R4000",
	pe.IMAGE_FILE_MACHINE_SH3:       "SH3",
	pe.IMAGE_FILE_MACHINE_SH3DSP:    "SH3DSP",
	pe.IMAGE_FILE_MACHINE_SH4:       "SH4",
	pe.IMAGE_FILE_MACHINE_SH5:       "SH5",
	pe.IMAGE_FILE_MACHINE_THUMB:     "THUMB",
	pe.IMAGE_FILE_MACHINE_WCEMIPSV2: "WCEMIPSV2",
}

func PE(file *os.File) (*Data, error) {
	f, err := pe.NewFile(file)
	if err != nil {
		return nil, err
	}

	arch := machineTypes[f.FileHeader.Machine]

	x64 := false
	switch f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		x64 = true
	}

	format := "PE32"
	if x64 {
		format += "+"
	}

	d := &Data{
		Format:   format,
		Platform: Platform_PE,
		Arch:     arch,
		X64:      x64,
		// Timestamp: f.TimeDateStamp,
	}
	return d, nil
}

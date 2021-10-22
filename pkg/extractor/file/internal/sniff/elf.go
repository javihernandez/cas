/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package sniff

import (
	"debug/elf"
	"os"
	"strings"
)

var elfosabiDesc = map[elf.OSABI]string{
	elf.ELFOSABI_HPUX:       "HP-UX operating system",
	elf.ELFOSABI_NETBSD:     "NetBSD",
	elf.ELFOSABI_LINUX:      "GNU/Linux",
	elf.ELFOSABI_HURD:       "GNU/Hurd",
	elf.ELFOSABI_86OPEN:     "86Open common IA32 ABI",
	elf.ELFOSABI_SOLARIS:    "Solaris",
	elf.ELFOSABI_AIX:        "AIX",
	elf.ELFOSABI_IRIX:       "IRIX",
	elf.ELFOSABI_FREEBSD:    "FreeBSD",
	elf.ELFOSABI_TRU64:      "TRU64 UNIX",
	elf.ELFOSABI_MODESTO:    "Novell Modesto",
	elf.ELFOSABI_OPENBSD:    "OpenBSD",
	elf.ELFOSABI_OPENVMS:    "Open VMS",
	elf.ELFOSABI_NSK:        "HP Non-Stop Kernel",
	elf.ELFOSABI_AROS:       "Amiga Research OS",
	elf.ELFOSABI_FENIXOS:    "The FenixOS highly scalable multi-core OS",
	elf.ELFOSABI_CLOUDABI:   "Nuxi CloudABI",
	elf.ELFOSABI_ARM:        "ARM",
	elf.ELFOSABI_STANDALONE: "Standalone (embedded) application",
}

func ELF(file *os.File) (*Data, error) {
	f, err := elf.NewFile(file)
	if err != nil {
		return nil, err
	}

	platform := elfosabiDesc[f.OSABI]
	if platform == "" {
		// https://refspecs.linuxfoundation.org/LSB_1.2.0/gLSB/noteabitag.html
		if abiTag := f.Section(".note.ABI-tag"); abiTag != nil {
			if data, err := abiTag.Data(); err == nil && strings.Contains(string(data), "GNU") {
				platform = "GNU/Linux"
			}
		}
	}

	d := &Data{
		Format:   "ELF",
		Type:     strings.TrimPrefix(f.Type.String(), "ET_"),
		Platform: platform,
		Arch:     strings.TrimPrefix(f.Machine.String(), "EM_"),
		X64:      f.Class == elf.ELFCLASS64,
	}
	return d, nil
}

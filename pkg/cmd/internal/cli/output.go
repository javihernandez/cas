/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/codenotary/cas/pkg/api"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"

	"github.com/codenotary/cas/pkg/cmd/internal/types"
)

func PrintError(output string, err *types.Error) error {
	if err == nil {
		return nil
	}
	switch output {
	case "":
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
	case "json":
		b, err := json.MarshalIndent(err, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func PrintWarning(output string, message string) error {
	switch output {
	case "":
		fallthrough
	case "json":
		fmt.Fprintf(os.Stderr, "\nWarning: %s\n", message)
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func PrintLc(output string, r *types.LcResult) error {
	switch output {
	case "":
		WriteLcResultTo(r, colorable.NewColorableStdout())
	case "json":
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func WriteLcResultTo(r *types.LcResult, out io.Writer) (n int64, err error) {
	if r == nil {
		return 0, nil
	}

	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 0, '\t', 0)

	printf := func(format string, a ...interface{}) error {
		m, err := fmt.Fprintf(w, format, a...)
		n += int64(m)
		return err
	}

	s := reflect.ValueOf(r).Elem()
	s = s.FieldByName("LcArtifact")
	typeOfT := s.Type()

	for i, l := 0, s.NumField(); i < l; i++ {
		f := s.Field(i)
		if key, ok := typeOfT.Field(i).Tag.Lookup("cas"); ok {
			var value string
			switch key {
			case "Size":
				if size, ok := f.Interface().(uint64); ok && size > 0 {
					value = humanize.Bytes(size)
				}
			case "Metadata":
				if metadata, ok := f.Interface().(api.Metadata); ok {
					for k, v := range metadata {
						if v == "" {
							continue
						}
						if vv, err := json.MarshalIndent(v, "\t", "    "); err == nil {
							value += fmt.Sprintf("\n\t\t%s=%s", k, string(vv))
						}
					}
					value = strings.TrimPrefix(value, "\n")
				}
			case "Apikey revoked":
				if f.IsZero() {
					value = color.New(meta.StyleWarning()).Sprintf("not available")
				} else {
					if revoked, ok := f.Interface().(*time.Time); ok {
						if revoked.IsZero() {
							value = color.New(meta.StyleAffordance()).Sprintf("no")
						} else {
							value = color.New(meta.StyleError()).Sprintf(revoked.Format(time.UnixDate))
						}
					}
				}
			case "Status":
				err = printf("Status:\t%s\n", meta.StatusNameStyled(r.Status))
				if err != nil {
					return
				}
			case "Included in":
				if included, ok := f.Interface().([]api.PackageDetails); ok {
					value += formatPackageDetails(included)
				}
			case "Dependencies":
				if deps, ok := f.Interface().([]api.PackageDetails); ok {
					value += formatPackageDetails(deps)
				}
			default:
				value = fmt.Sprintf("%s", f.Interface())
			}
			if value != "" {
				err = printf("%s:\t%s\n", key, value)
				if err != nil {
					return
				}
			}
		}
	}

	// here extra data when --verbose flag is provided
	if r.Verbose != nil {
		err = printf("\nAdditional details:\n")
		if err != nil {
			return
		}
		s = reflect.ValueOf(r.Verbose).Elem()
		typeOfT = s.Type()
		for i, l := 0, s.NumField(); i < l; i++ {
			if key, ok := typeOfT.Field(i).Tag.Lookup("cas"); ok {
				switch key {
				case "LedgerName":
					err = printf("Ledger Name:\t%s\n", r.Verbose.LedgerName)
					if err != nil {
						return
					}
				case "LocalSID":
					err = printf("Local SignerID:\t%s\n", r.Verbose.LocalSID)
					if err != nil {
						return
					}
				case "ApiKey":
					err = printf("Api-key:\t%s\n", r.Verbose.ApiKey)
					if err != nil {
						return
					}
				}
			}
		}
	}

	for _, e := range r.Errors {
		err = printf("Error:\t%s\n", color.New(meta.StyleError()).Sprintf(e.Error()))
		if err != nil {
			return
		}
	}

	return n, w.Flush()
}

func PrintLcSlice(output string, rs []*types.LcResult) error {
	switch output {
	case "":
		for _, r := range rs {
			WriteLcResultTo(r, colorable.NewColorableStdout())
			fmt.Println()
		}
	case "json":
		b, err := json.MarshalIndent(rs, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func formatPackageDetails(packages []api.PackageDetails) string {
	var ret string
	maxWidth := 0
	for _, pkg := range packages {
		width := len(pkg.Name) + len(pkg.Version)
		if width > maxWidth {
			maxWidth = width
		}
	}
	maxWidth++
	for i, pkg := range packages {
		if i != 0 {
			ret += "\n"
		}
		var line string
		if pkg.Version != "" {
			line = pkg.Name + "@" + pkg.Version
		} else {
			line = pkg.Name
		}
		ret += fmt.Sprintf("\t%-*s %s", maxWidth, line, pkg.Hash)
	}

	return ret
}

func outputNotSupportedErr(output string) error {
	return fmt.Errorf("output format not supported: %s", output)
}

func Print(output string, r *types.Result) error {
	switch output {
	case "":
		WriteResultTo(r, colorable.NewColorableStdout())
	case "json":
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		return outputNotSupportedErr(output)
	}
	return nil
}

func WriteResultTo(r *types.Result, out io.Writer) (n int64, err error) {
	if r == nil {
		return 0, nil
	}

	w := new(tabwriter.Writer)
	w.Init(out, 0, 8, 0, '\t', 0)

	printf := func(format string, a ...interface{}) error {
		m, err := fmt.Fprintf(w, format, a...)
		n += int64(m)
		return err
	}

	s := reflect.ValueOf(r).Elem()
	s = s.FieldByName("ArtifactResponse")
	typeOfT := s.Type()

	for i, l := 0, s.NumField(); i < l; i++ {
		f := s.Field(i)
		if key, ok := typeOfT.Field(i).Tag.Lookup("cas"); ok {
			var value string
			switch key {
			case "Size":
				if size, ok := f.Interface().(uint64); ok && size > 0 {
					value = humanize.Bytes(size)
				}
			case "Metadata":
				if metadata, ok := f.Interface().(api.Metadata); ok {
					for k, v := range metadata {
						if v == "" {
							continue
						}
						if vv, err := json.MarshalIndent(v, "\t", "    "); err == nil {
							value += fmt.Sprintf("\n\t\t%s=%s", k, string(vv))
						}
					}
					value = strings.TrimPrefix(value, "\n")
				}
			default:
				value = fmt.Sprintf("%s", f.Interface())
			}
			if value != "" {
				err = printf("%s:\t%s\n", key, value)
				if err != nil {
					return
				}
			}
		}
	}

	for _, e := range r.Errors {
		err = printf("Error:\t%s\n", color.New(meta.StyleError()).Sprintf(e.Error()))
		if err != nil {
			return
		}
	}

	return n, w.Flush()
}

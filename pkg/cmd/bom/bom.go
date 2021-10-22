/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package bom

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/codenotary/cas/pkg/bom"
	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/docker"
	"github.com/codenotary/cas/pkg/uri"
)

// NewCommand returns the cobra command for `cas info`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bom",
		Example: "  cas bom docker://alpine",
		Short:   "Collect BoM information",
		Long: `
Collect BoM (Bill of Material) information

It identifies dependencies of build artifact and produces the BoM. Dependencies can be
later authenticated by 'cas a --bom', and notarized together with artifact by 'cas n --bom'.
`,
		RunE: runBom,
		PreRun: func(cmd *cobra.Command, args []string) {
			// Bind to all flags to env vars (after flags were parsed),
			// but only ones retrivied by using viper will be used.
			viper.BindPFlags(cmd.Flags())
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactValidArgs(1)(cmd, args)
		},
	}

	// BoM output options
	cmd.Flags().String("bom-spdx", "", "name of the file to output BoM in SPDX format")
	cmd.Flags().String("bom-cyclonedx-json", "", "name of the file to output BoM in CycloneDX JSON format")
	cmd.Flags().String("bom-cyclonedx-xml", "", "name of the file to output BoM in CycloneDX XML format")

	return cmd
}

func runBom(cmd *cobra.Command, args []string) error {
	path := args[0]
	u, err := uri.Parse(path)
	if err != nil {
		return err
	}
	if _, ok := bom.BomSchemes[u.Scheme]; !ok {
		return fmt.Errorf("unsupported URI %s for --bom option", path)
	}
	if u.Scheme != "" {
		path = strings.TrimPrefix(u.Opaque, "//")
	}

	var bomArtifact artifact.Artifact
	if u.Scheme == "docker" {
		bomArtifact, err = docker.New(path)
		if err != nil {
			return err
		}
	} else {
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		bomArtifact = bom.New(path)
	}
	if bomArtifact == nil {
		return fmt.Errorf("unsupported artifact format/language")
	}

	outputOpts := artifact.Progress
	if viper.GetBool("silent") {
		outputOpts = artifact.Silent
	}

	fmt.Printf("Resolving dependencies...\n")
	_, err = bomArtifact.ResolveDependencies(outputOpts)
	if err != nil {
		return fmt.Errorf("cannot get dependencies: %w", err)
	}

	artifact.Display(bomArtifact, artifact.ColNameVersion)

	return bom.Output(bomArtifact) // process all possible BoM output options
}

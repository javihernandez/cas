package api

import (
	"fmt"

	"github.com/vchain-us/ledger-compliance-go/schema"

	"github.com/codenotary/cas/pkg/meta"
)

// Sign is invoked by the User to notarize an artifact using the given functional options,
// By default, the artifact is notarized using status = meta.StatusTrusted, visibility meta.VisibilityPrivate.
func (u LcUser) Sign(artifact Artifact, options ...LcSignOption) (uint64, error) {
	if artifact.Hash == "" {
		return 0, makeError("hash is missing", nil)
	}
	if artifact.Size < 0 {
		return 0, makeError("invalid size", nil)
	}

	o, err := makeLcSignOpts(options...)
	if err != nil {
		return 0, err
	}

	return u.createArtifacts([]*Artifact{&artifact}, []meta.Status{o.status}, [][]*schema.VCNDependency{o.bom})
}

// SignMulti ...
func (u LcUser) SignMulti(artifacts []*Artifact, options [][]LcSignOption) (uint64, error) {
	if len(artifacts) != len(options) {
		return 0, makeError("the number of options must be the same as the number artifacts", nil)
	}

	statuses := make([]meta.Status, len(artifacts))
	boms := make([][]*schema.VCNDependency, len(artifacts))
	for i, artifact := range artifacts {
		if artifact.Hash == "" {
			return 0, makeError(fmt.Sprintf("hash is missing for artifact %s", artifact.Name), nil)
		}
		if artifact.Size < 0 {
			return 0, makeError(fmt.Sprintf("invalid size for artifact %s", artifact.Name), nil)
		}

		o, err := makeLcSignOpts(options[i]...)
		if err != nil {
			return 0, err
		}
		statuses[i] = o.status
		boms[i] = o.bom
	}

	return u.createArtifacts(artifacts, statuses, boms)
}

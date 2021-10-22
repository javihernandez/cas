package types

import "github.com/codenotary/cas/pkg/api"

type Result struct {
	api.ArtifactResponse `yaml:",inline"`
	Errors               []error `json:"error,omitempty" yaml:"error,omitempty"`
}

func (r *Result) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

func NewResult(a *api.Artifact, ar *api.ArtifactResponse) *Result {

	var r Result

	switch true {
	case ar != nil:
		r = Result{*ar, nil}
	case a != nil:
		r = Result{api.ArtifactResponse{
			Name:        a.Name,
			Kind:        a.Kind,
			Hash:        a.Hash,
			Size:        a.Size,
			ContentType: a.ContentType,
			Metadata:    a.Metadata,
		}, nil}
	default:
		r = Result{}
	}

	r.Status = ""

	return &r
}

/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package api

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	immuschema "github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/vchain-us/ledger-compliance-go/schema"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (a Artifact) ToLcArtifact() *LcArtifact {
	aR := &LcArtifact{
		// root fields
		Kind:        a.Kind,
		Name:        a.Name,
		Hash:        a.Hash,
		Size:        a.Size,
		ContentType: a.ContentType,

		// custom metadata
		Metadata: a.Metadata,
	}

	return aR
}
func ItemToLcArtifact(item *schema.ItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(item.Item.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(item.Timestamp.GetSeconds(), int64(item.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if item.ApikeyRevoked != nil {
		if item.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(item.ApikeyRevoked.GetSeconds(), int64(item.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = item.LedgerName
	return &lca, nil
}

func ZItemToLcArtifact(ie *schema.ZItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(ie.Item.Entry.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(ie.Timestamp.GetSeconds(), int64(ie.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if ie.ApikeyRevoked != nil {
		if ie.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(ie.ApikeyRevoked.GetSeconds(), int64(ie.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = ie.LedgerName
	return &lca, nil
}

func VerifiableItemExtToLcArtifact(item *schema.VerifiableItemExt) (*LcArtifact, error) {
	var lca LcArtifact
	err := json.Unmarshal(item.Item.Entry.Value, &lca)
	if err != nil {
		return nil, err
	}
	ts := time.Unix(item.Timestamp.GetSeconds(), int64(item.Timestamp.GetNanos()))
	lca.Uid = strconv.Itoa(int(ts.UnixNano()))
	lca.Timestamp = ts.UTC()
	// if ApikeyRevoked == nil no revoked infos available. Old key type
	if item.ApikeyRevoked != nil {
		if item.ApikeyRevoked.GetSeconds() > 0 {
			t := time.Unix(item.ApikeyRevoked.GetSeconds(), int64(item.ApikeyRevoked.Nanos)).UTC()
			lca.Revoked = &t
		} else {
			lca.Revoked = &time.Time{}
		}
	}
	lca.Ledger = item.LedgerName
	lca.PublicKey = item.PublicKey
	return &lca, nil
}

type LcArtifact struct {
	// root fields
	Uid         string    `json:"uid" yaml:"uid" cas:"UID"`
	Kind        string    `json:"kind" yaml:"kind" cas:"Kind"`
	Name        string    `json:"name" yaml:"name" cas:"Name"`
	Hash        string    `json:"hash" yaml:"hash" cas:"Hash"`
	Size        uint64    `json:"size" yaml:"size" cas:"Size"`
	Timestamp   time.Time `json:"timestamp,omitempty" yaml:"timestamp" cas:"Timestamp"`
	ContentType string    `json:"contentType" yaml:"contentType" cas:"ContentType"`

	// custom metadata
	Metadata Metadata `json:"metadata" yaml:"metadata" cas:"Metadata"`

	Signer  string      `json:"signer" yaml:"signer" cas:"SignerID"`
	Revoked *time.Time  `json:"revoked,omitempty" yaml:"revoked" cas:"Apikey revoked"`
	Status  meta.Status `json:"status" yaml:"status" cas:"Status"`
	Ledger  string      `json:"ledger,omitempty" yaml:"ledger"`

	IncludedIn []PackageDetails `json:"included_in,omitempty" yaml:"included_in,omitempty" cas:"Included in"`
	Deps       []PackageDetails `json:"bom,omitempty" yaml:"bom,omitempty" cas:"Dependencies"`
	PublicKey  string
}

func (u LcUser) artifactToCasArtifact(
	artifact *Artifact,
	status meta.Status,
	deps []*schema.VCNDependency,
) (*schema.VCNArtifact, error) {

	casArtifact := schema.VCNArtifact{Dependencies: deps}

	aR := artifact.ToLcArtifact()
	aR.Status = status

	aR.Signer = GetSignerIDByApiKey(u.Client.ApiKey)

	arJSON, err := json.Marshal(aR)
	if err != nil {
		return nil, err
	}

	casArtifact.Artifact = arJSON

	return &casArtifact, nil
}

func (u LcUser) createArtifacts(
	artifacts []*Artifact,
	statuses []meta.Status,
	boms [][]*schema.VCNDependency,
) (uint64, error) {

	if len(artifacts) != len(statuses) || len(artifacts) != len(boms) {
		return 0, errors.New(
			"artifacts, statuses and bomTexts must have the same length")
	}

	req := schema.VCNArtifactsRequest{
		Artifacts: make([]*schema.VCNArtifact, 0, len(artifacts)),
	}

	now := time.Now().UTC()
	nowBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nowBytes, uint64(now.UnixNano()))

	for i := 0; i < len(artifacts); i++ {
		casArtifact, err := u.artifactToCasArtifact(artifacts[i], statuses[i], boms[i])
		if err != nil {
			return 0, err
		}

		req.Artifacts = append(req.Artifacts, casArtifact)
	}

	md := metadata.Pairs(
		meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue,
		meta.CasCmdHeaderName, meta.CasNotarizeCmdHeaderValue,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := u.Client.VCNSetArtifacts(ctx, &req)
	if err != nil {
		return 0, err
	}

	return resp.GetTransaction().GetId(), err
}

// LoadArtifact fetches and returns an *lcArtifact for the given hash and current u, if any.
func (u *LcUser) LoadArtifact(
	hash, signerID string,
	uid string,
	tx uint64,
	gRPCMetadata map[string][]string,
) (lc *LcArtifact, verified bool, err error) {

	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	if len(gRPCMetadata) > 0 {
		md = metadata.Join(md, gRPCMetadata)
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if signerID == "" {
		signerID = GetSignerIDByApiKey(u.Client.ApiKey)
	}

	key := AppendPrefix(meta.CasPluginTypeHeaderValue, []byte(signerID))
	key = AppendSignerId(hash, key)

	var jsonAr *schema.VerifiableItemExt
	if uid != "" {
		score, err := strconv.ParseFloat(uid, 64)
		if err != nil {
			return nil, false, err
		}
		zitems, err := u.Client.ZScanExt(ctx, &immuschema.ZScanRequest{
			Set:       key,
			SeekScore: math.MaxFloat64,
			SeekAtTx:  tx,
			Limit:     1,
			MinScore:  &immuschema.Score{Score: score},
			MaxScore:  &immuschema.Score{Score: score},
			SinceTx:   math.MaxUint64,
			NoWait:    true,
		})
		if err != nil {
			return nil, false, err
		}
		if len(zitems.Items) > 0 {
			jsonAr, err = u.Client.VerifiedGetExtAt(ctx, zitems.Items[0].Item.Key, zitems.Items[0].Item.AtTx)
		} else {
			return nil, false, ErrNotFound
		}
	} else {
		jsonAr, err = u.Client.VerifiedGetExtAt(ctx, key, tx)
	}
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Message() == "data is corrupted" {
			return nil, false, ErrNotVerified
		}
		if err.Error() == "data is corrupted" {
			return nil, false, ErrNotVerified
		}
		if ok && s.Message() == "key not found" {
			return nil, false, ErrNotFound
		}
		return nil, true, err
	}

	lcArtifact, err := VerifiableItemExtToLcArtifact(jsonAr)
	if err != nil {
		return nil, false, err
	}

	return lcArtifact, true, nil
}

// LoadArtifacts fetches and returns multiple *lcArtifact for the given hashes and current u, if any.
func (u *LcUser) LoadArtifacts(
	signerID string,
	hashes []string,
	gRPCMetadata map[string][]string,
) (artifacts []*LcArtifact, verified []bool, errs []error, err error) {

	md := metadata.Pairs(meta.CasPluginTypeHeaderName, meta.CasPluginTypeHeaderValue)
	if len(gRPCMetadata) > 0 {
		md = metadata.Join(md, gRPCMetadata)
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if signerID == "" {
		signerID = GetSignerIDByApiKey(u.Client.ApiKey)
	}

	prefixedSignerID := AppendPrefix(meta.CasPluginTypeHeaderValue, []byte(signerID))

	keys := make([][]byte, 0, len(hashes))
	for _, hash := range hashes {
		key := AppendSignerId(hash, prefixedSignerID)
		keys = append(keys, key)
	}

	itemsExt, errsMsgs, err := u.Client.VerifiedGetExtAtMulti(ctx, keys, make([]uint64, len(keys)))
	if err != nil {
		return nil, nil, nil, err
	}

	if len(itemsExt) != len(keys) || len(errsMsgs) != len(keys) {
		return nil, nil, nil, fmt.Errorf(
			"internal logic error: expected size of the reponse %d, got %d items and %d errors",
			len(keys), len(itemsExt), len(errsMsgs))
	}

	lcArtifacts := make([]*LcArtifact, len(itemsExt))
	verified = make([]bool, len(itemsExt))
	errs = make([]error, len(itemsExt))

	for i := 0; i < len(keys); i++ {
		if len(errsMsgs[i]) > 0 {
			switch {
			case strings.HasSuffix(errsMsgs[i], "data is corrupted"):
				errs[i] = ErrNotVerified
			case strings.HasSuffix(errsMsgs[i], "key not found"):
				errs[i] = ErrNotFound
			default:
				errs[i] = errors.New(errsMsgs[i])
				verified[i] = true
			}
			continue
		}

		lcArtifact, err := VerifiableItemExtToLcArtifact(itemsExt[i])
		if err != nil {
			return nil, nil, nil, err
		}
		lcArtifacts[i] = lcArtifact
		verified[i] = true
	}

	return lcArtifacts, verified, errs, nil
}

func AppendPrefix(prefix string, key []byte) []byte {
	var prefixed = make([]byte, len(prefix)+1+len(key))
	copy(prefixed[0:], prefix+".")
	copy(prefixed[len(prefix)+1:], key)
	return prefixed
}

func AppendSignerId(signerId string, k []byte) []byte {
	var prefixed = make([]byte, len(k)+len(signerId)+1)
	copy(prefixed[0:], k)
	copy(prefixed[len(k):], "."+signerId)
	return prefixed
}

// Date returns a RFC3339 formatted string of verification time (v.Timestamp), if any, otherwise an empty string.
func (lca *LcArtifact) Date() string {
	if lca != nil {
		ut := lca.Timestamp.UTC()
		if ut.Unix() > 0 {
			return ut.Format(time.RFC3339)
		}
	}
	return ""
}

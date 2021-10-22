package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/codenotary/cas/pkg/meta"
	"github.com/codenotary/cas/pkg/signature"
	"github.com/fatih/color"
	"google.golang.org/grpc/status"
)

// CheckConnectionPublicKey the aim of this method is to guarantee that the connection between cas and a CAS server are verified by the first login auto trusted signature.
// This method fetches an immudb state, checks if the public key provided to the immudb client match server signature and
// saves locally such key.
// In addition it checks if a previously trusted (local) key is the same to the current one used by client. This guarantee that the connection is established on a previously trusted server.
// If enforceSignatureVerify is TRUE it requires an explicit fingerprint confirmation.
// NOTE: if CAS_SIGNING_PUB_KEY_FILE or CAS_SIGNING_PUB_KEY environment flag or arguments are provided this method is not called.
func (u *LcUser) CheckConnectionPublicKey(enforceSignatureVerify bool) error {
	state, err := u.Client.CurrentState(context.Background())
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Message() == "unable to verify signature: no signature found" {
				// for security reason if is present a trusted public key we return an error also if enforceSignatureVerify = true. Client was using on a secure server so it's not secure anymore.
				return fmt.Errorf("Community Attestation Service is not signing messages but a public key %s was found in HOME folder. In order to continue with a not signed connection please remove such key", meta.CasSigningPubKeyFileName)
			}
			if st.Message() == "signature doesn't match provided public key" {
				color.Set(meta.StyleWarning())
				fmt.Printf("previously trusted Community Attestation Service changed its signature. In order to trust again the server please provide a new public key or remove %s stored in home folder.", meta.CasSigningPubKeyFileName)
				fmt.Println()
				color.Unset()
				return fmt.Errorf("operation aborted : %w", st.Err())
			}
		}
		return err
	}

	if state.Signature == nil && enforceSignatureVerify {
		return errors.New("Community Attestation Service is not signing messages. Operation aborted")
	}

	if state.Signature != nil && state.Signature.GetPublicKey() != nil {
		ECDSAPk := signature.UnmarshalKey(state.Signature.GetPublicKey())
		pk, err := signature.ConfirmFingerprint(ECDSAPk, enforceSignatureVerify)
		if err != nil {
			return err
		}
		if pk != nil {
			u.Client.SetServerSigningPubKey(pk)
		}
	}
	return nil
}

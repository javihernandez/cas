/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package signature

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/codenotary/immudb/pkg/signer"
	"github.com/codenotary/cas/pkg/meta"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// PrepareSignatureParams returns signing public keys from an environment var (string) or local file on file system.
// When public key is specified it disables signature comparison with the local public key previously trusted.
// If no external public key is explicit provided it try to the local public key, if present.
func PrepareSignatureParams(signingPubKeyString, signingPubKeyFile string) (*ecdsa.PublicKey, bool, error) {
	skipLocalPubKeyComp := false
	signingPubKey, err := getECDSAPublicKeyFromFlags(signingPubKeyString, signingPubKeyFile)
	if err != nil {
		return nil, false, err
	}
	if signingPubKey != nil {
		// if explicit public key is provided here is disabled public key comparison, confirmation and local saving of local public key.
		skipLocalPubKeyComp = true
	}

	// if an explicit public key is not provided cas try to fetch an already trusted public key stored locally
	if signingPubKey == nil {
		signingPubKey, err = getLocalPubKey()
		if err != nil && !os.IsNotExist(err) {
			return nil, false, err
		}
	}
	return signingPubKey, skipLocalPubKeyComp, nil
}

// ConfirmFingerprint Checks if the public key pk is equal to an already previously confirmed (local) key.
// If a local public key not exists cas automatically trust the provided key, except enforceSignatureVerify is true,
// in that case cas will require explicit acceptation of the signature that came from server.
func ConfirmFingerprint(pk *ecdsa.PublicKey, enforceSignatureVerify bool) (confirmed *ecdsa.PublicKey, err error) {
	if pk == nil {
		return nil, errors.New("fingerprint confirmation is not required")
	}

	// if public key does not exist the first provided is accepted and stored
	localPk, err := getLocalPubKey()
	if err != nil && os.IsNotExist(err) && !enforceSignatureVerify {
		err = setSigningPubKey(pk)
		if err != nil {
			return nil, err
		}
		color.Set(meta.StyleWarning())
		fmt.Println("CAS automatically trusted the signature found on current connection")
		color.Unset()
		return pk, nil
	}

	if !terminal.IsTerminal(int(os.Stdout.Fd())) && enforceSignatureVerify {
		return nil, errors.New("can't be run not interactively if CAS_SIGNING_PUB_KEY_FILE or CAS_SIGNING_PUB_KEY env is not provided and \"enforce signature verify\" is ON")
	}

	if localPk != nil && bytes.Compare(marshalKey(localPk), marshalKey(pk)) == 0 {
		return pk, nil
	}
	publicKey, err := ssh.NewPublicKey(pk)
	if err != nil {
		return nil, err
	}

	fingerprint := ssh.FingerprintSHA256(publicKey)

	fmt.Printf("This connection contains signed messages but the authenticity of the current Community Attestation Service service can't be established. \nECDSA fingerprint provided is:\n%s\nAre you sure you want to continue connecting? (Y\\n): ", fingerprint)

	var confirm string
	_, err = fmt.Scanln(&confirm)
	if err != nil {
		return nil, err
	}
	if confirm != "Y" && confirm != "n" {
		return nil, errors.New("please enter Y or n")
	}
	if confirm == "n" {
		return nil, errors.New("connection aborted")
	}
	err = setSigningPubKey(pk)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// GetECDSAPublicKeyFromBytes parses bytes content in order to have an ecdsa public key.
// A valid trailer and header is needed here. Use PreparePublicKey in case of need.
func GetECDSAPublicKeyFromBytes(publicKey []byte) (*ecdsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode((publicKey))
	if publicKeyBlock == nil {
		return nil, errors.New("no ecdsa key found")
	}
	cert, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.(*ecdsa.PublicKey), nil
}

// UnmarshalKey unmarshal an ecdsa public key contained in immudb status signature
func UnmarshalKey(publicKey []byte) *ecdsa.PublicKey {
	if publicKey == nil || len(publicKey) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKey)
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
}

func getECDSAPublicKeyFromFlags(signingPubKeyString, signingPubKeyFile string) (signingPubKey *ecdsa.PublicKey, err error) {
	if signingPubKeyString != "" && signingPubKeyFile != "" {
		return nil, fmt.Errorf("cannot use both --signing-pub-key-file and --signing-pub-key")
	}

	if signingPubKeyString != "" {
		signingPubKey, err = GetECDSAPublicKeyFromBytes(PreparePublicKey(signingPubKeyString))
		if err != nil {
			return nil, err
		}
	}

	if signingPubKey == nil {
		signingPubKey, err = getECDSAPublicKeyFromFile(signingPubKeyFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	return signingPubKey, nil
}

// getLocalPubKey it lookup a key in HOME folder.  Returns an ecdsa public key or error
func getLocalPubKey() (*ecdsa.PublicKey, error) {
	hd, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	keyFile := path.Join(hd, meta.CasSigningPubKeyFileName)

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return nil, err
	}
	return signer.ParsePublicKeyFile(keyFile)
}

func setSigningPubKey(pk *ecdsa.PublicKey) error {
	hd, err := homedir.Dir()
	if err != nil {
		return err
	}

	keyFile := path.Join(hd, meta.CasSigningPubKeyFileName)

	kf, err := os.OpenFile(keyFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	pb := marshalKey(pk)
	_, err = kf.Write(pb)
	if err := kf.Close(); err != nil {
		return err
	}
	color.Set(meta.StyleAffordance())
	fmt.Println("CAS saved locally the trusted public key")
	color.Unset()
	return nil
}

func marshalKey(publicKey *ecdsa.PublicKey) []byte {
	if publicKey == nil {
		return nil
	}
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncodedPub
}

func getECDSAPublicKeyFromFile(filePath string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return GetECDSAPublicKeyFromBytes(publicKeyBytes)
}

// PreparePublicKey append and prepend public key header and footer on a public key content
func PreparePublicKey(key string) []byte {
	var pk = "-----BEGIN PUBLIC KEY-----\n" +
		key +
		"\n-----END PUBLIC KEY-----"
	return []byte(pk)
}

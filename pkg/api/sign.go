/*
 * Copyright (c) 2018-2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package api

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/vchain-us/vcn/internal/blockchain"
	"github.com/vchain-us/vcn/internal/errors"
	"github.com/vchain-us/vcn/pkg/meta"
)

const walletNotSyncMsg = `
%s cannot be signed with CodeNotary. 
We are finalizing your account configuration. We will complete the 
configuration shortly and we will update you as soon as this is done.
We are sorry for the inconvenience and would like to thank you for 
your patience.
It only takes few seconds. Please try again in 1 minute.
`

func (u User) Sign(artifact Artifact, pubKey string, passphrase string, state meta.Status, visibility meta.Visibility) error {

	hasAuth, err := u.IsAuthenticated()
	if err != nil {
		return err
	}
	if !hasAuth {
		return makeError("authentication required, please login", nil)
	}

	if artifact.Hash == "" {
		return makeError("asset's hash is missing", nil)
	}
	if artifact.Name == "" {
		return makeError("asset's name is missing", nil)
	}

	keyin, err := u.cfg.OpenKey(pubKey)
	if err != nil {
		return err
	}

	synced, err := u.isWalletSynced(pubKey)
	if err != nil {
		return err
	}
	if !synced {
		return makeError(fmt.Sprintf(walletNotSyncMsg, artifact.Name), nil)
	}

	return u.commitHash(keyin, passphrase, artifact.Hash, artifact.Name, artifact.Size, state, visibility)
}

// todo(leogr): refactor
func (u User) commitHash(
	keyin io.Reader,
	passphrase string,
	hash string,
	name string,
	fileSize uint64,
	status meta.Status,
	visibility meta.Visibility,
) (err error) {
	transactor, err := bind.NewTransactor(keyin, passphrase)
	if err != nil {
		return
	}

	transactor.GasLimit = meta.GasLimit()
	transactor.GasPrice = meta.GasPrice()
	client, err := ethclient.Dial(meta.MainNetEndpoint())
	if err != nil {
		err = makeError("Cannot connect to blockchain", logrus.Fields{
			"error":   err,
			"network": meta.MainNetEndpoint(),
		})
		return
	}
	address := common.HexToAddress(meta.AssetsRelayContractAddress())
	instance, err := blockchain.NewAssetsRelay(address, client)
	if err != nil {
		err = makeFatal(
			"Cannot instantiate contract",
			logrus.Fields{
				"error":    err,
				"contract": meta.AssetsRelayContractAddress(),
			},
		)
		return
	}
	tx, err := instance.Sign(transactor, hash, big.NewInt(int64(status)))
	if err != nil {
		err = makeFatal(
			"method <sign> failed",
			logrus.Fields{
				"error": err,
				"hash":  hash,
			},
		)
		return
	}
	timeout, err := waitForTx(tx.Hash(), meta.TxVerificationRounds(), meta.PollInterval())
	if err != nil {
		// fixme(leogr): logging, and avoid to output directly
		errors.PrintErrorURLCustom("blockchain-permission", 403)
		err = makeFatal(
			"Could not write to blockchain",
			logrus.Fields{
				"error": err,
			},
		)
		return
	}
	if timeout {
		err = makeFatal(
			"Writing to blockchain timed out",
			logrus.Fields{
				"error": err,
			},
		)
		return
	}

	pubKey := transactor.From.Hex()
	verification, err := BlockChainVerifyMatchingPublicKey(hash, pubKey)
	if err != nil {
		return
	}

	err = u.CreateArtifact(verification, strings.ToLower(pubKey), name, hash, fileSize, visibility, status)
	if err != nil {
		return
	}

	// todo(ameingast): redundant tracking events?
	_ = TrackPublisher(&u, meta.VcnSignEvent)
	_ = TrackSign(&u, hash, name, status)
	return
}

func waitForTx(tx common.Hash, maxRounds uint64, pollInterval time.Duration) (timeout bool, err error) {
	client, err := ethclient.Dial(meta.MainNetEndpoint())
	if err != nil {
		return false, err
	}
	for i := uint64(0); i < maxRounds; i++ {
		_, pending, err := client.TransactionByHash(context.Background(), tx)
		if err != nil {
			return false, err
		}
		if !pending {
			return false, nil
		}
		time.Sleep(pollInterval)
	}
	return true, nil
}

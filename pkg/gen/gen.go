package gen

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
)

// Keys for map
const stakerKey = "staker.key_BASE64"
const stakerCert = "staker.crt_BASE64"
const signerKey = "signer.key_BASE64"
const publicKey = "publicKey" //hex encoded
const nodeID = "nodeID"

// GenerateKeys generates sensitive data
func GenerateKeys() (map[string]string, error) {

	cBytes, kBytes, err := staking.NewCertAndKeyBytes()
	if err != nil {
		return nil, err
	}

	id, err := getNodeID(cBytes, kBytes)
	if err != nil {
		return nil, err
	}

	blsSignerKey, err := bls.NewSecretKey()
	if err != nil {
		return nil, err
	}
	blsSignerBytes := bls.SecretKeyToBytes(blsSignerKey)

	blsPublicBytes := bls.PublicKeyToBytes(bls.PublicFromSecretKey(blsSignerKey))
	blsPublicHex := make([]byte, hex.EncodedLen(len(blsPublicBytes)))
	hex.Encode(blsPublicHex, blsPublicBytes)

	return map[string]string{
		stakerCert: base64.StdEncoding.EncodeToString(cBytes),
		stakerKey:  base64.StdEncoding.EncodeToString(kBytes),
		signerKey:  base64.StdEncoding.EncodeToString(blsSignerBytes),
		publicKey:  string(blsPublicBytes),
		nodeID:     id,
	}, nil
}

// getNodeID returns avalanchego nodeID based on certificate and key
func getNodeID(certBytes []byte, keyBytes []byte) (string, error) {
	cert, err := staking.LoadTLSCertFromBytes(keyBytes, certBytes)

	if err != nil {
		return "", err
	}

	return ids.NodeIDFromCert(cert.Leaf).String(), nil
}

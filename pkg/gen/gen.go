package gen

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/hashing"
)

const stakerKey = "staker.key_BASE64"
const stakerCert = "staker.crt_BASE64"
const signerKey = "signer.key_BASE64"
const publicKey = "publicKey"
const nodeID = "nodeID"

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

	return map[string]string{
		stakerCert: base64.StdEncoding.EncodeToString(kBytes),
		stakerKey:  base64.StdEncoding.EncodeToString(kBytes),
		signerKey:  base64.StdEncoding.EncodeToString(blsSignerBytes),
		publicKey:  string(blsPublicBytes),
		nodeID:     id,
	}, nil
}

func getNodeID(certBytes []byte, keyBytes []byte) (string, error) {
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return "", err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return "", err
	}
	nodeID, err := ids.ToShortID(hashing.PubkeyBytesToAddress(cert.Leaf.Raw))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", nodeID), nil
}

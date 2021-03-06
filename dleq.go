package dleq

/*
#cgo LDFLAGS: -L${SRCDIR}/lib ${SRCDIR}/lib/libdleq.so
#include "./lib/libdleq.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	Ed25519PrivateKeySize   = 32
	Ed25519PublicKeySize    = 32
	Secp256k1PrivateKeySize = 32
	Secp256k1PublicKeySize  = 33
)

// Proof represents a serialized DLEq proof
type Proof []byte

// PrivateKey represents a 32-byte private key
type PrivateKey []byte

// Ed25519PublicKey represents a 32-byte ed25519 public key
type Ed25519PublicKey []byte

// Secp256k1PublicKey represents a 33-byte secp256k1 public key
type Secp256k1PublicKey []byte

// Ed25519Secp256k1Prove creates a DLEq proof of the returned private key
// for ed25519 and secp256k1.
func Ed25519Secp256k1Prove() (Proof, PrivateKey, error) {
	// TOOD: malloc instread of using "make" cause of GC?
	proofSize := C.ed25519_secp256k1_proof_size()
	dst := make([]byte, proofSize)
	ptr := unsafe.Pointer(&dst[0])
	keyDst := make([]byte, Ed25519PrivateKeySize)
	keyPtr := unsafe.Pointer(&keyDst[0])

	ok := C.ed25519_secp256k1_prove((*C.char)(ptr), (*C.char)(keyPtr))
	if byte(ok) == 0 {
		return nil, nil, fmt.Errorf("failed to generate proof")
	}

	return Proof(dst), PrivateKey(keyDst), nil
}

// Ed25519Secp256k1Verify verifies the given DLEq proof.
// It returns an error if verification fails.
func Ed25519Secp256k1Verify(proof []byte) (Ed25519PublicKey, Secp256k1PublicKey, error) {
	dst := make([]byte, Ed25519PublicKeySize+Secp256k1PublicKeySize)
	dstPtr := unsafe.Pointer(&dst[0])
	proofPtr := unsafe.Pointer(&proof[0])

	ok := C.ed25519_secp256k1_verify((*C.char)(proofPtr), (*C.char)(dstPtr))
	if byte(ok) == 0 {
		return nil, nil, fmt.Errorf("failed to verify proof")
	}

	return Ed25519PublicKey(dst[:32]), Secp256k1PublicKey(dst[32:]), nil
}

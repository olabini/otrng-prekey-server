package prekeyserver

import (
	"crypto/rand"
	"math/big"

	"github.com/coyim/gotrax"
	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_generateEcdhProof_and_verify_generatesProofsThatValidates(c *C) {
	wr := gotrax.ReaderIntoWithRandom(gotrax.FixtureRand())
	values := make([]*gotrax.Keypair, 3)
	values[0] = gotrax.GenerateKeypair(wr)
	values[1] = gotrax.GenerateKeypair(wr)
	values[2] = gotrax.GenerateKeypair(wr)
	m := [64]byte{0x01, 0x02, 0x03}

	proof, e := generateEcdhProof(wr, values, m[:], usageProofMessageEcdh)
	c.Assert(e, IsNil)

	values2 := make([]*gotrax.PublicKey, 3)
	values2[0] = values[0].Pub
	values2[1] = values[1].Pub
	values2[2] = values[2].Pub

	c.Assert(proof.verify(values2, m[:], usageProofMessageEcdh), Equals, true)
	c.Assert(proof.verify(values2, m[:], usageProofSharedEcdh), Equals, false)

	m2 := [64]byte{0x02, 0x02, 0x03}
	c.Assert(proof.verify(values2, m2[:], usageProofMessageEcdh), Equals, false)

	wrongDL := gotrax.GenerateKeypair(wr)
	values2[1] = wrongDL.Pub
	c.Assert(proof.verify(values2, m[:], usageProofMessageEcdh), Equals, false)
}

func randomDhSecretValue(wr gotrax.WithRandom) *big.Int {
	res, _ := rand.Int(wr.RandReader(), dhQ)
	return res
}

func (s *GenericServerSuite) Test_generateDhProof_and_verify_generatesProofsThatValidates(c *C) {
	wr := gotrax.ReaderIntoWithRandom(gotrax.FixtureRand())
	valuesPriv := make([]*big.Int, 3)
	valuesPriv[0] = randomDhSecretValue(wr)
	valuesPriv[1] = randomDhSecretValue(wr)
	valuesPriv[2] = randomDhSecretValue(wr)

	valuesPub := make([]*big.Int, 3)
	valuesPub[0] = new(big.Int).Exp(g3, valuesPriv[0], dhP)
	valuesPub[1] = new(big.Int).Exp(g3, valuesPriv[1], dhP)
	valuesPub[2] = new(big.Int).Exp(g3, valuesPriv[2], dhP)

	m := [64]byte{0x01, 0x02, 0x03}

	proof, e := generateDhProof(wr, valuesPriv, valuesPub, m[:], usageProofMessageDh, nil)
	c.Assert(e, IsNil)

	c.Assert(proof.verify(valuesPub, m[:], usageProofMessageDh), Equals, true)
	c.Assert(proof.verify(valuesPub, m[:], usageProofSharedEcdh), Equals, false)

	m2 := [64]byte{0x02, 0x02, 0x03}
	c.Assert(proof.verify(valuesPub, m2[:], usageProofMessageDh), Equals, false)

	valuesPub[1].Mul(valuesPub[1], valuesPub[1])
	valuesPub[1].Mod(valuesPub[1], dhP)
	c.Assert(proof.verify(valuesPub, m[:], usageProofMessageDh), Equals, false)
}

func (s *GenericServerSuite) Test_generateEcdhProof_generatesSpecificValue(c *C) {
	wr := gotrax.ReaderIntoWithRandom(gotrax.FixtureRand())
	values := make([]*gotrax.Keypair, 5)
	values[0] = gotrax.GenerateKeypair(wr)
	values[1] = gotrax.GenerateKeypair(wr)
	values[2] = gotrax.GenerateKeypair(wr)
	values[3] = gotrax.GenerateKeypair(wr)
	values[4] = gotrax.GenerateKeypair(wr)
	m := [64]byte{0x03, 0x03, 0x01}

	proof, e := generateEcdhProof(wr, values, m[:], usageProofMessageEcdh)
	c.Assert(e, IsNil)

	c.Assert(proof.c, DeepEquals, []byte{
		0xbc, 0x8c, 0xb6, 0x80, 0xa5, 0x0c, 0x9e, 0x50,
		0xb3, 0x01, 0x8a, 0x36, 0x95, 0x20, 0xac, 0x54,
		0xfc, 0x30, 0xdf, 0x78, 0x0e, 0xc6, 0xdd, 0x1e,
		0xa7, 0x15, 0xae, 0x83, 0x09, 0x50, 0x22, 0xfe,
		0xd2, 0x9e, 0x44, 0x5a, 0x7b, 0x04, 0xbb, 0x4c,
		0x27, 0xe1, 0x55, 0x1c, 0x43, 0xf6, 0x25, 0x95,
		0x7a, 0xb9, 0xbf, 0xed, 0x2c, 0x90, 0x1e, 0x4f,
		0xfd, 0xfa, 0x54, 0x95, 0x04, 0x19, 0x02, 0x5e,
	})
	c.Assert(gotrax.SerializeScalar(proof.v), DeepEquals, []byte{
		0x85, 0x52, 0xb9, 0xd0, 0x72, 0x20, 0xed, 0x97,
		0x8e, 0x1a, 0xe5, 0x8f, 0x05, 0x51, 0x4a, 0x56,
		0x25, 0x08, 0xf3, 0xec, 0xd7, 0x7a, 0xbc, 0xd4,
		0x56, 0x5e, 0x21, 0x77, 0x2a, 0x2c, 0x66, 0x04,
		0x0f, 0xf9, 0xc2, 0xda, 0x5c, 0x9b, 0x24, 0x43,
		0xeb, 0xb7, 0x5f, 0xdf, 0xb4, 0x15, 0x3a, 0x99,
		0x02, 0x51, 0x94, 0x9b, 0x2b, 0x10, 0xef, 0x2f,
	})
}

func (s *GenericServerSuite) Test_generateDhProof_generatesSpecificValues(c *C) {
	wr := gotrax.ReaderIntoWithRandom(gotrax.FixtureRand())
	valuesPriv := make([]*big.Int, 5)
	valuesPriv[0] = randomDhSecretValue(wr)
	valuesPriv[1] = randomDhSecretValue(wr)
	valuesPriv[2] = randomDhSecretValue(wr)
	valuesPriv[3] = randomDhSecretValue(wr)
	valuesPriv[4] = randomDhSecretValue(wr)

	valuesPub := make([]*big.Int, 5)
	valuesPub[0] = new(big.Int).Exp(g3, valuesPriv[0], dhP)
	valuesPub[1] = new(big.Int).Exp(g3, valuesPriv[1], dhP)
	valuesPub[2] = new(big.Int).Exp(g3, valuesPriv[2], dhP)
	valuesPub[3] = new(big.Int).Exp(g3, valuesPriv[3], dhP)
	valuesPub[4] = new(big.Int).Exp(g3, valuesPriv[4], dhP)

	m := [64]byte{0x42, 0x02, 0x03}

	proof, e := generateDhProof(wr, valuesPriv, valuesPub, m[:], usageProofMessageDh, nil)
	c.Assert(e, IsNil)
	c.Assert(proof.c, DeepEquals, []byte{
		0x27, 0x74, 0x7e, 0x5c, 0x68, 0x38, 0xb, 0xd1,
		0xc9, 0x46, 0x44, 0xa1, 0x27, 0x44, 0x88, 0xde,
		0xc7, 0x41, 0x1a, 0x6e, 0xfa, 0xed, 0xf2, 0x71,
		0x4e, 0x3e, 0x37, 0x86, 0xa1, 0x3f, 0x3c, 0x6e,
		0x64, 0x1a, 0xc2, 0x7d, 0x16, 0xac, 0x2, 0x8d,
		0x59, 0xd, 0xe4, 0x8d, 0x6, 0xe2, 0xc0, 0xe4,
		0x47, 0xa2, 0xe, 0x92, 0x1e, 0x4a, 0x27, 0xba,
		0x9a, 0x1a, 0x39, 0x0, 0xa8, 0x4e, 0x47, 0x6c,
	})
	c.Assert(proof.v.Bytes(), DeepEquals, []byte{
		0x5c, 0xb9, 0x43, 0x7e, 0x78, 0x53, 0xd, 0x1c,
		0x68, 0xda, 0x77, 0x2, 0x49, 0xa8, 0x19, 0xf1,
		0xa2, 0xe3, 0x97, 0xc4, 0x2f, 0x3f, 0x6f, 0x13,
		0xae, 0xe4, 0x80, 0xf3, 0x42, 0xda, 0xff, 0x3,
		0x78, 0x65, 0x86, 0x5a, 0x8d, 0xe2, 0x5c, 0x18,
		0x4, 0x55, 0x25, 0x5a, 0x2, 0xd4, 0x57, 0x7c,
		0xf6, 0x1a, 0x30, 0xe1, 0x32, 0xf6, 0xe2, 0xbb,
		0x1c, 0x57, 0xb0, 0x6c, 0xa3, 0xbc, 0xc5, 0x5b,
		0xf4, 0x8c, 0x22, 0xb7, 0xb4, 0xbd, 0xcb, 0x73,
		0xd6, 0xa5, 0x53, 0x44, 0x61, 0x4a, 0x4, 0x7f,
		0x8a, 0x44, 0xfa, 0xf8, 0x20, 0x19, 0x50, 0x22,
		0xfb, 0x87, 0x8c, 0x73, 0x21, 0x95, 0xe, 0xed,
		0xc5, 0x95, 0xff, 0xe6, 0xa, 0x8b, 0x4f, 0x7d,
		0x2e, 0x36, 0xec, 0x6d, 0x47, 0x5a, 0x7c, 0x39,
		0xa7, 0xd5, 0xda, 0xe0, 0x3b, 0x13, 0xa5, 0x9b,
		0x66, 0x20, 0x27, 0xd8, 0x39, 0x83, 0xba, 0xf2,
		0x77, 0x4e, 0x38, 0x31, 0x4, 0xf6, 0x1, 0x5b,
		0x29, 0x5, 0x7, 0x50, 0xf1, 0x28, 0xf, 0xd7,
		0x3b, 0x4e, 0x64, 0x1e, 0x31, 0x85, 0xbf, 0xd,
		0x54, 0xb8, 0x4a, 0x73, 0xfc, 0x9c, 0xd3, 0x10,
		0xb, 0x6f, 0x65, 0x13, 0x48, 0x35, 0x2e, 0xe,
		0x78, 0x5b, 0xac, 0xda, 0xad, 0x7c, 0x9c, 0x7d,
		0x46, 0xc6, 0x82, 0x6b, 0x5b, 0x5d, 0xa3, 0xfd,
		0x85, 0x9b, 0x19, 0xc4, 0x5d, 0xf7, 0x24, 0x2c,
		0x31, 0xd9, 0xc0, 0x48, 0xbe, 0xb0, 0xcd, 0x68,
		0x9b, 0x8f, 0x6e, 0x5e, 0x6b, 0x18, 0x20, 0x96,
		0x19, 0x45, 0xd9, 0x3c, 0x73, 0x35, 0xcc, 0xe4,
		0xcd, 0x55, 0xb7, 0xba, 0x62, 0x6b, 0x82, 0x3a,
		0x42, 0x3c, 0xfd, 0x37, 0x59, 0x94, 0xb7, 0x6c,
		0x59, 0x60, 0x59, 0x39, 0x6c, 0xd0, 0x30, 0x82,
		0x62, 0x4d, 0x3a, 0x9f, 0xf8, 0x41, 0x3f, 0x73,
		0x58, 0xee, 0x78, 0x5a, 0x83, 0xa4, 0xe5, 0x46,
		0x1, 0x43, 0x30, 0xbf, 0x4, 0xd1, 0xf0, 0x95,
		0x2e, 0x45, 0xf, 0x22, 0xe4, 0xdc, 0xd3, 0x91,
		0x95, 0xbe, 0xc7, 0xe9, 0xc, 0x11, 0x9d, 0x8c,
		0x79, 0xd3, 0x2e, 0xe, 0x91, 0xfe, 0xf2, 0x9b,
		0x1a, 0x1d, 0xe8, 0x7f, 0xa, 0x85, 0xff, 0x41,
		0xd7, 0x33, 0x42, 0x7f, 0xe8, 0xd5, 0xc4, 0x9c,
		0x5f, 0x57, 0x83, 0x84, 0x3d, 0x57, 0xee, 0xee,
		0x7e, 0x14, 0x66, 0x3f, 0x20, 0x2c, 0x3a, 0xc1,
		0xe, 0x23, 0xdd, 0x91, 0x79, 0x39, 0x7b, 0x8e,
		0x7f, 0x53, 0x27, 0xc8, 0x8d, 0xc1, 0x7c, 0xc1,
		0x5a, 0x6a, 0xce, 0x85, 0x20, 0xfc, 0x33, 0xfc,
		0xa1, 0x7c, 0xb8, 0x77, 0x90, 0xcd, 0x4e, 0xba,
		0xd3, 0xa8, 0xc9, 0x84, 0xc5, 0x88, 0x76, 0x53,
		0x1a, 0x62, 0x8a, 0xce, 0x97, 0xd6, 0x17, 0x9e,
		0xed, 0xfa, 0xa8, 0xbe, 0x3f, 0x5b, 0xd3, 0xbe,
		0x97, 0x99, 0x80, 0x49, 0x8a, 0xb2, 0x39, 0xc6,
	})
}

func (s *GenericServerSuite) Test_dhProof_generatesSpecificValues2(c *C) {
	var privs [][]byte = make([][]byte, 3)
	var privsb []*big.Int = make([]*big.Int, 3)

	privs[0] = []byte{
		0x00, 0x00, 0x00, 0x4F, 0x01, 0x42, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00}

	privs[1] = []byte{
		0x00, 0x00, 0x00, 0x50, 0x22, 0x01, 0x42, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}

	privs[2] = []byte{
		0x00, 0x00, 0x00, 0x50, 0x66, 0x01, 0x42, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}

	_, privsb[0], _ = gotrax.ExtractMPI(privs[0])
	_, privsb[1], _ = gotrax.ExtractMPI(privs[1])
	_, privsb[2], _ = gotrax.ExtractMPI(privs[2])

	pubs := make([]*big.Int, 3)
	pubs[0] = new(big.Int).Exp(g3, privsb[0], dhP)
	pubs[1] = new(big.Int).Exp(g3, privsb[1], dhP)
	pubs[2] = new(big.Int).Exp(g3, privsb[2], dhP)

	m := []byte{
		0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	var wr gotrax.WithRandom = nil
	rr := func(_ gotrax.WithRandom) *big.Int {
		_, val, _ := gotrax.ExtractMPI([]byte{
			0x00, 0x00, 0x00, 0x50, 0x01, 0x02, 0x01, 0x04,
			0x01, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00})
		return val
	}

	proof, _ := generateDhProof(wr, privsb, pubs, m[:], usageProofMessageDh, rr)
	c.Assert(proof.serialize(), DeepEquals, []byte{
		// c
		0xBB, 0x35, 0x0B, 0x30, 0xFD, 0xE2, 0x68, 0x13,
		0xEA, 0xAC, 0x7C, 0xF9, 0x99, 0x3E, 0x95, 0xDE,
		0xDB, 0x2C, 0x22, 0xDE, 0x9C, 0xEE, 0x4D, 0x7D,
		0xF2, 0x8C, 0x9D, 0x60, 0xD7, 0xBF, 0xB9, 0x6A,
		0x74, 0xA8, 0x30, 0xA4, 0x01, 0x45, 0x25, 0x7B,
		0x47, 0x87, 0xC0, 0xBE, 0xAC, 0x96, 0x42, 0x0E,
		0x32, 0xD3, 0x5C, 0xF6, 0xCA, 0xA9, 0x3D, 0x26,
		0xB3, 0xC6, 0x3E, 0x6B, 0x0E, 0x7E, 0x1C, 0x90,

		// v
		0x00, 0x00, 0x00, 0x7C,
		0x52, 0x8F, 0x16, 0x63, 0x14, 0xDA, 0x1C, 0x3A,
		0xE8, 0xB6, 0xC8, 0x62, 0x53, 0xAA, 0x84, 0xFC,
		0xF3, 0xE4, 0xE4, 0x41, 0xB5, 0xC5, 0x52, 0x7F,
		0x3F, 0x62, 0x20, 0x91, 0xDE, 0x58, 0xC6, 0x28,
		0x59, 0x3E, 0x24, 0xCB, 0x51, 0x52, 0xF6, 0xB0,
		0xD1, 0x43, 0x72, 0x66, 0x52, 0x87, 0x77, 0x04,
		0x01, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00})
}

func (s *GenericServerSuite) _Test_dhProof_verify_verifiesSpecificValues(c *C) {
	pubs0 := []byte{
		0x00, 0x00, 0x01, 0x80, 0x3A, 0xF3, 0x6D, 0xBB,
		0x06, 0x33, 0x2F, 0x7B, 0x22, 0xD4, 0x93, 0x41,
		0x8D, 0xCD, 0xC7, 0x7E, 0x9B, 0xB9, 0x51, 0x89,
		0x2C, 0x3A, 0xAC, 0x99, 0x25, 0xE4, 0x4A, 0x52,
		0xF5, 0x61, 0x26, 0xDB, 0x13, 0xE1, 0x93, 0xBC,
		0x07, 0x4C, 0xB0, 0x0D, 0x8E, 0x36, 0xF2, 0x52,
		0x46, 0x3F, 0x0B, 0x4F, 0x6E, 0xF2, 0xCC, 0x99,
		0xA4, 0x92, 0xBA, 0x84, 0x71, 0xC2, 0x24, 0x3B,
		0xE3, 0x43, 0xE4, 0xB9, 0xE5, 0x6F, 0xA0, 0x74,
		0x57, 0xA5, 0x10, 0x85, 0xB9, 0x90, 0xBE, 0xAF,
		0x74, 0x38, 0x7F, 0x3F, 0x1B, 0x4E, 0xBC, 0x99,
		0x89, 0x46, 0x69, 0x56, 0x7E, 0x96, 0x86, 0x2A,
		0x50, 0x4D, 0xCF, 0x13, 0x55, 0x12, 0xCB, 0x9B,
		0x7E, 0x5D, 0xEA, 0x39, 0x68, 0x37, 0x75, 0x30,
		0xE6, 0x6E, 0x36, 0x3B, 0xDB, 0x24, 0xDA, 0x21,
		0x99, 0x11, 0x79, 0xF4, 0x7E, 0x61, 0x4B, 0x63,
		0x38, 0xB9, 0x99, 0xC6, 0x39, 0xF3, 0x46, 0xF4,
		0xD1, 0xFD, 0x70, 0x72, 0xB0, 0x81, 0xFA, 0x6F,
		0xDD, 0xD2, 0x56, 0x10, 0x9E, 0x49, 0xC0, 0xDB,
		0x16, 0x9D, 0x9F, 0xD4, 0xF0, 0xD5, 0x88, 0x82,
		0x91, 0xC5, 0xDE, 0x44, 0x64, 0x1A, 0xF3, 0x54,
		0x48, 0x20, 0xE2, 0x44, 0xF4, 0xC6, 0x2D, 0xEF,
		0xA6, 0xF2, 0x4A, 0xFD, 0xCB, 0x15, 0x26, 0x5D,
		0xEE, 0x1C, 0xFC, 0xB0, 0xFE, 0xE0, 0x13, 0x03,
		0x40, 0xA9, 0x06, 0x3F, 0x7C, 0xB4, 0xB1, 0x66,
		0x4F, 0x95, 0x10, 0xBC, 0xC7, 0x8B, 0x6F, 0xC5,
		0xB7, 0x19, 0x6B, 0x85, 0xF0, 0x75, 0x5B, 0xE3,
		0xA6, 0xA6, 0x48, 0xA4, 0xEA, 0xCD, 0xB7, 0x90,
		0x9C, 0xDD, 0xE1, 0x8F, 0x42, 0xA2, 0x3B, 0xCC,
		0x2B, 0xAA, 0xF2, 0x20, 0xFA, 0x95, 0xAC, 0x4B,
		0xBF, 0x73, 0xF1, 0x6A, 0xCA, 0x27, 0x3F, 0x00,
		0xF0, 0x86, 0x5B, 0xB3, 0x3E, 0x98, 0x10, 0x9C,
		0x1C, 0x46, 0xF3, 0x0C, 0x60, 0xBE, 0xC7, 0x0E,
		0x42, 0xCA, 0xE1, 0xE5, 0x2E, 0x63, 0x1E, 0x99,
		0x8D, 0x2D, 0x8A, 0x5A, 0x1F, 0x08, 0x76, 0x6C,
		0x8A, 0x0D, 0x6F, 0xE6, 0x7F, 0x81, 0xCD, 0x5F,
		0x40, 0x0D, 0xFF, 0xDE, 0x76, 0x2E, 0xB0, 0x23,
		0xBC, 0xF8, 0x06, 0x39, 0xB0, 0x6D, 0xC4, 0xEE,
		0x3A, 0x81, 0x1A, 0xB7, 0xE2, 0x0D, 0xC2, 0x5D,
		0xA8, 0xEB, 0x30, 0x4F, 0x98, 0xB0, 0x62, 0x28,
		0xCD, 0xEF, 0x6E, 0xA9, 0x67, 0x27, 0x7E, 0x22,
		0xC0, 0x1D, 0x82, 0xD5, 0xC1, 0x1A, 0xB7, 0xC0,
		0x48, 0x69, 0xC9, 0xD9, 0x60, 0x3A, 0x0B, 0x73,
		0x8F, 0x6B, 0x99, 0xCA, 0x15, 0xBA, 0x47, 0x5C,
		0xF4, 0xA7, 0xB7, 0x32, 0xBA, 0x83, 0x01, 0x4B,
		0x36, 0x7F, 0x64, 0x12, 0x97, 0x0F, 0xE6, 0xB5,
		0xFE, 0xAB, 0x28, 0x8F, 0x4B, 0x10, 0x2B, 0xC7,
		0x1E, 0x85, 0xFC, 0x16, 0xA7, 0x38, 0xAA, 0xDC,
		0x91, 0x16, 0x9D, 0xCC}

	pubs1 := []byte{
		0x00, 0x00, 0x01, 0x80, 0xBB, 0x2A, 0x95, 0x57,
		0x21, 0xA8, 0x24, 0x5F, 0x54, 0x17, 0xBB, 0x4B,
		0xFE, 0x7C, 0x36, 0x6A, 0x71, 0x88, 0x46, 0x1E,
		0x02, 0xA5, 0xEE, 0xC1, 0x65, 0xC1, 0xE4, 0xD8,
		0xB6, 0x95, 0xEF, 0xC7, 0xE4, 0x7C, 0xC7, 0x66,
		0x45, 0x8F, 0xAA, 0xF8, 0x40, 0x8A, 0x87, 0xD7,
		0x4A, 0xCD, 0x61, 0x37, 0x1E, 0x3C, 0x83, 0x34,
		0x36, 0x20, 0x2A, 0x1E, 0x7C, 0x4B, 0x9F, 0xAD,
		0xEB, 0x50, 0xB2, 0x15, 0xC7, 0x20, 0x66, 0x05,
		0x96, 0x46, 0x2B, 0xB2, 0x55, 0x17, 0x02, 0x26,
		0x7D, 0xD5, 0x81, 0x1B, 0xF1, 0x51, 0x5E, 0x90,
		0xBE, 0xE8, 0x68, 0x17, 0xDD, 0x4C, 0x57, 0x02,
		0x20, 0x44, 0x3D, 0xB0, 0x46, 0x9A, 0xB6, 0x97,
		0x6B, 0xC8, 0x90, 0x4B, 0x14, 0x1C, 0xAC, 0xAF,
		0xB9, 0xC3, 0xA7, 0xC2, 0x94, 0x20, 0x4F, 0x09,
		0x4A, 0xF6, 0x3A, 0xC8, 0x08, 0x06, 0x42, 0x6F,
		0xB5, 0xE8, 0x4A, 0xC8, 0xAD, 0xA0, 0x74, 0x91,
		0x90, 0x78, 0x70, 0xEB, 0x2F, 0x1C, 0xC2, 0x12,
		0x12, 0x70, 0xF9, 0x9E, 0x35, 0xAB, 0x57, 0x60,
		0xA2, 0x15, 0xFC, 0xF7, 0x22, 0x35, 0x08, 0xED,
		0x68, 0x5A, 0x53, 0x96, 0x49, 0x04, 0xCC, 0x77,
		0x96, 0x67, 0x4C, 0x3B, 0x26, 0x75, 0x7D, 0xB7,
		0x9B, 0x51, 0x9C, 0x57, 0x14, 0x39, 0x03, 0x8B,
		0x1B, 0x08, 0x2C, 0x3C, 0x98, 0x46, 0x05, 0x47,
		0x1D, 0xFC, 0xD7, 0xF7, 0x2F, 0x20, 0x78, 0x26,
		0x57, 0x22, 0xF2, 0xFD, 0xD4, 0x35, 0x08, 0x12,
		0x02, 0x81, 0xAF, 0xFD, 0xC3, 0xF8, 0xAC, 0x6A,
		0xAA, 0xE1, 0x6A, 0x32, 0x66, 0x51, 0xE2, 0x71,
		0x4C, 0x01, 0x44, 0xC7, 0x24, 0x2C, 0x8C, 0xF7,
		0x95, 0x71, 0x70, 0xF7, 0x24, 0xF7, 0xCF, 0x53,
		0x22, 0x5A, 0xFF, 0x08, 0xE8, 0x2D, 0xA9, 0x32,
		0x26, 0x25, 0x5D, 0xC4, 0x00, 0x7B, 0x11, 0x8C,
		0x98, 0x34, 0xD1, 0x2D, 0x6E, 0xF7, 0x0C, 0xE8,
		0x2C, 0xBA, 0x02, 0x6F, 0xCD, 0xD4, 0x7D, 0x3D,
		0xEB, 0x48, 0xDD, 0x08, 0xE8, 0x60, 0x16, 0x2A,
		0x63, 0xE8, 0xD8, 0xD4, 0xC4, 0xEC, 0x6C, 0x23,
		0xF6, 0xA1, 0x35, 0x8C, 0x10, 0x87, 0xEB, 0x20,
		0x03, 0xAB, 0x3E, 0x97, 0xBF, 0x7F, 0x7D, 0x2F,
		0x91, 0x49, 0x4B, 0x01, 0xE1, 0xA7, 0xEB, 0xC9,
		0xCD, 0xE5, 0xA2, 0xBD, 0x86, 0x8B, 0xDA, 0x7E,
		0xC0, 0x83, 0xDE, 0x13, 0xCA, 0x2A, 0x0C, 0x39,
		0xD3, 0xB8, 0xAD, 0xA0, 0x11, 0xE5, 0xB4, 0x24,
		0x36, 0xE1, 0x5B, 0x9A, 0x96, 0xFE, 0xC1, 0xEC,
		0x24, 0xA8, 0x92, 0xFD, 0xA8, 0xF1, 0xB8, 0x4C,
		0x6B, 0xBF, 0xFA, 0xB1, 0xBE, 0xF8, 0x5A, 0x79,
		0x87, 0xB4, 0xEF, 0x41, 0x8D, 0x0C, 0x7A, 0xB9,
		0x22, 0x7E, 0xF7, 0xD4, 0xE0, 0xAD, 0xDA, 0xA9,
		0x9F, 0x35, 0x54, 0x46, 0x6A, 0x08, 0x8E, 0x17,
		0xFF, 0x25, 0x1B, 0x29}

	pubs2 := []byte{
		0x00, 0x00, 0x01, 0x80, 0xE8, 0x9B, 0x8F, 0x50,
		0x77, 0xBB, 0x02, 0xD8, 0x4D, 0x07, 0x4D, 0x2D,
		0x8D, 0x70, 0xC6, 0x80, 0x3E, 0x42, 0x70, 0x07,
		0xE4, 0xD0, 0xAC, 0xCD, 0xAA, 0xB4, 0xDD, 0xB9,
		0x44, 0xBC, 0x28, 0x56, 0xE4, 0x2B, 0x29, 0x72,
		0x90, 0x47, 0x51, 0x8B, 0x69, 0x99, 0xCC, 0xEE,
		0x1D, 0x64, 0x35, 0x0D, 0x8A, 0xEA, 0xA0, 0x61,
		0xC7, 0x20, 0x8E, 0x4C, 0x07, 0x38, 0xCE, 0x37,
		0x71, 0x5E, 0x76, 0x71, 0xC7, 0x1D, 0x41, 0x4D,
		0x79, 0x2C, 0xE0, 0x8C, 0x49, 0x7F, 0x3E, 0xDA,
		0xB2, 0x10, 0x4B, 0xFD, 0xC3, 0x33, 0xDD, 0x7B,
		0xEB, 0xCF, 0x48, 0x1A, 0xBD, 0xF3, 0x6E, 0xAA,
		0xF6, 0xDA, 0x9B, 0x6E, 0x1F, 0x5C, 0xCE, 0x89,
		0x4F, 0x9C, 0x4B, 0xDD, 0x26, 0x37, 0xB1, 0x9C,
		0xA1, 0xD5, 0x5B, 0x3A, 0x59, 0xE9, 0xB1, 0x4A,
		0x08, 0x91, 0x97, 0xAB, 0x88, 0xB1, 0x7E, 0x87,
		0x3E, 0x7D, 0x7E, 0x20, 0xD4, 0xD0, 0x19, 0x07,
		0xC0, 0x31, 0x12, 0xD5, 0x62, 0x38, 0xB6, 0x73,
		0xD2, 0x92, 0x3B, 0xCE, 0xFB, 0xD3, 0x54, 0x49,
		0xE4, 0x0C, 0xCF, 0x14, 0xD2, 0x18, 0x80, 0x4C,
		0x9D, 0x14, 0x15, 0x33, 0x2A, 0x94, 0x22, 0xCE,
		0xB3, 0x0C, 0x21, 0xA3, 0x4E, 0x71, 0xF8, 0x3A,
		0x86, 0xA5, 0x09, 0x65, 0xF0, 0x28, 0x1F, 0xFB,
		0x17, 0x23, 0x11, 0xD9, 0x2A, 0x81, 0x3D, 0xB3,
		0x23, 0xB6, 0xD0, 0x3B, 0x85, 0x2E, 0x96, 0x0E,
		0x82, 0x98, 0x3A, 0x18, 0xB1, 0x41, 0x77, 0x71,
		0xA4, 0x72, 0x5A, 0x60, 0xE6, 0xFA, 0x6F, 0x5C,
		0x1D, 0x96, 0x34, 0xC9, 0x8D, 0x6C, 0xD5, 0x8C,
		0x44, 0xCF, 0xA3, 0xC9, 0x92, 0x44, 0x98, 0x09,
		0x3B, 0x23, 0x24, 0x2D, 0xB8, 0x9D, 0x47, 0xFF,
		0x96, 0xBD, 0xC9, 0xCD, 0xC2, 0xF3, 0x79, 0x10,
		0xDD, 0xCD, 0xE6, 0xD8, 0xDF, 0x62, 0x7B, 0x33,
		0xB8, 0x65, 0x6B, 0x62, 0xCE, 0x9D, 0x39, 0xDC,
		0xDC, 0xA0, 0x26, 0x66, 0xEF, 0xFA, 0x34, 0x2C,
		0x6E, 0x5C, 0x3C, 0xD1, 0xC6, 0x21, 0xA2, 0xE8,
		0x87, 0xCA, 0x11, 0x0C, 0x7F, 0x03, 0xE9, 0x20,
		0xAE, 0xCD, 0x90, 0xB7, 0x0C, 0xF7, 0xAD, 0xB3,
		0x17, 0x55, 0x83, 0x64, 0x8E, 0x83, 0x3A, 0x4D,
		0x30, 0x25, 0xD8, 0xB0, 0xCC, 0x89, 0xA7, 0x32,
		0xE0, 0x50, 0x8B, 0x18, 0xF5, 0x93, 0x02, 0x55,
		0x5D, 0x53, 0xA5, 0xDF, 0xD5, 0x8D, 0x00, 0x2C,
		0xA4, 0x54, 0xDF, 0x28, 0x18, 0xA0, 0xD5, 0x54,
		0x52, 0x8D, 0x55, 0x17, 0xFE, 0x5F, 0xEA, 0x02,
		0x6E, 0x06, 0x2B, 0xC3, 0x41, 0x21, 0xA4, 0x62,
		0x68, 0x2A, 0x49, 0x89, 0x78, 0xDB, 0x3D, 0x0E,
		0x41, 0xDC, 0x7C, 0xA2, 0xDC, 0x0B, 0xC9, 0xB2,
		0x78, 0xCB, 0xBA, 0x0F, 0xD5, 0xA1, 0xF7, 0x15,
		0x79, 0x13, 0xB0, 0xAF, 0x49, 0xCC, 0x5A, 0x98,
		0xDD, 0xC2, 0x3B, 0x8D}

	m := []byte{
		0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	px := []byte{
		0xE7, 0xF0, 0x42, 0x20, 0x19, 0x2B, 0x31, 0x39,
		0x84, 0xB0, 0x32, 0x41, 0xEB, 0xE8, 0x80, 0x72,
		0x97, 0x1E, 0x16, 0xAC, 0xB9, 0x13, 0x94, 0x03,
		0x64, 0x16, 0xE4, 0x38, 0xB3, 0x2E, 0xB3, 0x66,
		0xF6, 0x81, 0xFB, 0xC8, 0xC4, 0x85, 0xAB, 0x28,
		0xB6, 0x77, 0xD4, 0xB9, 0x97, 0x4F, 0x70, 0x70,
		0xF2, 0xC8, 0xA2, 0x93, 0x5C, 0xCE, 0xA6, 0xD4,
		0xFD, 0x88, 0x91, 0xD1, 0xA1, 0xE7, 0x49, 0xD3,
		0x00, 0x00, 0x00, 0x7C, 0x72, 0xAA, 0x65, 0x1B,
		0x7F, 0x3F, 0xB0, 0x94, 0xA8, 0x63, 0x19, 0x67,
		0x39, 0xCF, 0x1F, 0xD5, 0xAF, 0x9F, 0x96, 0xC7,
		0xFE, 0x63, 0xB8, 0x1B, 0xB8, 0xA2, 0x55, 0x6F,
		0x9E, 0x0B, 0x55, 0xBF, 0x2D, 0xC9, 0xFB, 0x3A,
		0x15, 0xAA, 0x42, 0x37, 0x4B, 0xC1, 0xDE, 0x26,
		0x0E, 0x57, 0x3E, 0xEA, 0xE5, 0xD6, 0x3E, 0xD2,
		0xF0, 0x7F, 0x82, 0x83, 0x4B, 0xB3, 0x19, 0x1D,
		0x2D, 0x35, 0xE5, 0x1A, 0xEF, 0x99, 0x43, 0x70,
		0x57, 0x18, 0xDC, 0x42, 0x96, 0x8A, 0xD9, 0xB8,
		0xD3, 0x34, 0x44, 0x9F, 0x33, 0x8A, 0x5D, 0x27,
		0x15, 0x75, 0x13, 0xBB, 0x55, 0x85, 0x96, 0x86,
		0x98, 0x02, 0x00, 0xBE, 0x5B, 0x28, 0xAC, 0xB9,
		0x2F, 0xDE, 0x6D, 0x1B, 0x27, 0xAB, 0x53, 0x01,
		0x70, 0x89, 0x0F, 0x54, 0x24, 0x48, 0xA2, 0xC1,
		0x7F, 0x92, 0x94, 0x88, 0x63, 0xC2, 0xF5, 0x2F}

	_, p0, _ := gotrax.ExtractMPI(pubs0)
	_, p1, _ := gotrax.ExtractMPI(pubs1)
	_, p2, _ := gotrax.ExtractMPI(pubs2)
	pubs := []*big.Int{p0, p1, p2}
	proof := &dhProof{}
	proof.deserialize(px)

	c.Assert(proof.verify(pubs, m, 0x13), Equals, true)
}

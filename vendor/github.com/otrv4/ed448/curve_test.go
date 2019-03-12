package ed448

import (
	"bytes"
	"encoding/hex"

	. "gopkg.in/check.v1"
)

func (s *Ed448Suite) Test_RadixBasePointIsOnCurve(c *C) {
	c.Assert(testPoint().OnCurve(), Equals, true)
}

func (s *Ed448Suite) Test_RadixMultiplyByBase(c *C) {
	scalar := scalar{}
	scalar[scalarWords-1] = 1000

	p := curve.multiplyByBase(scalar)

	c.Assert(p.OnCurve(), Equals, true)
}

func (s *Ed448Suite) Test_RadixGenerateKey(c *C) {
	buffer := make([]byte, symKeyBytes)
	buffer[0] = 0x10
	r := bytes.NewReader(buffer[:])

	privKey, err := curve.generateKey(r)

	expectedSymKey := make([]byte, symKeyBytes)
	expectedSymKey[0] = 0x10

	expectedPriv := []byte{
		0x06, 0x01, 0x3f, 0x3e, 0xb3, 0x3f, 0x9e, 0x10,
		0xde, 0xde, 0x34, 0x23, 0x6a, 0x9a, 0x75, 0x44,
		0x69, 0x41, 0x18, 0x4f, 0x79, 0xb7, 0x52, 0x50,
		0x03, 0xa0, 0x7d, 0xe2, 0x89, 0xee, 0x15, 0x8a,
		0xaf, 0x44, 0xf3, 0x39, 0x78, 0x2c, 0xa6, 0x9b,
		0xbe, 0x5b, 0xb4, 0x1d, 0x25, 0x6a, 0x83, 0x32,
		0x7c, 0xd0, 0xc0, 0x3d, 0xa5, 0x26, 0xf8, 0x37,
	}

	expectedPublic := []byte{
		0x4d, 0xdb, 0xad, 0x93, 0xb8, 0x95, 0x29, 0x61,
		0x67, 0xfc, 0xf4, 0xbd, 0x27, 0x94, 0xb9, 0x0f,
		0x06, 0x09, 0x05, 0xef, 0x8f, 0x32, 0x63, 0x2c,
		0xa6, 0xce, 0x45, 0xfb, 0x1c, 0x83, 0xc5, 0xe7,
		0x0f, 0xf9, 0xf4, 0x43, 0x2a, 0x0c, 0xaf, 0x82,
		0x7a, 0xf5, 0x19, 0xe9, 0x5e, 0x40, 0x17, 0x48,
		0x44, 0xb9, 0xf8, 0x11, 0x88, 0x9a, 0xc3, 0xa5,
	}

	c.Assert(err, IsNil)
	c.Assert(privKey.symKey(), DeepEquals, expectedSymKey)
	c.Assert(privKey.secretKey(), DeepEquals, expectedPriv)
	c.Assert(privKey.publicKey(), DeepEquals, expectedPublic)
}

func (s *Ed448Suite) Test_DeriveNonce(c *C) {
	msg := []byte("hey there")
	symKey := [symKeyBytes]byte{
		0x27, 0x54, 0xcd, 0xa7, 0x12, 0x98, 0x88, 0x3d,
		0x4e, 0xf5, 0x11, 0x23, 0x92, 0x74, 0xb8, 0xa7,
		0xef, 0x7e, 0x51, 0x7e, 0x31, 0x28, 0xd4, 0xf7,
		0xfc, 0xfd, 0x9c, 0x62, 0xff, 0x65, 0x09, 0x65,
	}

	expectedNonce := scalar{
		0xc7a99dbd, 0xb92054cc,
		0x79b10a3e, 0x38afe6b9,
		0x859aa259, 0x007e0791,
		0x91958009, 0x1ed45cd0,
		0xbbfa381b, 0x1f427b27,
		0xb194eb5c, 0x501789df,
		0x1616d689, 0x17db93b0,
	}

	nonce := deriveNonce(msg, symKey[:])

	c.Assert(nonce, DeepEquals, expectedNonce)
}

func (s *Ed448Suite) Test_DeriveChallenge(c *C) {
	msg := []byte("hey there")
	pubKey := [pubKeyBytes]byte{
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	}
	tmpSignature := [fieldBytes]uint8{
		0x66, 0x86, 0x04, 0xa8, 0x71, 0x4c, 0x39, 0xb9,
		0x42, 0x01, 0x7b, 0x45, 0xb6, 0xc7, 0xaf, 0xdb,
		0x7c, 0xad, 0x1f, 0x80, 0xa0, 0x23, 0x4d, 0xb5,
		0xab, 0x7c, 0x55, 0xf4, 0x38, 0x7d, 0xab, 0x60,
		0x25, 0x5a, 0x3d, 0xc9, 0xa1, 0x89, 0x85, 0xd1,
		0xc7, 0x4b, 0x19, 0x39, 0xbb, 0x08, 0x49, 0x09,
		0x0e, 0x0a, 0x31, 0x5a, 0x05, 0x5d, 0xe6, 0x47,
	}

	expectedChallenge := &scalar{
		0x6c226d73, 0x70edcfc3, 0x44156c47, 0x084f4695,
		0xe72606ac, 0x9d0ce5e5, 0xed96d3ba, 0x9ff3fa11,
		0x4a15c383, 0xca38a0af, 0xead789b3, 0xb96613ba,
		0x48ba4461, 0x34eb2031,
	}

	challenge := deriveChallenge(pubKey[:], tmpSignature, msg)

	c.Assert(challenge, DeepEquals, expectedChallenge)
}

func (s *Ed448Suite) Test_Sign(c *C) {
	msg := []byte("hey there")
	k := privateKey([privKeyBytes]byte{
		//secret
		0x1f, 0x44, 0xfd, 0x2e, 0xde, 0x47, 0xca, 0xa8,
		0x7c, 0x4c, 0x45, 0x88, 0x1a, 0x7e, 0x01, 0x5a,
		0xa9, 0x01, 0x37, 0xfb, 0x0d, 0xbe, 0xb9, 0xe0,
		0xeb, 0x47, 0x29, 0xf7, 0x74, 0x0b, 0x5c, 0x23,
		0x66, 0xaa, 0xfd, 0x39, 0x03, 0x38, 0x78, 0x80,
		0x8f, 0xb2, 0x06, 0x13, 0x4e, 0xfb, 0xcf, 0x02,
		0x11, 0x43, 0x11, 0x3a, 0xd1, 0xf8, 0xb8, 0x22,

		//public
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,

		//symmetric
		0x27, 0x54, 0xcd, 0xa7, 0x12, 0x98, 0x88, 0x3d,
		0x4e, 0xf5, 0x11, 0x23, 0x92, 0x74, 0xb8, 0xa7,
		0xef, 0x7e, 0x51, 0x7e, 0x31, 0x28, 0xd4, 0xf7,
		0xfc, 0xfd, 0x9c, 0x62, 0xff, 0x65, 0x09, 0x65,
	})

	expectedSignature := [signatureBytes]byte{
		0x66, 0x86, 0x04, 0xa8, 0x71, 0x4c, 0x39, 0xb9,
		0x42, 0x01, 0x7b, 0x45, 0xb6, 0xc7, 0xaf, 0xdb,
		0x7c, 0xad, 0x1f, 0x80, 0xa0, 0x23, 0x4d, 0xb5,
		0xab, 0x7c, 0x55, 0xf4, 0x38, 0x7d, 0xab, 0x60,
		0x25, 0x5a, 0x3d, 0xc9, 0xa1, 0x89, 0x85, 0xd1,
		0xc7, 0x4b, 0x19, 0x39, 0xbb, 0x08, 0x49, 0x09,
		0x0e, 0x0a, 0x31, 0x5a, 0x05, 0x5d, 0xe6, 0x47,
		0xc6, 0xb8, 0x18, 0x21, 0xd5, 0xac, 0x56, 0x43,
		0x3c, 0xe7, 0xd7, 0x26, 0xb7, 0x74, 0x91, 0x45,
		0x31, 0xea, 0x0b, 0xf1, 0xbb, 0x28, 0xe5, 0x83,
		0x95, 0xd6, 0x60, 0xb9, 0x28, 0x7e, 0xda, 0xd0,
		0xa1, 0xf9, 0xd7, 0xba, 0x01, 0xba, 0xf5, 0xe9,
		0x18, 0x15, 0xea, 0x94, 0xca, 0x8c, 0xc5, 0x12,
		0xeb, 0x76, 0x2c, 0x30, 0x3e, 0x36, 0xd0, 0x3b,
	}

	signature, err := curve.sign(msg, &k)

	c.Assert(err, IsNil)
	c.Assert(signature, DeepEquals, expectedSignature)
}

func (s *Ed448Suite) Test_Verify(c *C) {
	msg := []byte("hey there")
	k := publicKey([pubKeyBytes]byte{
		//public
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	})
	signature := [signatureBytes]byte{
		0x66, 0x86, 0x04, 0xa8, 0x71, 0x4c, 0x39, 0xb9,
		0x42, 0x01, 0x7b, 0x45, 0xb6, 0xc7, 0xaf, 0xdb,
		0x7c, 0xad, 0x1f, 0x80, 0xa0, 0x23, 0x4d, 0xb5,
		0xab, 0x7c, 0x55, 0xf4, 0x38, 0x7d, 0xab, 0x60,
		0x25, 0x5a, 0x3d, 0xc9, 0xa1, 0x89, 0x85, 0xd1,
		0xc7, 0x4b, 0x19, 0x39, 0xbb, 0x08, 0x49, 0x09,
		0x0e, 0x0a, 0x31, 0x5a, 0x05, 0x5d, 0xe6, 0x47,
		0xc6, 0xb8, 0x18, 0x21, 0xd5, 0xac, 0x56, 0x43,
		0x3c, 0xe7, 0xd7, 0x26, 0xb7, 0x74, 0x91, 0x45,
		0x31, 0xea, 0x0b, 0xf1, 0xbb, 0x28, 0xe5, 0x83,
		0x95, 0xd6, 0x60, 0xb9, 0x28, 0x7e, 0xda, 0xd0,
		0xa1, 0xf9, 0xd7, 0xba, 0x01, 0xba, 0xf5, 0xe9,
		0x18, 0x15, 0xea, 0x94, 0xca, 0x8c, 0xc5, 0x12,
		0xeb, 0x76, 0x2c, 0x30, 0x3e, 0x36, 0xd0, 0x3b,
	}

	valid := curve.verify(signature, msg, &k)

	c.Assert(valid, Equals, true)
}

func (s *Ed448Suite) Test_MultiplyMontgomery(c *C) {
	pk := mustDeserialize(serialized{
		0x0e, 0xe8, 0x29, 0x1c, 0xc5, 0x9d, 0x51, 0x9c,
		0xb2, 0x94, 0xdd, 0xc4, 0x5c, 0xb9, 0xf7, 0x0f,
		0xd1, 0xd9, 0x3e, 0x4c, 0x45, 0x55, 0x15, 0x70,
		0x84, 0x4d, 0x2e, 0x18, 0xad, 0x99, 0xc4, 0xf9,
		0xfe, 0xc7, 0xe8, 0x6f, 0x5c, 0xda, 0xac, 0xe9,
		0x55, 0xff, 0x42, 0x75, 0x52, 0x6c, 0x04, 0xb6,
		0xe1, 0xc8, 0x49, 0xb9, 0xc1, 0x86, 0x37, 0xd0,
	})

	sk := scalar{
		0x2efd441f, 0xa8ca47de, 0x88454c7c, 0x5a017e1a,
		0xfb3701a9, 0xe0b9be0d, 0xf72947eb, 0x235c0b74,
		0x39fdaa66, 0x80783803, 0x1306b28f, 0x02cffb4e,
		0x3a114311, 0x22b8f8d1,
	}

	bs, _ := hex.DecodeString("322d71661943b5e080abed64d9ed331874a975329aaf9b42815e793ac08691e478fe559b29593a5413d5a4475e3ae0735a6d9bc1dc192b7d")
	expectedPublic := new(bigNumber)
	expectedPublic.setBytes(bs)

	pk, ok := curve.multiplyMontgomery(pk, sk, scalarBits, 1)

	c.Assert(ok, Equals, word(0))
	c.Assert(pk, DeepEquals, expectedPublic)
}

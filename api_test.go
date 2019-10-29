package prekeyserver

import (
	"bytes"
	"crypto/rand"
	"errors"
	"os"
	"time"

	"github.com/otrv4/gotrx"
	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_realFactory_RandReader_returnsTheDefaultRandReader(c *C) {
	r := &realFactory{}
	c.Assert(r.RandReader(), Equals, rand.Reader)
}

func (s *GenericServerSuite) Test_realFactory_RandReader_returnsTheGivenReader(c *C) {
	f := gotrx.FixtureRand()
	r := &realFactory{r: f}
	c.Assert(r.RandReader(), Equals, f)
}

func (s *GenericServerSuite) Test_CreateFactory_returnsARealFactoryWithTheGivenRandomness(c *C) {
	f := gotrx.FixtureRand()
	fact := CreateFactory(f)
	rf, ok := fact.(*realFactory)
	c.Assert(ok, Equals, true)
	c.Assert(rf.r, Equals, f)
}

func (s *GenericServerSuite) Test_inMemoryStorageFactory_createStorage_returnsAnInMemoryStorageFactory(c *C) {
	res := (&inMemoryStorageFactory{}).createStorage()
	c.Assert(res, Not(IsNil))
	c.Assert(res, FitsTypeOf, &inMemoryStorage{})
}

func (s *GenericServerSuite) Test_realFactory_LoadStorageType_returnsInMemoryStorage(c *C) {
	res, _ := (&realFactory{}).LoadStorageType("in-memory")
	c.Assert(res, Not(IsNil))
	c.Assert(res, FitsTypeOf, &inMemoryStorageFactory{})
}

func (s *GenericServerSuite) Test_realFactory_LoadStorageType_returnsFileStorage(c *C) {
	os.Mkdir(testDir, 0700)
	defer os.RemoveAll(testDir)

	res, _ := (&realFactory{}).LoadStorageType("dir:" + testDir)
	c.Assert(res, Not(IsNil))
	c.Assert(res, FitsTypeOf, &fileStorageFactory{})
}

func (s *GenericServerSuite) Test_realFactory_LoadStorageType_returnsErrorForNonExistantDirectory(c *C) {
	_, e := (&realFactory{}).LoadStorageType("dir:unknown/dir/please/dont/create")
	c.Assert(e, ErrorMatches, "directory doesn't exist")
}

func (s *GenericServerSuite) Test_realFactory_LoadStorageType_givesErrorForUnknownStorageType(c *C) {
	res, e := (&realFactory{}).LoadStorageType("unknown-storage-please-don't-create")
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "unknown storage type")
}

func (s *GenericServerSuite) Test_realFactory_CreateKeypair_createsAKeypairFromTheGivenRandomness(c *C) {
	r := gotrx.FixtureRand()
	res := (&realFactory{r: r}).CreateKeypair()
	c.Assert(res, Not(IsNil))
	c.Assert(res, FitsTypeOf, &gotrx.Keypair{})
	c.Assert(res.(*gotrx.Keypair).Sym[:], DeepEquals, []byte{
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd, 0xab, 0xcd,
		0xab})
}

func (s *GenericServerSuite) Test_realFactory_NewServer_createsAServerWithTheGivenValues(c *C) {
	f := &realFactory{r: gotrx.FixtureRand()}
	kp := f.CreateKeypair()
	mockCalled := false
	mockRestrictor := func(string) bool {
		mockCalled = true
		return false
	}
	res := f.NewServer("foobar", kp, 42, &inMemoryStorageFactory{}, time.Duration(25), time.Duration(77), mockRestrictor)
	c.Assert(res, Not(IsNil))
	c.Assert(res, FitsTypeOf, &GenericServer{})
	gs := res.(*GenericServer)
	c.Assert(gs.identity, Equals, "foobar")
	c.Assert(gs.fingerprint, DeepEquals, gotrx.Fingerprint(kp.(*gotrx.Keypair).Fingerprint()))
	c.Assert(gs.key, Equals, kp)
	c.Assert(gs.fragLen, Equals, 42)
	c.Assert(gs.fragmentations, Not(IsNil))
	c.Assert(gs.sessions, Not(IsNil))
	c.Assert(gs.storageImpl, Not(IsNil))
	c.Assert(gs.sessionTimeout, Equals, time.Duration(25))
	c.Assert(gs.fragmentationTimeout, Equals, time.Duration(77))
	c.Assert(gs.messageHandler, Not(IsNil))
	c.Assert(gs.messageHandler.(*otrngMessageHandler).s, Equals, gs)
	c.Assert(gs.rest("bla"), Equals, false)
	c.Assert(mockCalled, Equals, true)

}

func (s *GenericServerSuite) Test_realFactory_NewServer_setsANullRestrictorIfNoneIsGiven(c *C) {
	f := &realFactory{r: gotrx.FixtureRand()}
	kp := f.CreateKeypair()
	res := f.NewServer("foobar", kp, 42, &inMemoryStorageFactory{}, time.Duration(25), time.Duration(77), nil)
	gs := res.(*GenericServer)
	c.Assert(gs.rest("bla"), Equals, false)
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_decodesACorrectMessage(c *C) {
	sym := [57]byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}

	expectedKp := gotrx.DeriveKeypair(sym)

	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA",
	}

	kp, e := kis.intoKeypair()
	c.Assert(e, IsNil)
	c.Assert(kp.Sym, DeepEquals, expectedKp.Sym)
	c.Assert(kp.Pub.K().Equals(expectedKp.Pub.K()), Equals, true)
	c.Assert(kp.Priv.K().Equals(expectedKp.Priv.K()), Equals, true)
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_generatesAnErrorForBadBase64OnSymmetric(c *C) {
	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA",
	}

	_, e := kis.intoKeypair()
	c.Assert(e, ErrorMatches, "couldn't decode symmetric key")
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_generatesAnErrorForBadBase64OnPrivate(c *C) {
	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA",
	}

	_, e := kis.intoKeypair()
	c.Assert(e, ErrorMatches, "couldn't decode private key")
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_generatesAnErrorForBadBase64OnPublic(c *C) {
	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0u",
	}

	_, e := kis.intoKeypair()
	c.Assert(e, ErrorMatches, "couldn't decode public key")
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_generatesAnErrorForBadPrivateScalar(c *C) {
	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04Sop",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA",
	}

	_, e := kis.intoKeypair()
	c.Assert(e, ErrorMatches, "couldn't decode scalar for private key")
}

func (s *GenericServerSuite) Test_keypairInStorage_intoKeypair_generatesAnErrorForBadPublicPoint(c *C) {
	kis := &keypairInStorage{
		Symmetric: "AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Private:   "S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=",
		Public:    "BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pw",
	}

	_, e := kis.intoKeypair()
	c.Assert(e, ErrorMatches, "couldn't decode point for public key")
}

func (s *GenericServerSuite) Test_realFactory_LoadKeypairFrom_canLoadAKeypairCorrectly(c *C) {
	sym := [57]byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}

	expectedKp := gotrx.DeriveKeypair(sym)

	b := bytes.NewBufferString("{\"Symmetric\":\"AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\"," +
		"\"Private\":\"S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=\"," +
		"\"Public\":\"BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA\"}\n")
	f := &realFactory{}
	kp, e := f.LoadKeypairFrom(b)
	c.Assert(e, IsNil)
	c.Assert(kp.(*gotrx.Keypair).Sym, DeepEquals, expectedKp.Sym)
	c.Assert(kp.(*gotrx.Keypair).Pub.K().Equals(expectedKp.Pub.K()), Equals, true)
	c.Assert(kp.(*gotrx.Keypair).Priv.K().Equals(expectedKp.Priv.K()), Equals, true)
}

type erroringReader struct{}

func (*erroringReader) Read([]byte) (int, error) {
	return 0, errors.New("something bad")
}

func (s *GenericServerSuite) Test_realFactory_LoadKeypairFrom_willReturnAnErrorFromReading(c *C) {
	f := &realFactory{}
	kp, e := f.LoadKeypairFrom(&erroringReader{})
	c.Assert(kp, IsNil)
	c.Assert(e, ErrorMatches, "something bad")
}

func (s *GenericServerSuite) Test_realFactory_LoadKeypairFrom_willReturnAnErrorFromParsing(c *C) {
	b := bytes.NewBufferString("{\"Symmetric\":\"AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\"," +
		"\"Private\":\"S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA\"," +
		"\"Public\":\"BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA\"}\n")
	f := &realFactory{}
	kp, e := f.LoadKeypairFrom(b)
	c.Assert(kp, IsNil)
	c.Assert(e, ErrorMatches, "couldn't decode private key")
}

func (s *GenericServerSuite) Test_realFactory_StoreKeysInto_willPrintTheExpectedJson(c *C) {
	sym := [57]byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}
	rf := &realFactory{}

	kp := gotrx.DeriveKeypair(sym)
	var b bytes.Buffer
	rf.StoreKeysInto(kp, &b)

	c.Assert(b.String(), Equals,
		"{\"Symmetric\":\"AQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA\","+
			"\"Private\":\"S0Cr1lAIHXdTixCTeWQAQRJksS0o9Ftr/EcO0yemXi9fJOTAWj+c9h9QVW5M0KDm9uH04SopxiA=\","+
			"\"Public\":\"BXLBTLEwd0S5Lzg3+ZY5q7sg/8Rx2J2dVFNJ3HAOASdJHwYTFPr4moXHB2C9AYilZp0aQ5Pwg0uA\"}\n")
}

type erroringWriter struct{}

func (*erroringWriter) Write([]byte) (int, error) {
	return 0, errors.New("something bad")
}

func (s *GenericServerSuite) Test_realFactory_StoreKeysInto_willReturnAnyErrorEncountered(c *C) {
	sym := [57]byte{
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}
	rf := &realFactory{}

	kp := gotrx.DeriveKeypair(sym)
	e := rf.StoreKeysInto(kp, &erroringWriter{})

	c.Assert(e, ErrorMatches, "something bad")
}

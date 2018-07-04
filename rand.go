package prekeyserver

import (
	"crypto/rand"
	"io"
)

type WithRandom interface {
	randReader() io.Reader
}

func (g *GenericServer) randReader() io.Reader {
	if g.rand != nil {
		return g.rand
	}
	return rand.Reader
}

func randomInto(r WithRandom, b []byte) error {
	if _, err := io.ReadFull(r.randReader(), b); err != nil {
		return errShortRandomRead
	}
	return nil
}

type realRandom struct{}

func defaultRandom() *realRandom {
	return &realRandom{}
}

func (*realRandom) randReader() io.Reader {
	return rand.Reader
}

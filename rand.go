package pooh

import (
	"math/rand"
	"time"

	"github.com/MichaelTJones/pcg"
)

var Rand = rand.New(newPcgSource(time.Now().UnixNano()))

type pcgSource struct {
	*pcg.PCG64
}

var _ rand.Source64 = &pcgSource{}

func newPcgSource(seed int64) *pcgSource {
	src := &pcgSource{
		PCG64: pcg.NewPCG64(),
	}
	src.Seed(seed)
	return src
}

func (p *pcgSource) Int63() int64 {
	return int64(p.PCG64.Random() >> 1)
}

func (p *pcgSource) Uint64() uint64 {
	return p.PCG64.Random()
}

func (p *pcgSource) Seed(seed int64) {
	p.PCG64.Seed(uint64(seed), 73171451331446857877, 5049815677178091608, uint64(^seed))
}

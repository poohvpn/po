package pooh

type Modulator interface {
	Modulate(a *Analog) error
}

type Demodulator interface {
	Demodulate(a *Analog) error
}

package pooh

type Demodulator interface {
	Demodulate(a *Analog) error
}

type Modulator interface {
	Modulate(a *Analog) error
}

package pooh

type Modulator interface {
	Modulate(a *Analog) error
	Proto() string
}

type Demodulator interface {
	Demodulate(a *Analog) error
}

package pooh

type Modulator interface {
	Modulate(a *Analog) (err error, ok bool)
	Proto() string
}

type Demodulator interface {
	Demodulate(a *Analog) (err error, ok bool)
}

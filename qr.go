package qr

type Qr struct {
	versioner *versioner
	encoder   *encoder
	moduler   *moduler
}

func New(*Qr) *Qr {
	return &Qr{
		versioner: newVersioner(),
		encoder:   newEncoder(),
		moduler:   nil,
	}
}

func (qr *Qr) Encode() {
}

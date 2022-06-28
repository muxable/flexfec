package fech

type FecHeader interface {
	Marshal() []byte
	Unmarshal(buf []byte)
}

package cpu

type CPU struct {
	memory [4096]byte
	v      [16]byte
	i      uint16
	st     byte
	dt     byte
	pc     uint16
	sp     byte
	stack  [16]uint16
}

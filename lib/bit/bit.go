// all of these are copilot generated, be careful :)

package bitmask

type UnsignedNumber interface {
	uint | uint8 | uint16 | uint32 | uint64
}

func SetBit[T UnsignedNumber](bitmask T, bit T) T {
	return bitmask | (1 << bit)
}

func ClearBit[T UnsignedNumber](bitmask T, bit T) T {
	return bitmask & ^(1 << bit)
}

func ToggleBit[T UnsignedNumber](bitmask T, bit T) T {
	return bitmask ^ (1 << bit)
}

func HasBit[T UnsignedNumber](bitmask T, bit T) bool {
	return bitmask&(1<<bit) != 0
}

func SetBitMask[T UnsignedNumber](bitmask T, bitMask T) T {
	return bitmask | bitMask
}

func ClearBitMask[T UnsignedNumber](bitmask T, bitMask T) T {
	return bitmask & ^bitMask
}

func ToggleBitMask[T UnsignedNumber](bitmask T, bitMask T) T {
	return bitmask ^ bitMask
}

func HasBitMask[T UnsignedNumber](bitmask T, bitMask T) bool {
	return bitmask&bitMask != 0
}

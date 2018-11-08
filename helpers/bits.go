package helpers

//Bits is the bit type (mask)
type Bits int64

//represents the possible bit positions
const (
	F0 Bits = 1 << iota
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
)

//Set sets a bit in the mask b
func Set(b, flag Bits) Bits { return b | flag }

//Clear clears a bit in the mask b
func Clear(b, flag Bits) Bits { return b &^ flag }

//Toggle toggles a bit in the mask b
func Toggle(b, flag Bits) Bits { return b ^ flag }

//Has checks is a bit in the mask b is set
func Has(b, flag Bits) bool { return b&flag != 0 }

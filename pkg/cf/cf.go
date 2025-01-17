package cf

// Min function to return the smaller of two numbers
func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

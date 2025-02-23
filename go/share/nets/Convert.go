package nets

func Bytes2Long(data []byte) int64 {
	v1 := int64(data[0]) << 56
	v2 := int64(data[1]) << 48
	v3 := int64(data[2]) << 40
	v4 := int64(data[3]) << 32
	v5 := int64(data[4]) << 24
	v6 := int64(data[5]) << 16
	v7 := int64(data[6]) << 8
	v8 := int64(data[7])
	return v1 + v2 + v3 + v4 + v5 + v6 + v7 + v8
}

func Long2Bytes(s int64) []byte {
	size := make([]byte, 8)
	size[7] = byte(s)
	size[6] = byte(s >> 8)
	size[5] = byte(s >> 16)
	size[4] = byte(s >> 24)
	size[3] = byte(s >> 32)
	size[2] = byte(s >> 40)
	size[1] = byte(s >> 48)
	size[0] = byte(s >> 56)
	return size
}

func Bytes2Int(data []byte) int32 {
	v1 := int32(data[0]) << 24
	v2 := int32(data[1]) << 16
	v3 := int32(data[2]) << 8
	v4 := int32(data[3])
	return v1 + v2 + v3 + v4
}

func Int2Bytes(s int32) []byte {
	size := make([]byte, 4)
	size[3] = byte(s)
	size[2] = byte(s >> 8)
	size[1] = byte(s >> 16)
	size[0] = byte(s >> 24)
	return size
}

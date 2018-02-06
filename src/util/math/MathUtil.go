package math

//int to bytes
func IntToBytes(i int) []byte {
	var buf = make([]byte, 8)

	buf[0] = byte(i & 0xff)
	buf[1] = (byte(i>>8) & 0xff)
	buf[2] = (byte(i>>16) & 0xff)
	buf[3] = (byte(i>>24) & 0xff)
	buf[4] = (byte(i>>32) & 0xff)
	buf[5] = (byte(i>>40) & 0xff)
	buf[6] = (byte(i>>48) & 0xff)
	buf[7] = (byte(i>>56) & 0xff)

	return buf
}

//bytes to int
//len(bytes) must be 8
func BytesToInt(bytes []byte) int {
	i := int(bytes[0] & 0xff)
	i |= (int(bytes[1]&0xff) << 8)
	i |= (int(bytes[2]&0xff) << 16)
	i |= (int(bytes[3]&0xff) << 24)
	i |= (int(bytes[4]&0xff) << 32)
	i |= (int(bytes[5]&0xff) << 40)
	i |= (int(bytes[6]&0xff) << 48)
	i |= (int(bytes[7]&0xff) << 56)

	return i
}

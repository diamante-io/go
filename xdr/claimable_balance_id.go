package xdr

func (e *EncodingBuffer) claimableBalanceCompressEncodeTo(cb ClaimableBalanceId) error {
	if err := e.xdrEncoderBuf.WriteByte(byte(cb.Type)); err != nil {
		return err
	}
	switch cb.Type {
	case ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		_, err := e.xdrEncoderBuf.Write(cb.V0[:])
		return err
	default:
		panic("Unknown type")
	}
}

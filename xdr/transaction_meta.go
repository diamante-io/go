package xdr

// Operations is a helper on TransactionMeta that returns operations
// meta from `TransactionMeta.Operations` or `TransactionMeta.V1.Operations`.
func (transactionMeta *TransactionMeta) OperationsMeta() []OperationMeta {
	switch transactionMeta.V {
	case 0:
		return *transactionMeta.Operations
	case 1:
		return transactionMeta.MustV1().Operations
	case 2:
		return transactionMeta.MustV2().Operations
	default:
		panic("Unsupported TransactionMeta version")
	}
}

package xdr

func (entry *ClaimableBalanceEntry) Flags() ClaimableBalanceFlags {
	switch entry.Ext.V {
	case 1:
		return ClaimableBalanceFlags(entry.Ext.V1.Flags)
	}
	return 0
}

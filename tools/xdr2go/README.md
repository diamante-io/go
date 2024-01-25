# xdr2go

`xdr2go` is a little CLI tool to transform base64 XDR objects into a pretty Go code. This helps in writing mocks and testing. It's using [`fmt.GoStringer`](https://golang.org/pkg/fmt/#GoStringer) interface to print pretty Go code compared to a standard library implementation.

### Why

Very often when writing tests we mock objects to make tests independent of other components. There are many ways we can create example XDR objects but they have disadvantages:
* For transactions we can use `txnbuild` package. The problem is that it allows building transactions only. Also, sometimes there's an existing transaction in the network we want to use in tests so building it using `txnbuild` requires extra time.
* `fmt.Printf("%#v", value)` prints a Go-syntax representation of the value. However, there are many problems with the way the values are printed using a standard `GoStringer`:
   * If the value is a struct with pointers it prints an hexadecimal memory address instead of it's value, ex:
        ```go
        &xdr.TransactionEnvelope{Type: 2, V0: (*xdr.TransactionV0Envelope)(nil), V1: (*xdr.TransactionV1Envelope)(0xc0000aa2d0), FeeBump: (*xdr.FeeBumpTransactionEnvelope)(nil)
        ```
    * It prints redundant `nil` values that could be skipped.
    * Binary values like ed25519 keys and unions (ex. `xdr.Asset`) are rendered in a way that's hard to read.
    * Union discriminants are printed as decimal, instead of using a type variable.
    * `uint32` values are printed in hex.
* I also checked [`https://godoc.org/github.com/kr/pretty`](https://godoc.org/github.com/kr/pretty). While it works fine, it also doesn't skip `nil` values and prints binary values and unions as standard Go printer.

Standard print:
```go
xdr.TransactionEnvelope{Type: 2, V0: (*xdr.TransactionV0Envelope)(nil), V1: (*xdr.TransactionV1Envelope)(0xc0000aa2d0), FeeBump: (*xdr.FeeBumpTransactionEnvelope)(nil)
```
With `GoStringer` implementations:
```go
xdr.TransactionEnvelope{Type: xdr.EnvelopeTypeEnvelopeTypeTx, V1: &xdr.TransactionV1Envelope{Tx: xdr.Transaction{SourceAccount: xdr.MustAddress("GAZ3T7HRWDBJ6SNQ7IWVUS65FP6QMCWHCALFYWX552KFV2O2RLOSRLKI"), Fee: 120, SeqNum: 122783846453215886, TimeBounds: &xdr.TimeBounds{MinTime: xdr.TimePoint(0), MaxTime: xdr.TimePoint(1594645065)}, Memo: xdr.Memo{Type: xdr.MemoTypeMemoNone}, Operations: []xdr.Operation{xdr.OperationBody{SourceAccount: &xdr.MustAddress("GBHC6AMZ3FWLYYHXITCIEZI6VXAU4IEMRCHLICXZXHOVSBFSWCRJ7JS7"), Body: &xdr.OperationBody{Type: xdr.OperationTypeManageSellOffer, ManageSellOfferOp: &xdr.ManageSellOfferOp{Selling: xdr.MustNewCreditAsset("LFEC", "GAG6FS3CR64QJHLHJU7HNXUB4KBLXVDFQBDXM5LG22WOM7CA2ITJAVD2"), Buying: xdr.MustNewNativeAsset(), Amount: 6, Price: xdr.Price{N: 9899999, D: 100000000}, OfferId: 0}}}}, Ext: xdr.TransactionExt{V: 0}}, Signatures: []xdr.DecoratedSignature{xdr.DecoratedSignature{Hint: xdr.SignatureHint{0xda, 0x8a, 0xdd, 0x28}, Signature: xdr.Signature{0x55, 0xb, 0xd0, 0x7d, 0xf3, 0x7, 0x71, 0x56, 0x99, 0x3c, 0x34, 0xfc, 0x47, 0xa0, 0xce, 0x2b, 0x39, 0xa, 0xc4, 0x8c, 0xb7, 0x80, 0x9f, 0x4c, 0xc8, 0x22, 0xae, 0xcc, 0xe9, 0x8b, 0x29, 0xb9, 0x80, 0x94, 0xab, 0x15, 0xbd, 0x6b, 0xc6, 0x3e, 0x2d, 0x12, 0x7a, 0x49, 0xa8, 0x83, 0x75, 0xdd, 0x21, 0x0, 0x14, 0x47, 0xdc, 0xf9, 0x4, 0xe7, 0xdb, 0x16, 0x5a, 0x6e, 0xb0, 0xd6, 0xc0, 0xd}}, xdr.DecoratedSignature{Hint: xdr.SignatureHint{0xb2, 0xb0, 0xa2, 0x9f}, Signature: xdr.Signature{0x6c, 0xb8, 0xe4, 0xd6, 0x39, 0xe2, 0x54, 0x3a, 0x51, 0xf1, 0xc1, 0x20, 0xb9, 0x5f, 0x7d, 0x8, 0xae, 0x31, 0x2c, 0xec, 0x19, 0xce, 0xb0, 0x3f, 0xb6, 0xe6, 0xfa, 0x25, 0xeb, 0x58, 0xf0, 0x33, 0xd1, 0x9f, 0x9a, 0xf5, 0xc2, 0x33, 0x7b, 0xa6, 0x45, 0x92, 0x18, 0xc0, 0x4d, 0xbf, 0x9d, 0xdf, 0xa0, 0xa5, 0x43, 0x73, 0x2a, 0x7c, 0x7d, 0xad, 0x67, 0x54, 0xad, 0xa6, 0xae, 0x30, 0xd1, 0xf}}}}}
```
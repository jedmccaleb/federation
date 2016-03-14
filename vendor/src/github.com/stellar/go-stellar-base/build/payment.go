package build

import (
	"errors"
	
	"github.com/stellar/go-stellar-base/amount"
	"github.com/stellar/go-stellar-base/xdr"
)

// Payment groups the creation of a new PaymentBuilder with a call to Mutate.
func Payment(muts ...interface{}) (result PaymentBuilder) {
	result.Mutate(muts...)
	return
}

// PaymentMutator is a interface that wraps the
// MutatePayment operation.  types may implement this interface to
// specify how they modify an xdr.PaymentOp object
type PaymentMutator interface {
	MutatePayment(*xdr.PaymentOp) error
}

// PaymentBuilder represents a transaction that is being built.
type PaymentBuilder struct {
	O   xdr.Operation
	P   xdr.PaymentOp
	Err error
}

// Mutate applies the provided mutators to this builder's payment or operation.
func (b *PaymentBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case PaymentMutator:
			err = mut.MutatePayment(&b.P)
		case OperationMutator:
			err = mut.MutateOperation(&b.O)
		default:
			err = errors.New("Mutator type not allowed")
		}

		if err != nil {
			b.Err = err
			return
		}
	}
}

// MutatePayment for Asset sets the PaymentOp's Asset field
func (m CreditAmount) MutatePayment(o *xdr.PaymentOp) (err error) {
	o.Amount, err = amount.Parse(m.Amount)
	if err != nil {
		return
	}

	length := len(m.Code)

	var issuer xdr.AccountId
	err = setAccountId(m.Issuer, &issuer)
	if err != nil {
		return
	}

	switch {
	case length >= 1 && length <= 4:
		var code [4]byte
		byteArray := []byte(m.Code)
		copy(code[:], byteArray[0:length])
		asset := xdr.AssetAlphaNum4{code, issuer}
		o.Asset, err = xdr.NewAsset(xdr.AssetTypeAssetTypeCreditAlphanum4, asset)
	case length >= 5 && length <= 12:
		var code [12]byte
		byteArray := []byte(m.Code)
		copy(code[:], byteArray[0:length])
		asset := xdr.AssetAlphaNum12{code, issuer}
		o.Asset, err = xdr.NewAsset(xdr.AssetTypeAssetTypeCreditAlphanum12, asset)
	default:
		err = errors.New("Asset code length is invalid")
	}

	return
}

// MutatePayment for Destination sets the PaymentOp's Destination field
func (m Destination) MutatePayment(o *xdr.PaymentOp) error {
	return setAccountId(m.AddressOrSeed, &o.Destination)
}

// MutatePayment for NativeAmount sets the PaymentOp's currency field to
// native and sets its amount to the provided integer
func (m NativeAmount) MutatePayment(o *xdr.PaymentOp) (err error) {
	o.Asset, err = xdr.NewAsset(xdr.AssetTypeAssetTypeNative, nil)
	if err != nil {
		return
	}

	o.Amount, err = amount.Parse(m.Amount)
	return
}

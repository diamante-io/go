package account

type Account struct {
	Address    string
	Identities []Identity
}

type Identity struct {
	Role        string
	AuthMethods []AuthMethod
}

type AuthMethodType string

func (t AuthMethodType) Valid() bool {
	return AuthMethodTypes[t]
}

const (
	AuthMethodTypeAddress     AuthMethodType = "diamcircle_address"
	AuthMethodTypePhoneNumber AuthMethodType = "phone_number"
	AuthMethodTypeEmail       AuthMethodType = "email"
)

var AuthMethodTypes = map[AuthMethodType]bool{
	AuthMethodTypeAddress:     true,
	AuthMethodTypePhoneNumber: true,
	AuthMethodTypeEmail:       true,
}

type AuthMethod struct {
	Type  AuthMethodType
	Value string
}

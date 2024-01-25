package keypair

import (
	"encoding/hex"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestFull_Hint(t *testing.T) {
	kp := MustParseFull("SBFGFF27Y64ZUGFAIG5AMJGQODZZKV2YQKAVUUN4HNE24XZXD2OEUVUP")
	assert.Equal(t, [4]byte{0x96, 0x47, 0x82, 0x96}, kp.Hint())
	assert.Equal(t, [4]byte{0x96, 0x47, 0x82, 0x96}, kp.FromAddress().Hint())
}

func TestFull_Equal(t *testing.T) {
	// A nil Full.
	var kp0 *Full

	// A Full with a value.
	kp1 := MustParseFull("SBFGFF27Y64ZUGFAIG5AMJGQODZZKV2YQKAVUUN4HNE24XZXD2OEUVUP")

	// Another Full with a value.
	kp2 := MustParseFull("SBPBTSQAIEA5HLWLVWA4TJ7RBKHCEERE2W2DZLB6AUUCEUIYWLJF2EUS")

	// A nil Full should be equal to a nil Full.
	assert.True(t, kp0.Equal(nil))

	// A non-nil Full is not equal to a nil KP with no type.
	assert.False(t, kp1.Equal(nil))

	// A non-nil Full is not equal to a nil Full.
	assert.False(t, kp1.Equal(nil))

	// A non-nil Full is equal to itself.
	assert.True(t, kp1.Equal(kp1))

	// A non-nil Full is equal to another Full containing the same address.
	assert.True(t, kp1.Equal(MustParseFull(kp1.seed)))

	// A non-nil Full is not equal a non-nil Full of different value.
	assert.False(t, kp1.Equal(kp2))
	assert.False(t, kp2.Equal(kp1))
}

var _ = Describe("keypair.Full", func() {
	var subject KP

	JustBeforeEach(func() {
		subject = MustParseFull(seed)
	})

	ItBehavesLikeAKP(&subject)

	type SignCase struct {
		Message   string
		Signature string
	}

	DescribeTable("Sign()",
		func(c SignCase) {
			sig, err := subject.Sign([]byte(c.Message))
			actual := hex.EncodeToString(sig)

			Expect(actual).To(Equal(c.Signature))
			Expect(err).To(BeNil())
		},

		Entry("hello", SignCase{
			"hello",
			"2e75cc20d519111caaaadddf464bb650d2eaf0a5d18d745693a16100f2a4937bc1dffa8b0b1f61a276996d7ee8deb2d0dd9ee510556077b02dec16792e915c0a",
		}),
		Entry("this is a message", SignCase{
			"this is a message",
			"7b7e99d3d660a53913064d5da96abcfa0c422a88f1dca7f14cdbd22045b550030e60fcd1aad85fd08bb7425d95ca690c8f63231895f6b0dd7c0c737227092a00",
		}),
	)

	DescribeTable("SignBase64()",
		func(c SignCase) {
			sig, err := subject.SignBase64([]byte(c.Message))

			Expect(sig).To(Equal(c.Signature))
			Expect(err).To(BeNil())
		},

		Entry("hello", SignCase{
			"hello",
			"LnXMINUZERyqqt3fRku2UNLq8KXRjXRWk6FhAPKkk3vB3/qLCx9honaZbX7o3rLQ3Z7lEFVgd7At7BZ5LpFcCg==",
		}),
		Entry("this is a message", SignCase{
			"this is a message",
			"e36Z09ZgpTkTBk1dqWq8+gxCKojx3KfxTNvSIEW1UAMOYPzRqthf0Iu3Ql2VymkMj2MjGJX2sN18DHNyJwkqAA==",
		}),
	)

	Describe("SignDecorated()", func() {
		It("returns the correct xdr struct", func() {
			sig, err := subject.SignDecorated(message)
			Expect(err).To(BeNil())
			Expect(sig.Hint).To(BeEquivalentTo(hint))
			Expect(sig.Signature).To(BeEquivalentTo(signature))
		})
	})

})

package internal_test

import (
	"encoding/hex"
	"strings"

	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
	"github.com/riposo/default-bucket/internal"
)

var _ = Describe("HashSecret", func() {
	It("should decode", func() {
		var s internal.HashSecret

		Expect((&s).Decode("7965")).To(Succeed())
		Expect(string(s)).To(Equal("ye"))

		Expect((&s).Decode("xx")).NotTo(Succeed())
		Expect((&s).Decode("abc")).NotTo(Succeed())

		Expect((&s).Decode(strings.Repeat("0", 128))).To(Succeed())
		Expect((&s).Decode(strings.Repeat("0", 130))).NotTo(Succeed())
	})

	It("should encode", func() {
		s1 := internal.HashSecret{}
		Expect(hex.EncodeToString(s1.Encode("foo"))).To(Equal("b8fe9f7f6255a6fa08f668ab632a8d081ad87983c77cd274e48ce450f0b349fd"))
		Expect(hex.EncodeToString(s1.Encode("bar"))).To(Equal("844181b39a1b15b417243e6231381b447a3f8b44aa15fbeb845c5d716696e71d"))
		Expect(hex.EncodeToString(s1.Encode("foo"))).To(Equal("b8fe9f7f6255a6fa08f668ab632a8d081ad87983c77cd274e48ce450f0b349fd"))
		Expect(hex.EncodeToString(s1.Encode(""))).To(Equal("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"))

		s2 := internal.HashSecret{'k', 'e', 'y'}
		Expect(hex.EncodeToString(s2.Encode("foo"))).To(Equal("5187fb9829208dafe6afa01a2799846493e76e252cfd05477cb590bee5cc4af2"))
		Expect(hex.EncodeToString(s2.Encode("bar"))).To(Equal("28833aea14891e88b06e1dd45a27cca7924ee3516666a339f098e0ef81bea4dc"))
	})
})

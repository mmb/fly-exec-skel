package flyexecskel_test

import (
	"github.com/mmb/fly-exec-skel/flyexecskel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnvVarName", func() {
	It("uppercases", func() {
		Expect(flyexecskel.EnvVarName("AbCdEf")).To(Equal("ABCDEF"))
	})

	It("converts dashes to underscores", func() {
		Expect(flyexecskel.EnvVarName("a-b-c")).To(Equal("A_B_C"))
	})
})

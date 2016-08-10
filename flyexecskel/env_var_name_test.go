package flyexecskel_test

import (
	"github.com/mmb/fly-exec-skel/flyexecskel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnvVarName", func() {
	Describe("InputEnvVarName", func() {
		It("uppercases and adds input suffix", func() {
			Expect(flyexecskel.InputEnvVarName("AbCdEf")).To(Equal(
				"ABCDEF_INPUT"))
		})

		It("converts dashes to underscores and adds input suffix", func() {
			Expect(flyexecskel.InputEnvVarName("a-b-c")).To(Equal(
				"A_B_C_INPUT"))
		})
	})

	Describe("OutputEnvVarName", func() {
		It("uppercases and adds output suffix", func() {
			Expect(flyexecskel.OutputEnvVarName("AbCdEf")).To(Equal(
				"ABCDEF_OUTPUT"))
		})

		It("converts dashes to underscores and adds output suffix", func() {
			Expect(flyexecskel.OutputEnvVarName("a-b-c")).To(Equal(
				"A_B_C_OUTPUT"))
		})
	})

	Describe("ParamEnvVarName", func() {
		It("uppercases and adds param suffix", func() {
			Expect(flyexecskel.ParamEnvVarName("AbCdEf")).To(Equal(
				"ABCDEF_PARAM"))
		})

		It("converts dashes to underscores and adds param suffix", func() {
			Expect(flyexecskel.ParamEnvVarName("a-b-c")).To(Equal(
				"A_B_C_PARAM"))
		})
	})
})

package flyexecskel_test

import (
	"github.com/mmb/fly-exec-skel/flyexecskel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Divider", func() {
	It("creates a commented out label padded with dashes", func() {
		Expect(flyexecskel.Divider("test label")).To(Equal(
			"# test label -------------------------------------------------------------------"))
	})
})

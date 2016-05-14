package flyexecskel_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFlyexecskel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flyexecskel Suite")
}

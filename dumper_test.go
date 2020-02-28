package dumper

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	// Initialize redis connection for logging
})

var _ = BeforeEach(func() {
	// fmt.Println("Before each")
})

var _ = AfterSuite(func() {
	// fmt.Println("After suite ")
})

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Logger")
}

var _ = Describe("Testing log instantiation", func() {
	It("should do some tests", func() {
		Expect(1).To(Equal(1))
	})
})
package parser_test

import (
	"code.cloudfoundry.org/cli/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command 'Parser'", func() {

	It("returns an unknown command error", func() {
		status := parser.ParseCommandFromArgs([]string{"howdy"})
		Expect(status).To(Equal(-666))
	})
})

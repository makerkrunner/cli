package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"testing"
)

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Parser Suite")
}

var _ = BeforeEach(func() {
	log.SetLevel(log.PanicLevel)
})

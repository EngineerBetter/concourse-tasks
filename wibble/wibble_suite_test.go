package wibble_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWibble(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wibble Suite")
}

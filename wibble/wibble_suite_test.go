package wibble_test

import (
	"github.com/EngineerBetter/wibble"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var specs []*wibble.TaskTestSuite
var yamlFile []byte
var setupErr error

func TestWibble(t *testing.T) {
	RegisterFailHandler(Fail)
	loadSpec("echo_spec.yml")
	loadSpec("file_write_spec.yml")
	loadSpec("input_spec.yml")
	RunSpecs(t, "Wibble Suite")
}

func loadSpec(filename string) {
	yamlFile, setupErr = ioutil.ReadFile(filename)
	expectErrToNotHaveOccurred(setupErr)

	var spec *wibble.TaskTestSuite
	setupErr = yaml.Unmarshal(yamlFile, &spec)
	expectErrToNotHaveOccurred(setupErr)

	specs = append(specs, spec)
}

var _ = BeforeSuite(func() {
	// Do validation that the test spec is valid here, where we can use assertions
	for _, spec := range specs {
		for _, specCase := range spec.Cases {
			for _, input := range specCase.It.HasInputs {
				if input.From != "" {
					Expect(input.From).To(BeADirectory())
				}
			}
		}
	}
})
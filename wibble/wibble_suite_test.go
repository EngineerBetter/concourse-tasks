package wibble_test

import (
	"flag"
	"fmt"
	"github.com/EngineerBetter/wibble"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var specs []*wibble.TaskTestSuite

var specsArg string
var targetArg string
func init() {
	flag.StringVar(&specsArg, "specs", "", "Comma-separated list of spec files to execute")
	flag.StringVar(&targetArg, "target", "", "fly target")
}

func TestWibble(t *testing.T) {
	RegisterFailHandler(Fail)
	specFiles := strings.Split(specsArg, ",")

	if specsArg == "" {
		log.Fatal("--specs must be provided")
	}

	if targetArg == "" {
		log.Fatal("--target must be provided")
	}

	for _, specFile := range specFiles {
		loadSpec(specFile)
	}
	RunSpecs(t, "Concourse Tasks")
}

func loadSpec(filename string) {
	if filename == "" {
		log.Fatalf("Spec file list (%s) contained empty element", specsArg)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("Spec file '%s' does not exist", filename)
	}

	yamlFile, setupErr := ioutil.ReadFile(filename)
	expectErrToNotHaveOccurred(setupErr)

	var spec *wibble.TaskTestSuite
	setupErr = yaml.Unmarshal(yamlFile, &spec)
	expectErrToNotHaveOccurred(setupErr)

	specs = append(specs, spec)
}

var _ = BeforeSuite(func() {
	// Do validation that the test spec is valid here, where we can use assertions
	Expect(specs).ToNot(BeEmpty())
	for _, spec := range specs {
		Expect(spec.Cases).ToNot(BeEmpty(), fmt.Sprintf("%s had no cases", spec.Config))
		for _, specCase := range spec.Cases {
			for _, input := range specCase.It.HasInputs {
				if input.From != "" {
					Expect(input.From).To(BeADirectory())
				}
			}
		}
	}
})
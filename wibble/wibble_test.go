package wibble_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/EngineerBetter/wibble"
	"io/ioutil"
	"os"
	"time"
)

var _ = Describe("", func() {
	for _, specX := range specs {
		// Re-assign as otherwise the single pointer is shared by the anonymous functions, and by the time that they
		// execute, range will have made the single pointer point to something else. Possibly.
		spec := specX

		Describe(spec.Config, func() {
			for _, specCaseX := range spec.Cases {
				specCase := specCaseX

				Describe("when "+specCase.When, func() {
					inputDirs := make(map[string]string)
					for _, input := range specCase.It.HasInputs {
						inputPath, err := ioutil.TempDir("", input.Name)
						expectErrToNotHaveOccurred(err)

						if input.From != "" {
							MustBash("cp -r " + input.From + "/. " + inputPath)
						}

						if input.Setup != "" {
							MustBashIn(inputPath, input.Setup)
						}

						inputDirs[input.Name] = inputPath
					}

					outputDirs := make(map[string]string)
					for _, outputExpectation := range specCase.It.HasOutputs {
						outputPath, err := ioutil.TempDir("", outputExpectation.Name)
						expectErrToNotHaveOccurred(err)
						outputDirs[outputExpectation.Name] = outputPath
					}

					var session *Session
					It(fmt.Sprintf("exits %d", specCase.It.Exits), func() {
						within := specCase.Within
						if within == "" {
							within = defaultTimeout
						}
						timeout, err := time.ParseDuration(within)
						Expect(err).ToNot(HaveOccurred())
						timeout = timeout * time.Duration(timeoutFactorArg)
						session = FlyExecute(targetArg, spec.Config, specCase.Params, inputDirs, outputDirs, timeout)
						Expect(session).To(Exit(specCase.It.Exits), OutErrMessage(session))
						Expect(session).To(Say("executing build"))
						Expect(session).To(Say("initializing"))
					})

					for _, sayExpectationX := range specCase.It.Says {
						sayExpectation := sayExpectationX

						It("says "+sayExpectation, func(){
							Expect(session).To(Say(sayExpectation))
						})
					}

					for _, outputExpectationX := range specCase.It.HasOutputs {
						outputExpectation := outputExpectationX

						Describe(fmt.Sprintf(", it has an output '%s'", outputExpectation.Name), func(){
							for _, forWhichX := range outputExpectation.ForWhich {
								forWhich := forWhichX

								Describe(fmt.Sprintf("for which '%s'", forWhich.Bash), func(){
									var assertionSession *Session
									It(fmt.Sprintf("exits %d", forWhich.Exits), func(){
										// THE REDIRECT IS ABSOLUTE CHEDDAR
										assertionSession = BashIn(outputDirs[outputExpectation.Name], forWhich.Bash+" 2>&1")
										Expect(assertionSession).To(Exit(forWhich.Exits), OutErrMessage(assertionSession))
									})

									for _, sayExpectationX := range forWhich.Says {
										sayExpectation := sayExpectationX
										It("says "+sayExpectation, func(){
											Expect(assertionSession).To(Say(sayExpectation))
										})
									}
								})
							}
						})
					}
				})
			}
		})
	}
})

func expectErrToNotHaveOccurred(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


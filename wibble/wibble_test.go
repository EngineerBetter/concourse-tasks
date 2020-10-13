package wibble_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	. "github.com/EngineerBetter/wibble"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

var _ = Describe("Wibble", func() {
	var spec *TaskTestSuite
	var specFile string

	JustBeforeEach(func() {
		yamlFile, err := ioutil.ReadFile(specFile)
		Expect(err).ToNot(HaveOccurred())

		err = yaml.Unmarshal(yamlFile, &spec)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("echo", func() {
		BeforeEach(func() {
			specFile = "echo_spec.yml"
		})

		It("loads correctly", func() {
			Expect(spec.Config).To(Equal("echo.yml"))
			Expect(len(spec.Cases)).To(Equal(2))
			Expect(spec.Cases[0].When).To(Equal("it is called"))
			Expect(spec.Cases[0].It.Exits).To(Equal(0))
			Expect(spec.Cases[0].It.Says).To(HaveLen(1))
			Expect(spec.Cases[0].Params).To(BeNil())
		})

		It("can be executed", func() {
			executeSpec(spec)
		})
	})

	Describe("file writer", func() {
		BeforeEach(func() {
			specFile = "file_write_spec.yml"
		})

		It("can be executed", func() {
			executeSpec(spec)
		})
	})

	Describe("inputs", func() {
		BeforeEach(func() {
			specFile = "input_spec.yml"
		})

		It("loads correctly", func() {
			Expect(spec.Config).To(Equal("existing_file_write.yml"))
			Expect(spec.Cases[0].It.HasInputs).To(HaveLen(1))
			Expect(spec.Cases[0].It.HasInputs[0].Name).To(Equal("input"))
			Expect(spec.Cases[0].It.HasInputs[0].From).To(Equal("fixtures/existing_file"))

			Expect(spec.Cases[0].It.HasOutputs[0].ForWhich[1].Exits).To(Equal(1))
		})

		It("can be executed", func() {
			executeSpec(spec)
		})
	})
})

func executeSpec(spec *TaskTestSuite) {
	Expect(spec.Config).ToNot(BeNil())
	Expect(spec.Cases).ToNot(HaveLen(0))

	for _, specCase := range spec.Cases {
		Describe(specCase.When, func() {
			inputDirs := make(map[string]string)
			for _, input := range specCase.It.HasInputs {
				inputPath, err := ioutil.TempDir("", input.Name)
				Expect(err).ToNot(HaveOccurred())

				if input.From != "" {
					Expect(input.From).To(BeADirectory())
					mustBash("cp -r "+input.From+"/. "+inputPath)
				}

				if input.Setup != "" {
					mustBashIn(inputPath, input.Setup)
				}

				inputDirs[input.Name] = inputPath
			}

			outputDirs := make(map[string]string)
			for _, outputExpectation := range specCase.It.HasOutputs {
				outputPath, err := ioutil.TempDir("", outputExpectation.Name)
				Expect(err).ToNot(HaveOccurred())
				outputDirs[outputExpectation.Name] = outputPath
			}

			session := flyExecute(spec.Config, specCase.Params, inputDirs, outputDirs)
			Expect(session).To(Exit(specCase.It.Exits), outErrMessage(session))
			Expect(session).To(Say("executing build"))
			Expect(session).To(Say("initializing"))

			for _, sayExpectation := range specCase.It.Says {
				Expect(session).To(Say(sayExpectation))
			}

			for _, outputExpectation := range specCase.It.HasOutputs {
				for _, forWhich := range outputExpectation.ForWhich {
					Expect(forWhich.Bash).ToNot(BeNil())
					// THE REDIRECT IS ABSOLUTE CHEDDAR
					assertionSession := bashIn(outputDirs[outputExpectation.Name], forWhich.Bash+" 2>&1")
					Expect(assertionSession).To(Exit(forWhich.Exits), outErrMessage(assertionSession))
					for _, sayExpectation := range forWhich.Says {
						Expect(assertionSession).To(Say(sayExpectation))
					}
				}
			}
		})
	}
}

func flyExecute(configPath string, params map[string]string, inputDirs, outputDirs map[string]string) *Session {
	pwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--include-ignored", "--input=this=" + pwd}

	for name, dir := range inputDirs {
		flyArgs = append(flyArgs, "--input="+name+"="+dir)
	}

	for name, dir := range outputDirs {
		flyArgs = append(flyArgs, "--output="+name+"="+dir)
	}

	cmd := exec.Command("fly", flyArgs...)
	cmd.Env = os.Environ()
	for key, value := range params {
		setEnv(key, value, cmd)
	}

	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 15*time.Second, time.Second).Should(Exit())
	return session
}

func outErrMessage(session *Session) string {
	return fmt.Sprintf("---\nSTDOUT:\n%v\nSTDERR:\n%v\n---", string(session.Out.Contents()), string(session.Err.Contents()))
}

func setEnv(key, value string, cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, key+"="+value)
}

func bashIn(dir, command string) *Session {
	return bash("cd " + dir + " && " + command)
}

func bash(command string) *Session {
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", command)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 15*time.Second, time.Second).Should(Exit())
	return session
}

func mustBashIn(dir, command string) *Session {
	return mustBash("cd " + dir + " && " + command)
}

func mustBash(command string) *Session {
	session := bash(command)
	Expect(session.ExitCode()).To(BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

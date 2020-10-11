package wibble_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type TaskTestSuite struct {
	Config string `yaml:"config"`
	Cases  []struct {
		When string `yaml:"when"`
		It   struct {
			Exits      int      `yaml:"exits"`
			Says       []string `yaml:"says"`
			HasOutputs []struct {
				Name     string `yaml:"name"`
				ForWhich []struct {
					Bash  string   `yaml:"bash"`
					Exits int      `yaml:"exits"`
					Says  []string `yaml:"says"`
				} `yaml:"for_which"`
			} `yaml:"has_outputs,omitempty"`
		} `yaml:"it,omitempty"`
		Params map[string]string `yaml:"params,omitempty"`
	} `yaml:"cases"`
}

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
})

func executeSpec(spec *TaskTestSuite) {
	Expect(spec.Config).ToNot(BeNil())
	Expect(spec.Cases).ToNot(HaveLen(0))

	for _, specCase := range spec.Cases {
		Describe(specCase.When, func() {
			outputDirs := make(map[string]string)
			for _, outputExpectation := range specCase.It.HasOutputs {
				outputPath, err := ioutil.TempDir("", outputExpectation.Name)
				Expect(err).ToNot(HaveOccurred())
				outputDirs[outputExpectation.Name] = outputPath
			}

			session := flyExecute(spec.Config, specCase.Params, outputDirs)
			Expect(session).To(Exit(specCase.It.Exits), outErrMessage(session))
			Expect(session).To(Say("executing build"))
			Expect(session).To(Say("initializing"))

			for _, sayExpectation := range specCase.It.Says {
				Expect(session).To(Say(sayExpectation))
			}

			for _, outputExpectation := range specCase.It.HasOutputs {
				for _, forWhich := range outputExpectation.ForWhich {
					Expect(forWhich.Bash).ToNot(BeNil())
					assertionSession := bashIn(outputDirs[outputExpectation.Name], forWhich.Bash)
					Expect(assertionSession).To(Exit(forWhich.Exits), outErrMessage(assertionSession))
					for _, sayExpectation := range forWhich.Says {
						Expect(assertionSession).To(Say(sayExpectation))
					}
				}
			}
		})
	}
}

func flyExecute(configPath string, params map[string]string, outputDirs map[string]string) *Session {
	pwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--include-ignored", "--input=this=" + pwd}

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
	Expect(session.ExitCode()).To(BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

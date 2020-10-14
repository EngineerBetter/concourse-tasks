package wibble

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
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
			HasInputs []struct {
				Name  string `yaml:"name"`
				From  string `yaml:"from"`
				Setup string `yaml:"setup"`
			} `yaml:"has_inputs,omitempty"`
		} `yaml:"it,omitempty"`
		Params map[string]string `yaml:"params,omitempty"`
	} `yaml:"cases"`
}

func FlyExecute(target, configPath string, params map[string]string, inputDirs, outputDirs map[string]string) *gexec.Session {
	pwd, err := os.Getwd()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	flyArgs := []string{"-t", target, "execute", "-c", configPath, "--include-ignored", "--input=this=" + pwd}

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

	session, err := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, 20*time.Second, time.Second).Should(gexec.Exit())
	return session
}

func OutErrMessage(session *gexec.Session) string {
	return fmt.Sprintf("---\nSTDOUT:\n%v\nSTDERR:\n%v\n---", string(session.Out.Contents()), string(session.Err.Contents()))
}

func setEnv(key, value string, cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, key+"="+value)
}

func BashIn(dir, command string) *gexec.Session {
	return Bash("cd " + dir + " && " + command)
}

func Bash(command string) *gexec.Session {
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", command)
	session, err := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, 20*time.Second, time.Second).Should(gexec.Exit())
	return session
}

func MustBashIn(dir, command string) *gexec.Session {
	return MustBash("cd " + dir + "; " + command)
}

func MustBash(command string) *gexec.Session {
	session := Bash(command)
	gomega.Expect(session.ExitCode()).To(gomega.BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

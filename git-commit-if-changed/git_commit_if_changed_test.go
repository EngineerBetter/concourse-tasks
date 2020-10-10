package git_commit_if_changed_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var _ = Describe("GitCommitIfChanged", func() {
	var err error
	var stuff map[string]interface{}
	var inputPath, outputPath string

	BeforeEach(func(){
		stuff = make(map[string]interface{})

		// Generate temp dir for each input
		inputPath, err = ioutil.TempDir("", "git-commit-if-changed-input")
		Expect(err).NotTo(HaveOccurred())
		stuff["inputPath"] = inputPath

		bashIn(inputPath, `
			git init
			echo foo > file
			git add -A
			git commit -m "first commit"`)

		// Generate temp dir for each output
		outputPath, err = ioutil.TempDir("", "git-commit-if-changed-output")
		Expect(err).NotTo(HaveOccurred())
		stuff["outputPath"] = outputPath
	})

	When("there is no change", func() {
		It("exits successfully", func(){
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			configPath := filepath.Join(pwd, "git-commit-if-changed.yml")

			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--include-ignored", "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 15*time.Second, time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))
		})
	})

	When("there is a change", func(){
		It("commits the configured change and exits successfully", func(){
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			configPath := filepath.Join(pwd, "git-commit-if-changed.yml")

			bashIn(inputPath,"echo bar >> file")

			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--include-ignored", "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)
			cmd.Env = os.Environ()
			setEnv("GIT_AUTHOR_EMAIL", "test@example.com", cmd)
			setEnv("GIT_AUTHOR_NAME", "Lesley", cmd)
			setEnv("GIT_COMMIT_MESSAGE", "automated commit", cmd)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 15*time.Second, time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))

			status := bashIn(outputPath,"git status")
			Expect(status.Out).ToNot(gbytes.Say("Changes not staged for commit"))
			Expect(status.Out).To(gbytes.Say("nothing to commit, working tree clean"))

			head := bashIn(outputPath,"git show HEAD")
			Expect(head.Out).To(gbytes.Say("Author: Lesley <test@example.com>"))
			Expect(head.Out).To(gbytes.Say("automated commit"))
		})
	})
})

func setEnv(key, value string, cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, key+"="+value)
}

func bashIn(dir, command string) *gexec.Session {
	return bash("cd "+dir+" && "+command)
}

func bash(command string) *gexec.Session {
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", command)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 15*time.Second, time.Second).Should(gexec.Exit())
	Expect(session.ExitCode()).To(BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

func message(flyArgs []string, session *gexec.Session) string {
	return fmt.Sprintf("fly args: %v\nSTDOUT:\n%v\nSTDERR:\n%v", flyArgs, string(session.Out.Contents()), string(session.Err.Contents()))
}

package git_commit_if_changed_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var _ = Describe("GitCommitIfChanged", func() {
	var err error
	var inputPath, outputPath string

	BeforeEach(func(){
		// Generate temp dir for each input
		inputPath, err = ioutil.TempDir("", "git-commit-if-changed-input")
		Expect(err).NotTo(HaveOccurred())

		bashIn(inputPath, `
			git init
			echo foo > file
			git add -A
			git commit -m "first commit"`)

		// Generate temp dir for each output
		outputPath, err = ioutil.TempDir("", "git-commit-if-changed-output")
		Expect(err).NotTo(HaveOccurred())
	})

	When("there is no change", func() {
		It("exits successfully", func(){
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			configPath := filepath.Join(pwd, "git-commit-if-changed.yml")

			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--include-ignored", "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 15*time.Second, time.Second).Should(Exit())
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

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 15*time.Second, time.Second).Should(Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))

			status := bashIn(outputPath,"git status")
			Expect(status.Out).ToNot(Say("Changes not staged for commit"))
			Expect(status.Out).To(Say("nothing to commit, working tree clean"))

			head := bashIn(outputPath,"git show HEAD")
			Expect(head.Out).To(Say("Author: Lesley <test@example.com>"))
			Expect(head.Out).To(Say("automated commit"))
		})
	})
})

func setEnv(key, value string, cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, key+"="+value)
}

func bashIn(dir, command string) *Session {
	return bash("cd "+dir+" && "+command)
}

func bash(command string) *Session {
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", command)
	session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 15*time.Second, time.Second).Should(Exit())
	Expect(session.ExitCode()).To(BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

func message(flyArgs []string, session *Session) string {
	return fmt.Sprintf("fly args: %v\nSTDOUT:\n%v\nSTDERR:\n%v", flyArgs, string(session.Out.Contents()), string(session.Err.Contents()))
}

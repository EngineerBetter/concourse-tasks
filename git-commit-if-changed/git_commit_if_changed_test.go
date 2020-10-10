package git_commit_if_changed_test

import (
	"bytes"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"html/template"
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

		bash(`cd {{.inputPath}}
			git init
			echo foo > file
			git add -A
			git commit -m "first commit"`, stuff)

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

			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 2*time.Minute, time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))
		})
	})

	When("there is a change", func(){
		It("commits the change and exits successfully", func(){
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			configPath := filepath.Join(pwd, "git-commit-if-changed.yml")

			bash(`cd {{.inputPath}}
			echo bar >> file`, stuff)

			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 2*time.Minute, time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))

			after := bash(`cd {{.outputPath}}
			git status`, stuff)
			Expect(after.Out).To(gbytes.Say("wibble"))
		})
	})
})

func bash(command string, stuff map[string]interface{}) *gexec.Session {
	tmpl, err := template.New("command").Parse(command)
	Expect(err).NotTo(HaveOccurred())
	var buff bytes.Buffer
	tmpl.Execute(&buff, stuff)
	parsed := buff.String()
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", parsed)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 2*time.Minute, time.Second).Should(gexec.Exit())
	Expect(session.ExitCode()).To(BeZero(), "bash command: %v\nSTDOUT:\n%v\nSTDERR:\n%v", command, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}

func message(flyArgs []string, session *gexec.Session) string {
	return fmt.Sprintf("fly args: %v\nSTDOUT:\n%v\nSTDERR:\n%v", flyArgs, string(session.Out.Contents()), string(session.Err.Contents()))
}

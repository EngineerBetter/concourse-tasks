package git_commit_if_changed_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var _ = Describe("GitCommitIfChanged", func() {
	When("there is no change", func() {
		It("exits successfully", func(){
			pwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			configPath := filepath.Join(pwd, "git-commit-if-changed.yml")

			// fly login
			// Generate temp dir for each input
			inputPath, err := ioutil.TempDir("", "git-commit-if-changed-input")
			Expect(err).NotTo(HaveOccurred())

			// Generate temp dir for each output
			outputPath, err := ioutil.TempDir("", "git-commit-if-changed-output")
			Expect(err).NotTo(HaveOccurred())

			// fly execute -c configPath --input=NAME=PATH --output=NAME=PATH
			flyArgs := []string{"-t", "eb", "execute", "-c", configPath, "--input=this="+pwd, "--input=input="+inputPath, "--output=output="+outputPath}
			cmd := exec.Command("fly", flyArgs...)

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, 2*time.Minute, time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(BeZero(), message(flyArgs, session))
		})
	})
})

func message(flyArgs []string, session *gexec.Session) string {
	return fmt.Sprintf("fly args: %v\nSTDOUT:\n%v\nSTDERR:\n%v", flyArgs, string(session.Out.Contents()), string(session.Err.Contents()))
}

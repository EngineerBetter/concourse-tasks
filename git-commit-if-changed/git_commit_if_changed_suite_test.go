package git_commit_if_changed_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitCommitIfChanged(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitCommitIfChanged Suite")
}

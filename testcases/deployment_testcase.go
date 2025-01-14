package testcases

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/fixtures"
	. "github.com/cloudfoundry-incubator/bosh-disaster-recovery-acceptance-tests/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DeploymentTestcase struct{}

func (t DeploymentTestcase) Name() string {
	return "deployment_testcase"
}

func (t DeploymentTestcase) BeforeBackup(config Config) {
	By("uploading stemcell", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellSrc,
		)
	})

	By("deploying sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh deploy sdk",
			config,
			"-n",
			"-d",
			"small-deployment",
			"deploy",
			fmt.Sprintf("--var=vm_type=%s", config.BOSH.CloudConfig.DefaultVMType),
			fmt.Sprintf("--var=network_name=%s", config.BOSH.CloudConfig.DefaultNetwork),
			fmt.Sprintf("--var=az_name=%s", config.BOSH.CloudConfig.DefaultAZ),
			fixtures.Path("small-deployment.yml"),
		)
	})
}

func (t DeploymentTestcase) AfterBackup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"delete-deployment",
		)
	})

	By("cleaning up", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"clean-up",
			"--all",
		)
	})
}

func (t DeploymentTestcase) AfterRestore(config Config) {
	By("re-uploading stemcell", func() {
		RunBoshCommandSuccessfullyWithFailureMessage(
			"bosh upload stemcell",
			config,
			"upload-stemcell",
			config.StemcellSrc,
			"--fix",
		)
	})

	By("doing cck to bring back instances", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh cck sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"cck",
			"--auto",
		)
	})

	By("validate deployment instances are back", func() {
		session := RunBoshCommandSuccessfullyWithFailureMessage("bosh get sdk instances",
			config,
			"-n",
			"-d",
			"small-deployment",
			"instances",
		)
		Expect(string(session.Out.Contents())).To(MatchRegexp("small-job/[a-z0-9-]+[ \t]+running"))
	})
}

func (t DeploymentTestcase) Cleanup(config Config) {
	By("deleting sdk deployment ", func() {
		RunBoshCommandSuccessfullyWithFailureMessage("bosh delete sdk deployment",
			config,
			"-n",
			"-d",
			"small-deployment",
			"delete-deployment",
			"--force",
		)
	})
}

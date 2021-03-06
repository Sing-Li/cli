package organization_test

import (
	. "cf/commands/organization"
	"cf/configuration"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestShowOrgRequirements", func() {
		args := []string{"my-org"}

		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}
		callShowOrg(args, reqFactory)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		reqFactory = &testreq.FakeReqFactory{LoginSuccess: false}
		callShowOrg(args, reqFactory)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})
	It("TestShowOrgFailsWithUsage", func() {

		org := models.Organization{}
		org.Name = "my-org"
		org.Guid = "my-org-guid"
		reqFactory := &testreq.FakeReqFactory{Organization: org, LoginSuccess: true}

		args := []string{"my-org"}
		ui := callShowOrg(args, reqFactory)
		Expect(ui.FailedWithUsage).To(BeFalse())

		args = []string{}
		ui = callShowOrg(args, reqFactory)
		Expect(ui.FailedWithUsage).To(BeTrue())
	})
	It("TestRunWhenOrganizationExists", func() {

		developmentSpaceFields := models.SpaceFields{}
		developmentSpaceFields.Name = "development"
		stagingSpaceFields := models.SpaceFields{}
		stagingSpaceFields.Name = "staging"
		domainFields := models.DomainFields{}
		domainFields.Name = "cfapps.io"
		cfAppDomainFields := models.DomainFields{}
		cfAppDomainFields.Name = "cf-app.com"
		org := models.Organization{}
		org.Name = "my-org"
		org.Guid = "my-org-guid"
		org.QuotaDefinition = models.NewQuotaFields("cantina-quota", 512)
		org.Spaces = []models.SpaceFields{developmentSpaceFields, stagingSpaceFields}
		org.Domains = []models.DomainFields{domainFields, cfAppDomainFields}

		reqFactory := &testreq.FakeReqFactory{Organization: org, LoginSuccess: true}

		args := []string{"my-org"}
		ui := callShowOrg(args, reqFactory)

		Expect(reqFactory.OrganizationName).To(Equal("my-org"))

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting info for org", "my-org", "my-user"},
			{"OK"},
			{"my-org"},
			{"  domains:", "cfapps.io", "cf-app.com"},
			{"  quota: ", "cantina-quota", "512M"},
			{"  spaces:", "development", "staging"},
		})
	})
})

func callShowOrg(args []string, reqFactory *testreq.FakeReqFactory) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("org", args)

	token := configuration.TokenInfo{Username: "my-user"}

	spaceFields := models.SpaceFields{}
	spaceFields.Name = "my-space"

	orgFields := models.OrganizationFields{}
	orgFields.Name = "my-org"

	configRepo := testconfig.NewRepositoryWithAccessToken(token)
	configRepo.SetSpaceFields(spaceFields)
	configRepo.SetOrganizationFields(orgFields)

	cmd := NewShowOrg(ui, configRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}

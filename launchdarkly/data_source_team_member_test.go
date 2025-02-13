package launchdarkly

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	ldapi "github.com/launchdarkly/api-client-go"
	"github.com/stretchr/testify/require"
)

func testAccDataSourceTeamMemberConfig(email string) string {
	return fmt.Sprintf(`
data "launchdarkly_team_member" "test" {
  email = "%s"
}
`, email)
}

func testAccDataSourceTeamMemberCreate(client *Client, email string) (*ldapi.Member, error) {
	membersBody := ldapi.MembersBody{
		Email:     email,
		FirstName: "Test",
		LastName:  "Account",
	}
	members, _, err := client.ld.TeamMembersApi.PostMembers(client.ctx, []ldapi.MembersBody{
		membersBody,
	})
	if err != nil {
		return nil, err
	}
	return &members.Items[0], nil
}

func testAccDataSourceTeamMemberDelete(client *Client, id string) error {
	_, err := client.ld.TeamMembersApi.DeleteMember(client.ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func TestAccDataSourceTeamMember_noMatchReturnsError(t *testing.T) {
	email := "does-not-exist@example.com"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTeamMemberConfig(email),
				ExpectError: regexp.MustCompile(`failed to find team member`),
			},
		},
	})
}

func TestAccDataSourceTeamMember_exists(t *testing.T) {
	accTest := os.Getenv("TF_ACC")
	if accTest == "" {
		t.SkipNow()
	}
	randomEmail := fmt.Sprintf("%s@example.com", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	client, err := newClient(os.Getenv(LAUNCHDARKLY_ACCESS_TOKEN), os.Getenv(LAUNCHDARKLY_API_HOST), false)
	require.NoError(t, err)
	member, err := testAccDataSourceTeamMemberCreate(client, randomEmail)
	require.NoError(t, err)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTeamMemberConfig(randomEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.launchdarkly_team_member.test", "email"),
					resource.TestCheckResourceAttr("data.launchdarkly_team_member.test", "email", randomEmail),
					resource.TestCheckResourceAttr("data.launchdarkly_team_member.test", "first_name", member.FirstName),
					resource.TestCheckResourceAttr("data.launchdarkly_team_member.test", "last_name", member.LastName),
					resource.TestCheckResourceAttr("data.launchdarkly_team_member.test", "id", member.Id),
				),
			},
		},
	})
	err = testAccDataSourceTeamMemberDelete(client, member.Id)
	require.NoError(t, err)
}

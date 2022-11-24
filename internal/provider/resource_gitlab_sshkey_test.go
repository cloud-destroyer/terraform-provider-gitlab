//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabSSHKey_basic(t *testing.T) {
	var key gitlab.SSHKey

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabUserSSHKeyDestroy,
		Steps: []resource.TestStep{
			// Create a user + sshkey
			{
				Config: testAccGitlabSSHKeyConfig(testRSAPubKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabSSHKeyExists("gitlab_sshkey.foo_key", &key),
					testAccCheckGitlabSSHKeyAttributes(&key, &testAccGitlabSSHKeyExpectedAttributes{
						Title: "foo-key",
						Key:   testRSAPubKey,
					}),
				),
			},
			// Only update key comment (which is a no-op plan)
			{
				Config:   testAccGitlabSSHKeyConfig(testRSAPubKeyUpdatedComment),
				PlanOnly: true,
			},
			// Update the key and title
			{
				Config: testAccGitlabSSHKeyUpdateConfig(updatedRSAPubKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabSSHKeyExists("gitlab_sshkey.foo_key", &key),
					testAccCheckGitlabSSHKeyAttributes(&key, &testAccGitlabSSHKeyExpectedAttributes{
						Title:     "key",
						Key:       updatedRSAPubKey,
						ExpiresAt: "3016-01-21T00:00:00Z",
					}),
				),
			},
			{
				ResourceName:      "gitlab_sshkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Change pub key to one without a comment
			{
				Config: testAccGitlabSSHKeyConfig(updatedRSAPubKeyWithoutComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabSSHKeyExists("gitlab_sshkey.foo_key", &key),
					testAccCheckGitlabSSHKeyAttributes(&key, &testAccGitlabSSHKeyExpectedAttributes{
						Title: "foo-key",
						Key:   updatedRSAPubKeyWithoutComment,
					}),
				),
			},
			{
				ResourceName:      "gitlab_sshkey.foo_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabSSHKey_ignoreTrailingWhitespaces(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabSSHKeyDestroy,
		Steps: []resource.TestStep{
			// Create a user + sshkey
			{
				Config: fmt.Sprintf(`
					resource "gitlab_sshkey" "this" {
						title   = "test"
						key     = <<EOF
						%s
						EOF
					}
				`, testKeyWithTrailingNewline),
			},
			// Check for no-op plan
			{
				Config: fmt.Sprintf(`
					resource "gitlab_sshkey" "this" {
						title   = "test"
						key     = <<EOF
						%s
						EOF
					}
				`, testKeyWithTrailingNewline),
				PlanOnly: true,
			},
			// Verify Import
			{
				ResourceName:      "gitlab_sshkey.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabSSHKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_sshkey" {
			continue
		}

		keyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse ssh key resource ID: %s", err)
		}

		keys, _, err := testGitlabClient.Users.ListSSHKeys()
		if err != nil {
			return err
		}

		var gotKey *gitlab.SSHKey

		for _, k := range keys {
			if k.ID == keyID {
				gotKey = k
				break
			}
		}
		if gotKey != nil {
			return fmt.Errorf("SSH Key still exists")
		}

		return nil
	}
	return nil
}

func testAccCheckGitlabSSHKeyExists(n string, key *gitlab.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		keyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse ssh key resource ID: %s", err)
		}

		keys, _, err := testGitlabClient.Users.ListSSHKeys()
		if err != nil {
			return err
		}

		var gotKey *gitlab.SSHKey

		for _, k := range keys {
			if k.ID == keyID {
				gotKey = k
				break
			}
		}
		if gotKey == nil {
			return fmt.Errorf("Could not find sshkey %d for currently authenticated user.", keyID)
		}

		*key = *gotKey
		return nil
	}
}

type testAccGitlabSSHKeyExpectedAttributes struct {
	Title     string
	Key       string
	CreatedAt string
	ExpiresAt string
}

func testAccCheckGitlabSSHKeyAttributes(key *gitlab.SSHKey, want *testAccGitlabSSHKeyExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key.Title != want.Title {
			return fmt.Errorf("got title %q; want %q", key.Title, want.Title)
		}

		k := strings.Join(strings.Split(key.Key, " ")[:2], " ")
		wk := strings.Join(strings.Split(want.Key, " ")[:2], " ")

		if k != wk {
			return fmt.Errorf("got key %q; want %q", k, wk)
		}

		return nil
	}
}

func testAccGitlabSSHKeyConfig(pubKey string) string {
	return fmt.Sprintf(`
resource "gitlab_sshkey" "foo_key" {
  title = "foo-key"
  key = "%s"
}
  `, pubKey)
}

func testAccGitlabSSHKeyUpdateConfig(pubKey string) string {
	return fmt.Sprintf(`
resource "gitlab_sshkey" "foo_key" {
  title = "key"
  key = "%s"
  expires_at = "3016-01-21T00:00:00Z"
}
  `, pubKey)
}

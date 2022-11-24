package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_sshkey", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_sshkey` + "`" + ` resource allows to manage the lifecycle of an SSH key assigned to the currently authenticated user.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/users.html#single-ssh-key)`,

		CreateContext: resourceGitlabSSHKeyCreate,
		ReadContext:   resourceGitlabSSHKeyRead,
		DeleteContext: resourceGitlabSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: gitlabSSHKeySchema(),
	}
})

func resourceGitlabSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.AddSSHKeyOptions{
		Title: gitlab.String(d.Get("title").(string)),
		Key:   gitlab.String(d.Get("key").(string)),
	}

	if expiresAt, ok := d.GetOk("expires_at"); ok {
		parsedExpiresAt, err := time.Parse(time.RFC3339, expiresAt.(string))
		if err != nil {
			return diag.Errorf("failed to parse created_at: %s. It must be in valid RFC3339 format.", err)
		}
		gitlabExpiresAt := gitlab.ISOTime(parsedExpiresAt)
		options.ExpiresAt = &gitlabExpiresAt
	}

	key, _, err := client.Users.AddSSHKey(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", key.ID))
	log.Printf("Created key has id %d.", key.ID)
	return resourceGitlabSSHKeyRead(ctx, d, meta)
}

func resourceGitlabSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	keyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("unable to parse user ssh key resource id: %s: %v", d.Id(), err)
	}

	options := &gitlab.ListSSHKeysForUserOptions{
		Page:    1,
		PerPage: 20,
	}

	var key *gitlab.SSHKey
	for options.Page != 0 && key == nil {
		keys, resp, err := client.Users.ListSSHKeys(gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		for _, k := range keys {
			if k.ID == keyID {
				key = k
				break
			}
		}

		options.Page = resp.NextPage
	}

	if key == nil {
		log.Printf("Could not find sshkey %d for currently authenticated user", keyID)
		d.SetId("")
		return nil
	}

	d.Set("key_id", keyID)
	d.Set("title", key.Title)
	d.Set("key", key.Key)
	if key.ExpiresAt != nil {
		d.Set("expires_at", key.ExpiresAt.Format(time.RFC3339))
	}
	if key.CreatedAt != nil {
		d.Set("created_at", key.CreatedAt.Format(time.RFC3339))
	}
	return nil
}

func resourceGitlabSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	keyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("unable to parse ssh key resource id: %s: %v", d.Id(), err)
	}

	if _, err := client.Users.DeleteSSHKey(keyID, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

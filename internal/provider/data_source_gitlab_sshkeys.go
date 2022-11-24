package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_sshkeys", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_sshkeys`" + ` data source allows a list of SSH keys to be retrieved for the currently authenticated user.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/users.html#list-ssh-keys)`,

		ReadContext: dataSourceGitlabKeysRead,
		Schema: map[string]*schema.Schema{
			"keys": {
				Description: "The user's keys.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabSSHKeySchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[INFO] Reading Gitlab user")

	keys, _, err := client.Users.ListSSHKeys(gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("current_user")
	if err := d.Set("keys", flattenSSHKeysForState(keys)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

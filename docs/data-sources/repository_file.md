---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_repository_file Data Source - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_repository_file data source allows details of a file in a repository to be retrieved.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/ee/api/repository_files.html
---

# gitlab_repository_file (Data Source)

The `gitlab_repository_file` data source allows details of a file in a repository to be retrieved.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/repository_files.html)

## Example Usage

```terraform
data "gitlab_repository_file" "example" {
  project   = "example"
  ref       = "main"
  file_path = "README.md"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `file_path` (String) The full path of the file. It must be relative to the root of the project without a leading slash `/`.
- `project` (String) The name or ID of the project.
- `ref` (String) The name of branch, tag or commit.

### Optional

- `id` (String) The ID of this resource.

### Read-Only

- `blob_id` (String) The blob id.
- `commit_id` (String) The commit id.
- `content` (String) base64 encoded file content. No other encoding is currently supported, because of a [GitLab API bug](https://gitlab.com/gitlab-org/gitlab/-/issues/342430).
- `content_sha256` (String) File content sha256 digest.
- `encoding` (String) The file content encoding.
- `file_name` (String) The filename.
- `last_commit_id` (String) The last known commit id.
- `size` (Number) The file size.


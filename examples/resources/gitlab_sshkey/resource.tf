provider "gitlab" {
  alias    = "not_an_admin"
  base_url = "https://gitlab.example.com/api/v4/"
  token    = data.aws_ssm_parameter.gitlab_token_not_an_admin.value
  # token = "gl-some-token-2345652"
}

// Can only manage "own" keys.
// Explicit alias of provider is not required. Merely an example.
resource "gitlab_sshkey" "example" {
  provider   = gitlab.not_an_admin
  title      = "example-key"
  key        = "ssh-rsa AAAA..."
  expires_at = "2016-01-21T00:00:00.000Z"
}


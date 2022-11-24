resource "gitlab_sshkey" "example" {
  title      = "example-key"
  key        = "ssh-rsa AAAA..."
  expires_at = "2016-01-21T00:00:00.000Z"
}


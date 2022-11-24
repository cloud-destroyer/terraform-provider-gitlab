data "gitlab_sshkeys" "current_user" {
}

output "sshkeys" {
  value = data.gitlab_sshkeys.current_user
}

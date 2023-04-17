# resource "spectrocloud_registry_helm" "github-registry" {
#   name       = var.image_registry_name
#   endpoint   = var.image_registry_endpoint
#   is_private = true
#   credentials {
#     credential_type = "basic"
#     username        = var.github_registry_username
#     password        = var.github_registry_password
#   }
# }
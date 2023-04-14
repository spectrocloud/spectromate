# resource "spectrocloud_registry_oci" "image_registry" {
#   name       = var.image_registry_name
#   type       = "basic"
#   endpoint   = "ghcr.io"
#   is_private = true
#   credentials {
#     credential_type = "sts"
#     arn             = "arn:aws:iam::123456:role/stage-demo-ecr"
#     external_id     = "sofiwhgowbrgiornM="
#   }
# }
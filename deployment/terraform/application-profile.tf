# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

resource "spectrocloud_application_profile" "spectromate" {
  name        = "spectromate"
  description = "The Spectromate application profile."
  version     = var.app_version
  pack {
    name            = "redis"
    type            = data.spectrocloud_pack_simple.redis_service.type
    source_app_tier = data.spectrocloud_pack_simple.redis_service.id
    properties = {
      "userPassword"    = base64encode(random_password.password.result)
      "redisVolumeSize" = var.redis_database_volume_size
    }
    values = templatefile("manifests/redis.yaml", {})
  }
  pack {
    name            = "spectromate"
    type            = data.spectrocloud_pack_simple.container_pack.type
    registry_uid    = data.spectrocloud_registry.public_registry.id
    source_app_tier = data.spectrocloud_pack_simple.container_pack.id
    values = templatefile("manifests/container.yaml", {
      image                = var.image
      slack_signing_secret = var.slack_signing_secret
      mendable_api_key     = var.mendable_api_key
      trace_level          = var.trace_level
    })
  }
  pack {
    name            = "spectromate-community"
    type            = data.spectrocloud_pack_simple.container_pack.type
    registry_uid    = data.spectrocloud_registry.public_registry.id
    source_app_tier = data.spectrocloud_pack_simple.container_pack.id
    values = templatefile("manifests/container.yaml", {
      image                = var.image
      slack_signing_secret = var.slack_signing_secret_community
      mendable_api_key     = var.mendable_api_key
      trace_level          = "INFOs"
    })
  }
  tags = var.tags
}
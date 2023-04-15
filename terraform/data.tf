data "spectrocloud_cluster_group" "beehive" {
  name    = var.cluster-group-name
  context = "system"
}

data "spectrocloud_registry" "public_registry" {
  name = "Public Repo"
}

data "spectrocloud_pack_simple" "container_pack" {
  type         = "container"
  name         = "container"
  version      = "1.0.2"
  registry_uid = data.spectrocloud_registry.public_registry.id
}

data "spectrocloud_pack_simple" "redis_service" {
  name         = "redis-operator"
  type         = "operator-instance"
  version      = "6.2.19-1"
  registry_uid = data.spectrocloud_registry.public_registry.id
}
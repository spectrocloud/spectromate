resource "spectrocloud_application" "spectromate" {
  name                    = var.application_name
  application_profile_uid = spectrocloud_application_profile.spectromate.id

  config {
    cluster_name = spectrocloud_virtual_cluster.cluster-1.name
    cluster_uid  = spectrocloud_virtual_cluster.cluster-1.id
  }
  tags = var.tags

  timeouts {
    create = "10m"
    update = "10m"
  }
}
resource "spectrocloud_virtual_cluster" "cluster-1" {
  name              = var.cluster_name
  cluster_group_uid = data.spectrocloud_cluster_group.beehive.id

  resources {
    max_cpu           = 8
    max_mem_in_mb     = 12288
    min_cpu           = 0
    min_mem_in_mb     = 0
    max_storage_in_gb = "12"
    min_storage_in_gb = "0"
  }

  tags = var.tags

  timeouts {
    create = "15m"
    delete = "15m"
  }
}
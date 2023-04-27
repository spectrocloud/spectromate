# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

resource "spectrocloud_virtual_cluster" "cluster-1" {
  name              = var.cluster_name
  cluster_group_uid = data.spectrocloud_cluster_group.beehive.id

  resources {
    max_cpu           = var.cluster_resources["resources"].max_cpu
    max_mem_in_mb     = var.cluster_resources["resources"].max_mem_in_mb
    min_cpu           = var.cluster_resources["resources"].min_cpu
    min_mem_in_mb     = var.cluster_resources["resources"].min_mem_in_mb
    max_storage_in_gb = var.cluster_resources["resources"].max_storage_in_gb
    min_storage_in_gb = var.cluster_resources["resources"].min_storage_in_gb
  }

  tags = var.tags

  timeouts {
    create = "15m"
    delete = "15m"
  }
}
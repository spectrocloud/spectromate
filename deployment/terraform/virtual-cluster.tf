# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

resource "spectrocloud_virtual_cluster" "cluster-1" {
  name              = var.cluster_name
  cluster_group_uid = data.spectrocloud_cluster_group.beehive.id

  resources {
    max_cpu           = 4
    max_mem_in_mb     = 4096
    min_cpu           = 0
    min_mem_in_mb     = 0
    max_storage_in_gb = "4"
    min_storage_in_gb = "0"
  }

  tags = var.tags

  timeouts {
    create = "15m"
    delete = "15m"
  }
}
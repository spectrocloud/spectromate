# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    spectrocloud = {
      version = ">= 0.17.4"
      source  = "spectrocloud/spectrocloud"
    }
    random = {
      version = ">= 3.6.0"
      source  = "hashicorp/random"
    }
  }
}

provider "spectrocloud" {
  project_name = var.project
}

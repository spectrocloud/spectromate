# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    spectrocloud = {
      version = ">= 0.13.0"
      source  = "spectrocloud/spectrocloud"
    }
    random = {
      version = ">= 3.5.0"
      source  = "hashicorp/random"
    }
  }
}

provider "spectrocloud" {
  project_name = var.project
}
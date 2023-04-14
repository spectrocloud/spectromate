terraform {
  required_providers {
    spectrocloud = {
      version = ">= 0.12.0"
      source  = "spectrocloud/spectrocloud"
    }
  }
}

provider "spectrocloud" {
  project_name = var.project
}
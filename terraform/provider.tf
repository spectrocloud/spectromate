terraform {
  required_providers {
    spectrocloud = {
      version = ">= 0.12.0"
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
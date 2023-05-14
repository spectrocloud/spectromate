# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

variable "project" {
  type        = string
  description = "The name of the Spectro Cloud project to use."
  default     = "Default"
}

variable "app_version" {
  type        = string
  description = "The version of the Spectromate application profile in Palette."
  default     = "1.0.0"
}

variable "cluster-group-name" {
  type        = string
  default     = "beehive"
  description = "The name of the cluster group to use for the virtual cluster"
}

variable "slack_signing_secret" {
  type        = string
  description = "The value of the Slack Signing Secret. Set using TF_VAR environment variable."
  sensitive   = true
}

variable "mendable_api_key" {
  type        = string
  description = "The value of the Mendable API Key. Set using TF_VAR environment variable."
  sensitive   = true
}

variable "image" {
  type        = string
  description = "The Spectromate image to deploy."
  default     = "ghcr.io/spectrocloud/spectromate:v1.0.1"
}

variable "cluster_name" {
  type        = string
  description = "The name of the cluster to create."
  default     = "cluster-1"
}

variable "application_name" {
  type        = string
  description = "The name of the application to create."
  default     = "spectromate-app"
}

variable "redis_database_volume_size" {
  type        = string
  description = "The size of the Redis database volume in GiB."
  default     = "8"
}

variable "trace_level" {
  type        = string
  description = "The trace level for the Spectromate application."
  default     = "DEBUG"
}

variable "cluster_resources" {
  description = "The resources to allocate to the virtual cluster"
  type = map(object({
    max_cpu           = number
    max_mem_in_mb     = number
    min_cpu           = number
    min_mem_in_mb     = number
    max_storage_in_gb = string
    min_storage_in_gb = string
  }))
  default = {
    resources = {
      max_cpu           = 4
      max_mem_in_mb     = 4096
      min_cpu           = 0
      min_mem_in_mb     = 0
      max_storage_in_gb = "4"
      min_storage_in_gb = "0"
    }
  }
}


variable "tags" {
  type        = list(string)
  description = "The default tags to apply to Palette resources"
  default     = ["spectro-cloud-education", "app:spectromate", "repository:spectrocloud/spectromate", "terraform_managed:true"]
}
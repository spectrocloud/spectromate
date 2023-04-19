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

variable "image_registry_name" {
  type        = string
  description = "The name of the image registry."
  default     = "github-image-registry-private"
}

variable "image_registry_endpoint" {
  type        = string
  description = "The endpoint of the image registry."
  default     = "ghcr.io"
}


variable "github_registry_username" {
  type        = string
  description = "The username for the image registry. Set using TF_VAR environment variable."
  sensitive   = true
}


variable "github_registry_password" {
  type        = string
  description = "The password for the image registry. Set using TF_VAR environment variable."
  sensitive   = true
}


variable "image" {
  type        = string
  description = "The Spectromate image to deploy."
  default     = "ghcr.io/spectrocloud/spectromate:dev"
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

variable "redis_database_password" {
  type        = string
  description = "The password for the Redis database. Set using TF_VAR environment variable."
}


variable "trace_level" {
  type        = string
  description = "The trace level for the Spectromate application."
  default     = "DEBUG"
}

variable "tags" {
  type        = list(string)
  description = "The default tags to apply to Palette resources"
  default     = ["spectro-cloud-education", "app:spectromate", "repository:spectrocloud/spectromate", "terraform_managed:true"]
}
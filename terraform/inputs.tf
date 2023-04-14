variable "project" {
    type        = string
    description = "The name of the Spectro Cloud project to use."
    default     = "Default"
}

variable "image_registry_name" {
  type        = string
  description = "The name of the image registry."
  default     = "github-image-registry-private"
}


variable "tags" {
  type        = list(string)
  description = "The default tags to apply to Palette resources"
  default     = ["spectro-cloud-education", "app:spectromate", "repository:spectrocloud/spectromate/", "terraform_managed:true"]
}
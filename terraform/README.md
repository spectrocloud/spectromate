## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_random"></a> [random](#requirement\_random) | >= 3.5.0 |
| <a name="requirement_spectrocloud"></a> [spectrocloud](#requirement\_spectrocloud) | >= 0.12.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.5.1 |
| <a name="provider_spectrocloud"></a> [spectrocloud](#provider\_spectrocloud) | 0.12.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [random_password.password](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [spectrocloud_application.spectromate](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/resources/application) | resource |
| [spectrocloud_application_profile.spectromate](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/resources/application_profile) | resource |
| [spectrocloud_registry_helm.github-registry](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/resources/registry_helm) | resource |
| [spectrocloud_cluster_group.beehive](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/data-sources/cluster_group) | data source |
| [spectrocloud_pack_simple.container_pack](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/data-sources/pack_simple) | data source |
| [spectrocloud_pack_simple.redis_service](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/data-sources/pack_simple) | data source |
| [spectrocloud_registry.public_registry](https://registry.terraform.io/providers/spectrocloud/spectrocloud/latest/docs/data-sources/registry) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_application_name"></a> [application\_name](#input\_application\_name) | The name of the application to create. | `string` | `"spectromate-app"` | no |
| <a name="input_cluster-group-name"></a> [cluster-group-name](#input\_cluster-group-name) | The name of the cluster group to use for the virtual cluster | `string` | `"beehive"` | no |
| <a name="input_cluster_name"></a> [cluster\_name](#input\_cluster\_name) | The name of the cluster to create. | `string` | `"cluster-1"` | no |
| <a name="input_github_registry_password"></a> [github\_registry\_password](#input\_github\_registry\_password) | The password for the image registry. | `string` | n/a | yes |
| <a name="input_github_registry_username"></a> [github\_registry\_username](#input\_github\_registry\_username) | The username for the image registry. | `string` | n/a | yes |
| <a name="input_image"></a> [image](#input\_image) | The Spectromate image to deploy. | `string` | `"spectrocloud/spectromate/spectromate:dev"` | no |
| <a name="input_image_registry_endpoint"></a> [image\_registry\_endpoint](#input\_image\_registry\_endpoint) | The endpoint of the image registry. | `string` | `"ghcr.io"` | no |
| <a name="input_image_registry_name"></a> [image\_registry\_name](#input\_image\_registry\_name) | The name of the image registry. | `string` | `"github-image-registry-private"` | no |
| <a name="input_mendable_api_key"></a> [mendable\_api\_key](#input\_mendable\_api\_key) | value of the Mendable API Key | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | The name of the Spectro Cloud project to use. | `string` | `"Default"` | no |
| <a name="input_redis_database_volume_size"></a> [redis\_database\_volume\_size](#input\_redis\_database\_volume\_size) | The size of the Redis database volume in GiB. | `string` | `"8"` | no |
| <a name="input_slack_signing_secret"></a> [slack\_signing\_secret](#input\_slack\_signing\_secret) | value of the Slack Signing Secret | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | The default tags to apply to Palette resources | `list(string)` | <pre>[<br>  "spectro-cloud-education",<br>  "app:spectromate",<br>  "repository:spectrocloud/spectromate/",<br>  "terraform_managed:true"<br>]</pre> | no |
| <a name="input_trace_level"></a> [trace\_level](#input\_trace\_level) | The trace level for the Spectromate application. | `string` | `"DEBUG"` | no |

## Outputs

No outputs.
# loki-cf-logdrain module

Use this module to deploy loki-cf-logdrain in your Cloud foundry space

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_cloudfoundry"></a> [cloudfoundry](#requirement\_cloudfoundry) | >= 0.50.4 |
| <a name="requirement_hsdp"></a> [hsdp](#requirement\_hsdp) | >= 0.42.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_cloudfoundry"></a> [cloudfoundry](#provider\_cloudfoundry) | >= 0.50.4 |
| <a name="provider_hsdp"></a> [hsdp](#provider\_hsdp) | >= 0.42.1 |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [cloudfoundry_app.loki_cf_logdrain](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/resources/app) | resource |
| [cloudfoundry_route.loki_cf_logdrain](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/resources/route) | resource |
| [cloudfoundry_user_provided_service.logdrain](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/resources/user_provided_service) | resource |
| [random_password.token](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
| [cloudfoundry_domain.domain](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/data-sources/domain) | data source |
| [cloudfoundry_org.org](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/data-sources/org) | data source |
| [cloudfoundry_space.space](https://registry.terraform.io/providers/cloudfoundry-community/cloudfoundry/latest/docs/data-sources/space) | data source |
| [hsdp_config.cf](https://registry.terraform.io/providers/philips-software/hsdp/latest/docs/data-sources/config) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cf_org_name"></a> [cf\_org\_name](#input\_cf\_org\_name) | The CF org name to deplo to. | `string` | n/a | yes |
| <a name="input_cf_password"></a> [cf\_password](#input\_cf\_password) | The CF password to use. | `string` | n/a | yes |
| <a name="input_cf_space_name"></a> [cf\_space\_name](#input\_cf\_space\_name) | The CF space name to deploy into. | `string` | n/a | yes |
| <a name="input_cf_user"></a> [cf\_user](#input\_cf\_user) | The CF username to use. | `string` | n/a | yes |
| <a name="input_disk"></a> [disk](#input\_disk) | The amount of Disk space to allocate for Grafana Loki (MB) | `number` | `1024` | no |
| <a name="input_loki_password"></a> [loki\_password](#input\_loki\_password) | The Loki password used for basic auth. | `string` | `""` | no |
| <a name="input_loki_push_endpoint"></a> [loki\_push\_endpoint](#input\_loki\_push\_endpoint) | The Loki push endpoint. This should include /loki/api/v1/push | `string` | n/a | yes |
| <a name="input_loki_username"></a> [loki\_username](#input\_loki\_username) | The Loki username used for basic auth. Default: loki | `string` | `"loki"` | no |
| <a name="input_memory"></a> [memory](#input\_memory) | n/a | `number` | `256` | no |
| <a name="input_name_postfix"></a> [name\_postfix](#input\_name\_postfix) | The name postfix to apply | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | n/a | `string` | n/a | yes |
| <a name="input_tag"></a> [tag](#input\_tag) | n/a | `string` | `"latest"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_logdrain_service_id"></a> [logdrain\_service\_id](#output\_logdrain\_service\_id) | The uuid of the logdrain service. You can bind this to your app to enable logdraining |
| <a name="output_logdrain_url"></a> [logdrain\_url](#output\_logdrain\_url) | Logdrain URL |
<!-- END_TF_DOCS -->

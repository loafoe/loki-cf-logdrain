# loki-cf-logdrain

![loki-cf-logdrain.excalidraw.svg](resources%2Floki-cf-logdrain.excalidraw.svg)

User deployable service which implements a pipeline consisting of a small Go app and a promtail sidecar process. It presents a CF compatible logdrainer endpoint which accepts RFC5424 messages, forwards them to the promtail sidecard process, which in turn forwards the log messages to [Loki](https://grafana.com/oss/loki/) , done.

## Usage

Deployment should be performed using e [bundled Terraform module](modules/loki-cf-logdrain). It injects the proper promtail config for you.

```hcl
module "loki_logdrain" {
  source = "./modules/loki-cf-logdrain"
  name_postfix           = local.postfix
  cf_domain              = var.cf_domain
  cf_space_id            = var.cf_space_id
  
  loki_username = "loki"
  loki_password = "some-secret-password"
  
  loki_push_endpoint = "https://loki.some-fiesta-cluster.terrakube.com/loki/api/v1/push"
}
```


## License

License is MIT

data "cloudfoundry_space" "space" {
  name = var.cf_space_name
  org  = data.cloudfoundry_org.org.id
}

data "cloudfoundry_org" "org" {
  name = var.cf_org_name
}

data "cloudfoundry_domain" "domain" {
  name = data.hsdp_config.cf.domain
}

resource "random_password" "token" {
  length  = 32
  special = false
}

resource "cloudfoundry_app" "loki_cf_logdrain" {
  name         = "tf-loki-logdrain-${var.name_postfix}"
  space        = data.cloudfoundry_space.space.id
  memory       = var.memory
  disk_quota   = var.disk
  docker_image = "loafoe/loki-cf-logdrain:${var.tag}"
  environment = merge({
    TOKEN = random_password.token.result
    PROMTAIL_YAML_BASE64 = base64encode(templatefile("${path.module}/templates/promtail.yaml", {
      loki_push_endpoint = var.loki_push_endpoint
      username           = "loki"
      password           = var.loki_password
    }))
  }, {})
  strategy = "rolling"

  //noinspection HCLUnknownBlockType
  routes {
    route = cloudfoundry_route.loki_cf_logdrain.id
  }
}

resource "cloudfoundry_route" "loki_cf_logdrain" {
  domain   = data.cloudfoundry_domain.domain.id
  space    = data.cloudfoundry_space.space.id
  hostname = "tf-loki-logdrain-${var.name_postfix}"
}

resource "cloudfoundry_user_provided_service" "logdrain" {
  name  = "tf-loki-logdrain-${var.name_postfix}"
  space = data.cloudfoundry_space.space.id
  //noinspection HILUnresolvedReference
  syslog_drain_url = "https://${cloudfoundry_route.loki_cf_logdrain.endpoint}/syslog/drain/${random_password.token.result}"
}

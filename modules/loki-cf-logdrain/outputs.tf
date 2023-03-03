output "logdrain_service_id" {
  description = "The uuid of the logdrain service. You can bind this to your app to enable logdraining"
  value       = cloudfoundry_user_provided_service.logdrain.id
}

output "logdrain_url" {
  description = "Logdrain URL"
  sensitive   = true
  value       = cloudfoundry_user_provided_service.logdrain.syslog_drain_url
}

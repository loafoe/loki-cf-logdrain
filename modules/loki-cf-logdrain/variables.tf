variable "cf_user" {
  description = "The CF username to use."
  type        = string
}

variable "cf_password" {
  description = "The CF password to use."
  type        = string
  sensitive   = true
}

variable "tag" {
  type    = string
  default = "latest"
}

variable "memory" {
  type    = number
  default = 256
}

variable "cf_org_name" {
  description = "The CF org name to deplo to."
  type        = string
}

variable "cf_space_name" {
  description = "The CF space name to deploy into."
  type        = string
}

variable "region" {
  type = string
}

variable "name_postfix" {
  description = "The name postfix to apply"
  type        = string
}

variable "disk" {
  description = "The amount of Disk space to allocate for Grafana Loki (MB)"
  type        = number
  default     = 1024
}

variable "loki_password" {
  description = "The Loki password used for basic auth."
  type        = string
  sensitive   = true
  default     = ""
}

variable "loki_username" {
  description = "The Loki username used for basic auth. Default: loki"
  type        = string
  default     = "loki"
}

variable "loki_push_endpoint" {
  description = "The Loki push endpoint. This should include /loki/api/v1/push"
  type        = string
}

variable "project" {
  type = string
}

variable "notification_channel_id" {
  type = string
}

variable "allowed_list" {
  type    = set(string)
  default = []
}

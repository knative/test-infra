resource "google_monitoring_alert_policy" "boskos_alerts" {
  count        = length(var.allowed_list) == 0 ? 1 : 0
  project      = var.project
  display_name = "boskos-alerts"
  combiner     = "OR" # required

  conditions {
    display_name = "Boskos ran out of resources"

    condition_monitoring_query_language {
      duration = "0s"
      query    = <<-EOT
      fetch prometheus_target
      | metric 'prometheus.googleapis.com/boskos_resources/gauge'
      | {
          t_0:
          filter state == 'free'
          ;
          t_1:
          ident
      }
      | group_by [metric.type]
      | outer_join 0
      | condition t_0.value_boskos_resources_aggregate == 0 && t_1.value_boskos_resources_aggregate > 5
      | window 1m
      EOT
      trigger {
        count = 1
      }
    }
  }

  documentation {
    content   = "Boskos ran out of resources"
    mime_type = "text/markdown"
  }

  # gcloud beta monitoring channels list --project=oss-prow
  #   notification_channels = ["projects/${var.project}/notificationChannels/${var.notification_channel_id}"]
}

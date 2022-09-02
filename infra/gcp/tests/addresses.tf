resource "google_compute_global_address" "prow_v4" {
  project    = module.project.project_id
  name       = "prow-ingress"
  ip_version = "IPV4"
}

resource "google_compute_global_address" "prow_v6" {
  project    = module.project.project_id
  name       = "prow-ingress-v6"
  ip_version = "IPV6"
}

resource "google_compute_global_address" "argo_v4" {
  project    = module.project.project_id
  name       = "argo-ingress"
  ip_version = "IPV4"
}

resource "google_compute_global_address" "grafana_v4" {
  project    = module.project.project_id
  name       = "grafana-ingress"
  ip_version = "IPV4"
}

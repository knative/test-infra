resource "google_certificate_manager_certificate_map" "knative_dev" {
  name        = "knative-dev"
  description = "knative.dev map"
  project     = module.project.project_id
}

resource "google_certificate_manager_dns_authorization" "knative_dev" {
  name        = "knative-dev"
  description = "knative.dev dns"
  domain      = "knative.dev"
  project     = module.project.project_id
}

resource "google_certificate_manager_certificate" "knative_dev" {
  name        = "knative-dev-cert"
  description = "knative.dev wildcard certificate"
  project     = module.project.project_id
  managed {
    domains = [
      "*.knative.dev",
      "knative.dev"
    ]
    dns_authorizations = [
      google_certificate_manager_dns_authorization.knative_dev.id,
    ]
  }
}

resource "google_certificate_manager_certificate_map_entry" "knative_dev" {
  name         = "knative-dev"
  description  = "knative.dev map entry"
  project      = module.project.project_id
  map          = google_certificate_manager_certificate_map.knative_dev.name
  certificates = [google_certificate_manager_certificate.knative_dev.id]
  matcher      = "PRIMARY"
}

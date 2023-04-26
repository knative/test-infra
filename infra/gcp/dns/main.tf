module "project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 12"

  name            = "Knative DNS"
  project_id      = "knative-dns"
  folder_id       = "1055082993535"
  org_id          = "22054930418"
  billing_account = var.billing_account

  # Sane project defaults
  default_service_account     = "keep"
  disable_services_on_destroy = false
  create_project_sa           = false
  random_project_id           = false


  activate_apis = [
    "dns.googleapis.com",
    "compute.googleapis.com"
  ]
}

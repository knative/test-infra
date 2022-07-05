module "project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 12"

  name            = "Knative Shared Infra"
  project_id      = "knative-shared"
  folder_id       = "1055082993535"
  org_id          = "22054930418"
  billing_account = "018CEF-0F96A6-4D1A14"

  # Sane project defaults
  default_service_account     = "keep"
  disable_services_on_destroy = false
  create_project_sa           = false
  random_project_id           = false
  auto_create_network         = true


  activate_apis = [
    "compute.googleapis.com",
    "container.googleapis.com"
  ]
}

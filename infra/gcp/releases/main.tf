module "project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 12"

  name            = "Knative Releases"
  project_id      = "knative-releases"
  folder_id       = "1055082993535"
  org_id          = "22054930418"
  billing_account = "018CEF-0F96A6-4D1A14"

  # Sane project defaults
  default_service_account     = "keep"
  disable_services_on_destroy = false
  create_project_sa           = false
  random_project_id           = false


  activate_apis = [
    "containerregistry.googleapis.com",
  ]
}

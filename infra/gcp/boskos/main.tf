resource "google_folder" "boskos" {
  display_name = "Boskos"
  parent       = "organizations/22054930418"
}

locals {
  boskos_projects = [
    // projects less than 100 were created manually, need to be imported later
    for i in range("101", "110") : format("knative-boskos-%d", i)
  ]
}


module "project" {
  for_each = toset(local.boskos_projects)
  source   = "terraform-google-modules/project-factory/google"
  version  = "~> 12"

  name            = each.key
  project_id      = each.key
  folder_id       = google_folder.boskos.id
  org_id          = "22054930418"
  billing_account = "018CEF-0F96A6-4D1A14"

  # Sane project defaults
  default_service_account     = "keep"
  disable_services_on_destroy = false
  create_project_sa           = false
  random_project_id           = false


  activate_apis = [
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "cloudscheduler.googleapis.com",
    "cloudbuild.googleapis.com",
    "artifactregistry.googleapis.com"
  ]
}

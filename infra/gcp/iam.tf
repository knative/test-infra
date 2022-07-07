module "iam" {
  source        = "terraform-google-modules/iam/google//modules/organizations_iam"
  version       = "~> 7.4"
  organizations = [data.google_organization.org.org_id]
  mode          = "authoritative"

  bindings = {
    "roles/iam.organizationRoleAdmin" = [
      "group:kn-infra-gcp-org-admins@knative.team",
      "serviceAccount:terraform@knative-seed.iam.gserviceaccount.com",
    ]

    "roles/resourcemanager.organizationAdmin" = [
      "group:kn-infra-gcp-org-admins@knative.team",
      "serviceAccount:terraform@knative-seed.iam.gserviceaccount.com",
    ]

    "roles/cloudsupport.techSupportEditor" = [
      "domain:knative.team",
      "group:kn-infra-gcp-org-admins@knative.team"
    ]

    "roles/owner" = [
      "group:kn-infra-gcp-org-admins@knative.team",
      "serviceAccount:terraform@knative-seed.iam.gserviceaccount.com",
    ]

    "roles/resourcemanager.folderAdmin" = [
      "group:kn-infra-gcp-org-admins@knative.team",
      "serviceAccount:terraform@knative-seed.iam.gserviceaccount.com",
    ]

    "roles/browser" = [
      "domain:knative.team",
      "group:gke-security-groups@knative.dev",
    ]

    "roles/resourcemanager.projectCreator" = [
      "group:kn-infra-gcp-org-admins@knative.team",
      "serviceAccount:terraform@knative-seed.iam.gserviceaccount.com",
    ]
  }
}

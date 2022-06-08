module "iam" {
  source  = "terraform-google-modules/iam/google//modules/projects_iam"
  version = "~> 7"

  projects = ["knative-analytics"]

  mode = "authoritative"

  bindings = {

    "roles/bigquery.dataOwner" = [
      "serviceAccount:service-1096014308856@container-engine-robot.iam.gserviceaccount.com",
    ]
  }
}

resource "google_secret_manager_secret" "gsuite" {
  secret_id = "gsuite-group-manager_key"
  project   = module.project.project_id

  labels = {
    environment = "prod"
    system      = "gsuite"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_iam_binding" "gsuite_iam" {
  project   = google_secret_manager_secret.gsuite.project
  secret_id = google_secret_manager_secret.gsuite.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members = [
    "serviceAccount:gsuite-groups-manager@knative-tests.iam.gserviceaccount.com",
  ]
}

// Secret has been loaded manually.
// TODO:create key in TF and pass it to secret manager

resource "google_bigquery_dataset" "billing" {
  dataset_id                  = "billing"
  friendly_name               = "Billing Exports Dataset"
  description                 = "This dataset containts the billing exports."
  location                    = "US"
  project      = module.project.project_id
  labels = {
    environment = "production"
  }
}

resource "google_bigquery_dataset" "gke" {
  dataset_id                  = "gke_usage"
  friendly_name               = "GKE Usage Dataset"
  description                 = "This dataset containts the data of our GKE Usage"
  location                    = "US"
  project      = module.project.project_id
  labels = {
    environment = "production"
  }
}

resource "google_bigquery_dataset" "audits" {
  dataset_id                  = "audit_logging"
  friendly_name               = "Audit Logging Dataset"
  description                 = "This dataset holds critical audit logs for future analysis."
  location                    = "US"
  project      = module.project.project_id
  labels = {
    environment = "production"
  }
}

provider "google" {
}

provider "google-beta" {
}

terraform {
  required_version = "1.2.2"

  backend "gcs" {
    bucket = "knative-state"
    prefix = "boskos"
  }

  required_providers {
    google = {
      version = "4.33.0"
    }
    google-beta = {
      version = "4.33.0"
    }
  }
}

provider "google" {
}

provider "google-beta" {
}

terraform {
  required_version = "0.14.11"

  backend "gcs" {
    bucket = "knative-state"
    prefix = "root"
  }

  required_providers {
    google = {
      version = "4.15"
    }
    google-beta = {
      version = "4.15"
    }
  }
}

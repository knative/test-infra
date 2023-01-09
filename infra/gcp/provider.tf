provider "google" {
}

provider "google-beta" {
}

terraform {
  required_version = "1.3.7"

  backend "gcs" {
    bucket = "knative-state"
    prefix = "root"
  }

  required_providers {
    google = {
      version = "4.47.0"
    }
    google-beta = {
      version = "4.47.0"
    }
  }
}

provider "google" {
}

provider "google-beta" {
}

terraform {
  required_version = "1.1.4"

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

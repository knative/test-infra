# data "google_billing_account" "main" {
#   billing_account = var.billing_account
# }

data "google_organization" "org" {
  domain = var.domain
}

module "dns" {
  source          = "./dns"
  billing_account = var.billing_account
}

module "gsuite" {
  source          = "./gsuite"
  billing_account = var.billing_account

}

module "tests" {
  source          = "./tests"
  billing_account = var.billing_account
}

module "analytics" {
  source          = "./analytics"
  billing_account = var.billing_account

}

module "releases" {
  source          = "./releases"
  billing_account = var.billing_account

}

module "nightly" {
  source          = "./nightly"
  billing_account = var.billing_account

}

module "shared" {
  source          = "./shared"
  billing_account = var.billing_account

}

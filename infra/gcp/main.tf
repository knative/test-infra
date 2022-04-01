data "google_billing_account" "main" {
  billing_account = var.billing_account
}

data "google_organization" "org" {
  domain = var.domain
}


module "boskos" {
  source = "./boskos"
}

module "dns" {
  source = "./dns"
}

module "gsuite" {
  source = "./gsuite"
}

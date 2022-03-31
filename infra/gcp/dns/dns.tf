// knative.dev DNS
module "knative_dev" {
  source     = "terraform-google-modules/cloud-dns/google"
  version    = "3.0.0"
  project_id = module.project.project_id
  type       = "public"
  name       = "knative-dev"
  domain     = "knative.dev."

  recordsets = [
    {
      name = ""
      type = "A"
      ttl  = 300
      records = [
        "75.2.60.5",
      ]
    },
    {
      name = "www"
      type = "CNAME"
      ttl  = 300
      records = [
        "knative.netlify.com.",
      ]
    },
    {
      name = ""
      type = "CAA"
      ttl  = 300
      records = [
        "0 issue \"letsencrypt.org\"",
        "0 issue \"pki.goog\"",
      ]
    },
    {
      name = ""
      type = "TXT"
      ttl  = 300
      records = [
        "\"v=spf1 ?all\"",
        "google-site-verification=w5KR-YluNH94Htu_LcKidfaDfQhlyzRaCp4-_VI5yFY"
      ]
    },
    {
      name = "prow"
      type = "A"
      ttl  = 300
      records = [
        "35.201.93.215",
      ]
    },
    {
      name = "elections"
      type = "CNAME"
      ttl  = 300
      records = [
        "elb.apps.ospo-osci.z3b1.p1.openshiftapps.com.",
      ]
    },
    {
      name = "testgrid"
      type = "CNAME"
      ttl  = 300
      records = [
        "ghs.googlehosted.com.",
      ]
    },
    {
      name = "gubernator"
      type = "CNAME"
      ttl  = 300
      records = [
        "ghs.googlehosted.com.",
      ]
    },
    {
      name = "slack"
      type = "CNAME"
      ttl  = 300
      records = [
        "ghs.googlehosted.com.",
      ]
    },
    {
      name = "blog"
      type = "CNAME"
      ttl  = 300
      records = [
        "ghs.googlehosted.com.",
      ]
    },
    {
      name = "stats"
      type = "CNAME"
      ttl  = 300
      records = [
        "ghs.googlehosted.com.",
      ]
    },
  ]
}

// knative.dev DNS
module "knative_team" {
  source     = "terraform-google-modules/cloud-dns/google"
  version    = "3.0.0"
  project_id = module.project.project_id
  type       = "public"
  name       = "knative-team"
  domain     = "knative.team."

  recordsets = [
    {
      name = ""
      type = "MX"
      ttl  = 300
      records = [
        "1 aspmx.l.google.com.",
        "5 alt1.aspmx.l.google.com.",
        "5 alt2.aspmx.l.google.com.",
        "10 alt3.aspmx.l.google.com.",
        "10 alt4.aspmx.l.google.com.",
      ]
    },
    {
      name = "www"
      type = "CNAME"
      ttl  = 300
      records = [
        "knative.netlify.com.",
      ]
    },
  ]
}

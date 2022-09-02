local util = import 'config_util.libsonnet';

//
// Edit configuration in this object.
//
local config = {
  local comps = util.consts.components,

  // Instance specifics
  instance: {
    name: "Knative Prow",
    botName: "knative-prow-robot",
    url: "https://prow.knative.dev",
  },

  // SLO compliance tracking config
  slo: {
    components: [
      comps.deck,
      comps.hook,
      comps.plank,
      comps.sinker,
      comps.tide,
      comps.monitoring,
    ],
  },

  // Heartbeat jobs
  heartbeatJobs: [
    {name: 'ci-knative-heartbeat', interval: '5m', alertInterval: '20m'},
  ],

  // Tide pools that are important enough to have their own graphs on the dashboard.
  tideDashboardExplicitPools: [],

  // Additional scraping endpoints
  probeTargets: [
  # ATTENTION: Keep this in sync with the list in ../../additional-scrape-configs_secret.yaml
    {url: 'https://prow.knative.dev', labels: {slo: comps.deck}},
  ],

  // Boskos endpoints to be monitored
  boskosResourcetypes: [
    {type: "gke-project", friendly: "GKE projects"},
  ],

  // How long we go during work hours without seeing a webhook before alerting.
  webhookMissingAlertInterval: '60m',

  // How many days prow hasn't been bumped.
  prowImageStaleByDays: {daysStale: 14, eventDuration: '24h'},

  kubernetesExternalSecretServiceAccount: "kubernetes-external-secrets-sa@knative-tests.iam.gserviceaccount.com",
};

// Generate the real config by adding in constant fields and defaulting where needed.
{
  _config+:: util.defaultConfig(config),
  _util+:: util,
}

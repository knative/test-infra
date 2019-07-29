# Automating releases for a new Knative repository

**Note:** Throughout this document, MODULE is a Knative module name, e.g.
`serving` or `eventing`.

By using the release automation already in place, a new Knative repository can
get nightly and official releases with little effort. All automated releases are
monitored through [TestGrid](http://testgrid.knative.dev).

- **Nightly releases** are built every night between 2AM and 3 AM (PST), from
  HEAD on the master branch. They are referenced by a date/commit label, in the
  form `vYYYYMMDD-<commit_short_hash>`. The job status can be checked in the
  `nightly` tab in the corresponding repository dashboard in TestGrid. Images
  are published to the `gcr.io/knative-nightly/MODULE` registry and manifests to
  the `knative-nightly/MODULE` GCS bucket.

- **Versioned releases** are usually built against a branch in the repository.
  They are referenced by a _vX.Y.Z_ label, and published in the _Releases_ page
  of the repository. Images are published to the `gcr.io/knative-release/MODULE`
  registry and manifests to the `knative-releases/MODULE` GCS bucket.

Versioned releases can be one of two kinds:

- **Major or minor releases** are those with changes to the `X` or `Y` values in
  the version. They are cut only when a new release branch (which must be named
  `release-X.Y`) is created from the master branch of a repository. Within about
  2 to 3 hours the new release will be built and published. The job status can
  be checked in the `auto-release` tab in the corresponding repository dashboard
  in TestGrid. The release notes published to GitHub are empty, so you must
  manually edit it and add the relevant markdown content.

- **Patch or dot releases** are those with changes to the `Z` value in the
  version. They are cut automatically, every Tuesday night between 2AM and 3 AM
  (PST). For example, if the latest release on release branch `release-0.2` is
  `v0.2.1`, the next minor release will be named `v0.2.2`. A minor release is
  only created if there are new commits to the latest release branch of a
  repository. The job status can be checked in the `dot-release` tab in the
  corresponding repository dashboard in TestGrid. The release notes published to
  GitHub are a copy of the previous release notes, so you must manually edit it
  and adjust its content.

## Setting up automated releases

1. Have the
   [//test/presubmit-tests.sh](prow_setup.md#setting-up-jobs-for-a-new-repo)
   script added to your repo, as it's used as a release gateway. Alternatively,
   have some sort of validation and set `$VALIDATION_TESTS` in your release
   script (see below).

1. Write your release script, which will publish your artifacts. For details,
   see the
   [helper script documentation](../scripts/README.md#using-the-releasesh-helper-script).

1. Enable `nightly`, `auto-release` and `dot-release` jobs for your repo in the
   [config_knative.yaml](prow/config_knative.yaml) file. For example:

   ```
   knative/MODULE:
    - nightly: true
    - dot-release: true
    - auto-release: true
   ```

1. Run `make config` to regenerate [config.yaml](prow/config.yaml), otherwise
   the presubmit test will fail. Merge such pull request and ask one of the
   owners of _knative/test-infra_ to:

   1. Run `make update-config` in `ci/prow`.

   1. Run `make update-config` in `ci/testgrid`.

   Within two hours the 3 new jobs (nightly, auto-release and dot-release) will
   appear on TestGrid.

   The jobs can also be found in the
   [Prow status page](https://prow.knative.dev) under the names
   `ci-knative-MODULE-nightly-release`, `ci-knative-MODULE-auto-release` and
   `ci-knative-MODULE-dot-release`.

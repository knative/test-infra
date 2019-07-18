# Automating releases for a new Knative repository

By using the release automation already in place, a new Knative repository can
get nightly and official releases with little effort. All automated releases
are monitored through [TestGrid](http://testgrid.knative.dev).

* **Nightly releases** are built every night from HEAD on the master branch.
  They are referenced by a date/commit label, in the form `vYYYYMMDD-<commit_short_hash>`.
* **Versioned releases** are usually built against a branch in the repository.
  They are referenced by a *vX.Y.Z* label, and published in the *Releases* page
  of the repository.

## Setting up automated releases

1. Have the [//test/presubmit-tests.sh](prow_setup.md#setting-up-jobs-for-a-new-repo)
   script added to your repo, as it's used as a release gateway. Alternatively,
   have some sort of validation and set `$VALIDATION_TESTS` in your release script
   (see below).

1. Write your release script, which will publish your artifacts. For details, see
   the [helper script documentation](../scripts/README.md#using-the-releasesh-helper-script).

1. Enable `nightly`, `auto-release` and `dot-release` jobs for your repo in the
   [config_knative.yaml](prow/config_knative.yaml) file. For example:

   ```
   knative/MODULE:
    - nightly: true
    - dot-release: true
    - auto-release: true
   ```

1. Run `make config` to regenerate [config.yaml](prow/config.yaml), otherwise the presubmit
   test will fail. Merge such pull request and ask one of the owners of *knative/test-infra*
   to:

   1. Run `make update-config` in `ci/prow`.

   1. Run `make update-config` in `ci/testgrid`.

   Within two hours the new 3 jobs will appear on TestGrid.

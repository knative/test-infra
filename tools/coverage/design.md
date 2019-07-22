# Design of Test Coverage Tool

![design.svg](design.svg)
We pack the test coverage tool in a container, that is triggered by prow. It runs test coverage profiling on target repository. Afterward, it calculates coverage change or summarizes coverage data, depends on workflow type, explained below.  

## Post-submit workflow

Produces and stores coverage profile, as master branch, for later presubmit jobs to compare against

1. A PR is merged.
1. Post-submit prow job started.
1. Test coverage profile generated. Completion marker generated upon successful run.

## Pre-submit workflow

The pre-submit workflow is triggered when a pull request is created or updated.
It runs code coverage tool to report coverage change through a comment by a robot github account.
The robot account will be referred to as "metrics-bot" in later text.

1. Developer submits a new commit to an open PR on GitHub.
1. Matching pre-submit prow job is started.
1. Test coverage profile generated.
1. Calculate coverage change against master branch. Compare the coverage file generated in this cycle against the most recent successful master branch build.
1. Use PR data from GitHub, git-attributes, as well as coverage change data calculated above, to produce a list of files that we care about in the line-by-line coverage report. produce line by line coverage html and add link to metrics-bot report.
1. Let metrics-bot post presubmit coverage on GitHub, under that conversation of the PR.

## Periodical workflow

Produces periodic coverage result as input for TestGrid

1. Periodical prow job starts periodically.
The frequency and start time can be configured in [the config file](../../ci/prow/config.yaml)
1. Test coverage profile and metadata generated.
1. Generate and store per-file coverage data.
  - Stores in the XML format that is used by TestGrid.
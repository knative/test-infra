#Design of Test Coverage Tool
See design.svg for the design diagram.

We pack the test coverage tool in a container, that is triggered by prow. It runs test coverage profiling on target repository. Prow clones the target repository as the current working directory for the container.

 Prow has handlers for different github events. We add a pre-submit prow job that is triggered by any new commit to a PR to run test coverage on the new build to compare it with the master branch and previous commit for pre-submit coverage. We add a post-submit prow job that is triggered by merge events, to run test coverage on the nodes of master branch. Test coverage data on the master branch is supplied to TestGrid for displaying the coverage change over time, as well as serve as the basis of comparison for pre-submit coverage mentioned in the pre-submit scenario.

 Here is the step-by-step description of the pre-submit and post-submit workflows

 ##Pre-submit workflow
Runs code coverage tool to report coverage change in a new commit or new PR
1. Developer submit new commit to an open PR on github
2. Matching pre-submit prow job is started 
3. Test coverage profile generated
4. Calculate coverage change against master branch. Compare the coverage file generated in this cycle against the most recent successful post-submit build. Coverage file for post-submit commits were generated in post-submit workflow
5. Use PR data from github, git-attributes, as well as coverage change data calculated above, to produce a list of files that we care about in the line-by-line coverage report. produce line by line coverage html and add link to metrics-bot report.
6. Let metrics-bot post presubmit coverage on github, under that conversation of the PR. 

 ##Post-submit workflow
Produces & stores coverage profile for later presubmit jobs to compare against
1. A PR is merged
2. Post-submit prow job started
3. Test coverage profile generated. Completion marker generated upon successful run.

 ##Periodical workflow
Produces periodical coverage result as input for TestGrid
1. Periodical prow job starts periodically based on the specification in prow job config
2. Test coverage profile & metadata generated 
3. Generate / store per-file coverage data
  - Stores in the XML format that is used by TestGrid
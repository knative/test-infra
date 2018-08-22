# GitHub Helper Tool

This tool is designed to interact with GitHub, providing useful data for a Prow job.

Currently the tool makes unauthenticated requests to GitHub API.

## Flags

* `-list-changed-files` will list the files that are touched by the current PR in a Prow job.
* `-verbose` will dump extra info on output when executing the comments; it is intended for debugging

#!/bin/bash
if [ -z "$1" ]
then
  echo "$0: missing first parameter: lower case project name"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 1
fi

if [ -z "$2" ]
then
  echo "$0: missing second parameter: 'Project Name'"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 2
fi

if [ -z "$3" ]
then
  echo "$0: missing third parameter: 'project_org/main_repo'"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 3
fi

if [ -z "$4" ]
then
  echo "$0: missing fourth parameter: start-date in YYYY-MM-DD format (>= 2015-01-01)"
  echo "Usage: $0 lower_case_project_name 'Project Name' 'project_org/main_repo' start-date"
  exit 4
fi

./set_project_name.sh homebrew Homebrew 'Homebrew/brew' '2015-01-01' "$1" "$2" "$3" "$4"

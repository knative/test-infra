import subprocess
import typing
import argparse

ROLES = [
    ("owner", ["prime-engprod-sea@google.com"]),
    ("editor", [
        "knative-productivity-admins@googlegroups.com",
        "knative-tests@appspot.gserviceaccount.com",
        "prow-job@knative-tests.iam.gserviceaccount.com",
        "prow-job@knative-nightly.iam.gserviceaccount.com",
        "prow-job@knative-releases.iam.gserviceaccount.com",
    ])
]

def addRole(project, role, accounts:typing.List[str]):
  for account in accounts:
    accountType = None
    if account.endswith("@googlegroups.com") or account.endswith("prime-engprod-sea@google.com"):
      accountType = "group"
    elif account.endswith(".gserviceaccount.com"):
      accountType = "serviceAccount"
    else:
      print("account '{}' matched no known account type", account)
      raise NameError

    print(f"[{project}] Adding new {role}: {account}")
    completed = subprocess.run(['gcloud',
                    'projects', 'add-iam-policy-binding', project,
                    '--member', f'{accountType}:{account}',
                    '--role', f'roles/{role}'], capture_output=True)
    if completed.returncode == 1:
      print(completed)
      completed.check_returncode()

def addRoles(project, roles):
  for role, accounts in roles:
    addRole(project, role, accounts)

def addRolesToProjects(start, end):
  projects = [f"knative-boskos-{number}" for number in range(start, end)]
  for project in projects:
    addRoles(project, ROLES)

if __name__ == "__main__":
  parser = argparse.ArgumentParser(description="Add owner and editor roles")
  parser.add_argument("start", type=int,
                      help="which project number to start from.")
  parser.add_argument("end", type=int, help="which project number to end with.")

  args = parser.parse_args()
  addRolesToProjects(args.start, args.end + 1)

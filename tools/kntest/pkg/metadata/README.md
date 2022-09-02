## kntest metadata

`kntest metadata` command is used for operations on the metadata.json file in
Prow job artifacts.

## Subcommands

### Set

`kntest metadata set` subcommand can be invoked with following parameters:

- `--key`: meta info key
- `--value`: meta info value

Note the value can be overwritten if the key already exists in the metadata.json
file.

### Get

`kntest metadata get` subcommand can be invoked with following parameters:

- `--key`: meta info key

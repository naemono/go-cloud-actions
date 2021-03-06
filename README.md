# Public Cloud Common Commands/Actions

This repo contains some common commands that are run in public clouds (currently azure, gcp, and aws).

```bash
$ ./bin/cloud -h
A cli to interact with users, permissions, and peering in public cloud.

Usage:
  cloud [flags]
  cloud [command]

Available Commands:
  compute     Control compute in public clouds
  help        Help about any command
  identity    Control identity (users and permissions) in public clouds
  network     Control networks in public clouds
  peering     Control peering of VPCs/VNets in public clouds
  resources   Control resources in public clouds

Flags:
  -h, --help              help for cloud
  -l, --loglevel string   logging level (default "info")
  -v, --version           version for cloud

Use "cloud [command] --help" for more information about a command.
```

| Command       | SubCommands                   | Description    |
| -----------   | -----------                   | ----------      |
| compute       | create-container-instance, create-cluster     | Create Container Instances, Create GKE cluster |
| identity      | applications [add, add-credentials], roles [list], users  [add]  | Add Appications/Users |
| network       | network-profile  [add, list], vpc [create, create-subnet, delete, list, list-subnets], regions [az-list]  | Add/List Network Profiles, CRUD operations on AWS VPCs, Availability zone listing |
| peering       | [create, list]                | Add/List Network Peerings |
| resources     | resource-groups [add]         | Add Resource Groups |
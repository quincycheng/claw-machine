# claw-machine
CLI for migrating local users on Privilege Cloud to Identity Security Platform

![go run main.go](./assets/screencap.png.png)

## Setup
1. Clone this repo
2. Copy `config.json.sample` to `conjur.json`
2. Update configuration of `conjur.json`
3. Execute `go run main.go`

## Conversion 

### Entities
Due to platform difference, below table summarizes the conversion of various entities

| From (Privilege Cloud) | To (Identity Security Platform) |
|------------------------|---------------------------------|
| Group                  | Role                            |
| User                   | User                            |
| User > Role            | Role                            |

### Default Group to Role Mapping

| From (Privilege Cloud) | To (Identity Security Platform) |
|------------------------|---------------------------------|
| Auditors               | Privilege Cloud Auditors        |
| Operators              | Privilege Cloud Safe Managers   |
| Vault Admins           | Privilege Cloud Administrators  |

### Default User Role to Role Mapping

| From (Privilege Cloud) | To (Identity Security Platform) |
|------------------------|---------------------------------|
| Admin                  | Privilege Cloud Administrators  |
| Auditor                | Privilege Cloud Auditors        |
| Safe manager           | Privilege Cloud Safe Managers   |
| Standard user          | Privilege Cloud Users           |

## Default Groups to be skipped

Below are the built-in group from Privilege Cloud that can be skipped during migration
- Backup Users
- DR Users
- Notification Engines
- PVWAAppUsers
- PVWAGWAccounts
- PVWAMonitor
- PVWAUsers
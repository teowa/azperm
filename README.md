## azperm - get Azure resource permissions from terraform config or resource types

### Usage

```
> azperm
NAME:
   azperm - Retrieves Azure resource permissions from Terraform configurations (Azure/azapi provider only) or specified resource types. Requires Azure CLI login.

USAGE:
   azperm [--file-name <file_name>] [--resource-type <resource_type>]

COMMANDS:
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --file-name , [ --file-name , ]          One or more Terraform configuration files to parse, separated by commas(,).
   --resource-type , [ --resource-type , ]  One or more resource types to parse, separated by commas(,).
   --help, -h                               show help
```

### Examples

```
> azperm --file-name=main.tf
{
  "Microsoft.KeyVault/vaults/secrets@2023-02-01": {
    "action": [
      "Microsoft.KeyVault/vaults/secrets/read",
      "Microsoft.KeyVault/vaults/secrets/write"
    ],
    "dataActions": [
      "Microsoft.KeyVault/vaults/secrets/delete",
      "Microsoft.KeyVault/vaults/secrets/backup/action",
      "Microsoft.KeyVault/vaults/secrets/purge/action",
      "Microsoft.KeyVault/vaults/secrets/update/action",
      "Microsoft.KeyVault/vaults/secrets/recover/action",
      "Microsoft.KeyVault/vaults/secrets/restore/action",
      "Microsoft.KeyVault/vaults/secrets/readMetadata/action",
      "Microsoft.KeyVault/vaults/secrets/getSecret/action",
      "Microsoft.KeyVault/vaults/secrets/setSecret/action"
    ]
  }
}

```

```

> azperm --resource-type=Microsoft.KeyVault/vaults/secrets@2023-02-01
{
  "Microsoft.KeyVault/vaults/secrets@2023-02-01": {
    "action": [
      "Microsoft.KeyVault/vaults/secrets/read",
      "Microsoft.KeyVault/vaults/secrets/write"
    ],
    "dataActions": [
      "Microsoft.KeyVault/vaults/secrets/delete",
      "Microsoft.KeyVault/vaults/secrets/backup/action",
      "Microsoft.KeyVault/vaults/secrets/purge/action",
      "Microsoft.KeyVault/vaults/secrets/update/action",
      "Microsoft.KeyVault/vaults/secrets/recover/action",
      "Microsoft.KeyVault/vaults/secrets/restore/action",
      "Microsoft.KeyVault/vaults/secrets/readMetadata/action",
      "Microsoft.KeyVault/vaults/secrets/getSecret/action",
      "Microsoft.KeyVault/vaults/secrets/setSecret/action"
    ]
  }
}

```

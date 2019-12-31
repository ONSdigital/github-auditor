# GitHub Auditor
This repository contains a Go application that consumes the GraphQL GitHub audit log API.

## Running
### Environment Variables
The environment variables below are required:

```
GITHUB_ORG_NAME # Name of the GitHub Enterprise organisation
GITHUB_TOKEN    # GitHub personal access token
```

### Token Scopes
The GitHub personal access token for using this application requires the following scopes:

- `admin:enterprise`
- `admin:org`
- `repo`
- `user`

## Copyright
Copyright (C) 2019 Crown Copyright (Office for National Statistics)
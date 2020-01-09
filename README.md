# GitHub Auditor
This repository contains a Go application that consumes the GraphQL GitHub audit log API.

## Running
### Environment Variables
The environment variables below are required:

```
FIRESTORE_CREDENTIALS # Path to the GCP service account JSON key
FIRESTORE_PROJECT     # Name of the GCP project containing the Firestore project
GITHUB_ORG_NAME       # Name of the GitHub Enterprise organisation
GITHUB_TOKEN          # GitHub personal access token
SLACK_ALERTS_CHANNEL  # Name of the Slack channel to post alerts to
SLACK_WEBHOOK         # Used for accessing the Slack Incoming Webhooks API
```

### Token Scopes
The GitHub personal access token for using this application requires the following scopes:

- `admin:enterprise`
- `admin:org`
- `repo`
- `user`

## Copyright
Copyright (C) 2020 Crown Copyright (Office for National Statistics)
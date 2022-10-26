# GitHub Auditor
This repository contains a [Go](https://golang.org/) application that consumes the [GitHub audit log API](https://developer.github.com/v4/interface/auditentry/) and posts Slack alerts for events of interest. A [Cloud Firestore](https://cloud.google.com/firestore/) database is used to store a small amount of state to ensure duplicate alerts aren't created.

## Building
Use `make` to compile binaries for macOS and Linux.

## Running
### Environment Variables
The environment variables below are required:

```
FIRESTORE_PROJECT     # Name of the GCP project containing the Firestore project
GITHUB_ORG_NAME       # Name of the GitHub Enterprise organisation
GITHUB_TOKEN          # GitHub personal access token
SLACK_ALERTS_CHANNEL  # Name of the Slack channel to post alerts to
SLACK_WEBHOOK         # Used for accessing the Slack Incoming Webhooks API
```

The environment variable below is optional:

```
FIRESTORE_CREDENTIALS # Path to the GCP service account JSON key (used when running locally)
```

### Token Scopes
The GitHub personal access token for using this application requires the following scopes:

- `admin:org`
- `repo`
- `user`

## Copyright
Copyright (C) 2020 Crown Copyright (Office for National Statistics)
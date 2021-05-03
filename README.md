# synthetic

[![Go Report Card](https://goreportcard.com/badge/github.com/ifosch/synthetic)](https://goreportcard.com/report/github.com/ifosch/synthetic)
[![Maintainability](https://api.codeclimate.com/v1/badges/010f4a83de5d0d6bcad0/maintainability)](https://codeclimate.com/github/ifosch/synthetic/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/010f4a83de5d0d6bcad0/test_coverage)](https://codeclimate.com/github/ifosch/synthetic/test_coverage)

## Setup

### Requirements

You will need:
- A bot account on your team's workspace on Slack, with the token
  stored in the `SLACK_TOKEN` environment variable.
- A Jenkins user. The Jenkins URL will be stored in the `JENKINS_URL`
  environment variable; the username, in the `JENKINS_USER` one; and,
  the password in the `JENKINS_PASSWORD` one.

### Using the docker image

You can use the [Synthetic Docker
image](https://hub.docker.com/repository/docker/natx/synthetic/tags?page=1&ordering=last_updated).

The
`latest` tag of this image is built on every update of the master
branch, so it's considered development and, thus, unstable.

For every tag in this repo, there is an automatically generated Docker
image. Like for `v0.0.1`, there is a corresponding tag on this image,
with the same name.

## Roadmap

Things to come are:
- Start providing some k8s management commands.
- Fix Jenkins UI when Folders are present.
- Validate arguments when using Jenkins `build` command.
- Poll Jenkins for running jobs updates.
- Improve test coverage.
- Add instrumentation and monitoring with Prometheus.
- Improve logging.
- Improve message processing techniques.
- Handle more Slack events.
- Improve artifacts provided to the user.
- Implement a plugin system.


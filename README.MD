# K8S Report

[![Go Report Card](https://goreportcard.com/badge/github.com/jsenon/k8sslackreport)](https://goreportcard.com/report/github.com/jsenon/k8sslackreport)
[![CircleCI](https://circleci.com/gh/jsenon/k8sslackreport.svg?style=svg)](https://circleci.com/gh/jsenon/k8sslackreport)

## Func

Report in slack:

- Number of pods with details (running, failed)
- Number of nodes with details (running, failed)
- Number of namespaces
- Number of jobs
- Number of persistent volume claims (bound, not bound)

## Build

```sh
make build
```

## Usage

Right now pod watcher is available on default and kube-system namespaces

Create Webhook slack application and export URL

```sh
export SLACK_URL=https://hooks.slack.com/xxx
```

For internal k8s:

```sh
k8sslackreport report --api internal
```

For External k8s, its based on current context in your .kube/config file:

```sh
k8sslackreport report --api external
```

### Env

- `SLACK_URL` = Slack URL Hook

## TODO

- [ ] Reduce complexity for report.go
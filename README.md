<h1 align="center">devents 💫  </h1>

<h5 align="center">Collects Docker Events across Hosts</h5>

<br/>


`devents` aims at providing visibility into `docker events` across multiple hosts. It allows you to answer the questions like:

> "At a given point in time, what was happening in this daemon? Did it create a container? Did it attach the network to the container?"
> "What is the rate of container creation that we're seeing in the set of docker daemon labelled with `com.testing=test-machines`?"
> "What images are being pulled the most?"

###  Table of Contents
 
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Usage](#usage)
  - [Container](#container)
- [Aggregators](#aggregators)
  - [Stdout](#stdout)
  - [Fluentd](#fluentd)
- [Metrics](#metrics)
  - [Label Retrieval](#label-retrieval)
    - [label](#label)
- [LICENSE](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


### Usage

```
Usage: devents [--fluentdhost FLUENTDHOST] [--fluentdtag FLUENTDTAG] [--fluentdport FLUENTDPORT] [--dockerhost DOCKERHOST] [--aggregator AGGREGATOR] [--metricspath METRICSPATH] [--metricsport METRICSPORT] [--metricslabel METRICSLABEL]

Options:
  --fluentdhost FLUENTDHOST
                         fluentd host to connect to [default: localhost]
  --fluentdtag FLUENTDTAG
                         fluentd tag to add to the messages [default: devents]
  --fluentdport FLUENTDPORT
                         fluentd port to connect to [default: 24224]
  --dockerhost DOCKERHOST
                         docker daemon to connect to [default: unix://var/run/docker.sock]
  --aggregator AGGREGATOR, -a AGGREGATOR
                         aggregators to use (stdout|fluentd|prometheus) [default: []]
  --metricspath METRICSPATH
                         path to use for prometheus scrapping [default: /metrics]
  --metricsport METRICSPORT
                         port to listen for prometheus scrapping [default: 9103]
  --metricslabel METRICSLABEL
                         includes labels from containers|images in the timeseries
  --help, -h             display this help and exit
```


#### Container

- `labels`: specify a list of labels to extract volumes from.


### Aggregators

#### Stdout

Events are simply flushed to `stdout`:

```
devents \
        --aggregator stdout
```


#### Fluentd

You can tie `fluentd` with `devents` agents by specifying fluentd configuration:

```
devents \
        --aggregator fluentd \
        --fluentd-port=24224 \
        --fluentd-host=localhost \
        --fluentd-tag=com.mytag
```


### Metrics

`devents` is also able to perform metric collection. It does so via `prometheus`. By default it exposes the metrics endpoint at `/metrics` on port `9090`. These are configurable:

```
devents \
        --aggregator prometheus \
        --metrics-path /prometheus-metrics \
        --metrics-port 1337
```

#### Label Retrieval

Some event types support the extraction of extra parameters (attributes).

##### label

> Supported by: `image` and `container`

A great use of this is in metrics of user-defined namespacing. For instance:

- containers are created with labels `com.mypaas.project=<project-id`
- `devents` is initiated with `--metrics-label com.mypaas.project`
- query for the instant rate of `project-specific` container creation: `irate(devents_container_start{com-mypaas-project="prjectId"}[5m])`

Note.: prometheus labels can't have `.`, so, `.`s in the labels are replaced by `-`. For instance:


1. Start `devents` with capturing of `com.docker.swarm.service.id` label:

```
devents \
        --aggregator prometheus \
        --metricslabel com.docker.swarm.service.id
```

2. Run a container with a label: 

```sh
docker run \
        --label com.docker.swarm.service.id=123 \
        -d \
        nginx:alpine`
```

3. Check the /metrics endpoint:

```sh
# HELP devents_container_action Docker container actions performed
# TYPE devents_container_action counter
devents_container_action{action="create",com_docker_swarm_service_id=""} 1
devents_container_action{action="create",com_docker_swarm_service_id="123"} 1
devents_container_action{action="start",com_docker_swarm_service_id=""} 1
devents_container_action{action="start",com_docker_swarm_service_id="123"} 1
```


### LICENSE

MIT


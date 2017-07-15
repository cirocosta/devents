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


- [Collection](#collection)
- [Aggregators](#aggregators)
  - [Stdout](#stdout)
  - [Fluentd](#fluentd)
- [Metrics](#metrics)
- [LICENSE](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### Collection

Once events are received `devents` is capable of enhancing the event message with more information. For instance, `{type=container,action=start}` can be enhanced with environment variables, ports and node information.


### Aggregators


#### Stdout

Events are simply flushed to `stdout`:

```
devents \
        --stdout
```


#### Fluentd

You can tie `fluentd` with `devents` agents by specifying fluentd configuration:

```
devents \
        --fluentd-port=24224 \
        --fluentd-host=localhost \
        --fluentd-tag=com.mytag
```


### Metrics

`devents` is also able to perform metric collection. It does so via `prometheus`. By default it exposes the metrics endpoint at `/metrics` on port `9090`. These are configurable:

```
devents \
        --metrics-path /prometheus-metrics \
        --metrics-port 1337
```


### LICENSE

MIT


<h1 align="center">devents 💫  </h1>

<h5 align="center">Collects Docker Events across Hosts</h5>

<br/>


`devents` aims at providing visibility into `docker events` across multiple hosts. It allows you to answer the questions like:

- "At a given point in time, what was happening in this daemon? Did it create a container? Did it attach the network to the container?"
- "What is the rate of container creation that we're seeing in the set of docker daemon labelled with `com.testing=test-machines`?"
- "What images are being pulled the most?"


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


### LICENSE

MIT


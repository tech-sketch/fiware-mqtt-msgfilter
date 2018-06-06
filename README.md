# fiware-mqtt-msgfilter
This REST API service works with [coreos/etcd](https://coreos.com/etcd/docs/latest/) in order to check the message duplication.

[![TravisCI Status](https://travis-ci.org/tech-sketch/fiware-mqtt-msgfilter.svg?branch=master)](https://travis-ci.org/tech-sketch/fiware-mqtt-msgfilter)
[![DockerHub Status](https://dockerbuildbadges.quelltext.eu/status.svg?organization=techsketch&repository=fiware-mqtt-msgfilter)](https://hub.docker.com/r/techsketch/fiware-mqtt-msgfilter/builds/)

## Description
This REST API service accepts the **POST** request to `/distinct/`.

The duplication check flow is like below:

1. This service returns `200 OK` when the received message has not been stored to etcd cluster yet, or has been expired already from etcd cluster.
1. Otherwise, this service returns `409 Conflict`

## Environment Variables
This REST API accept Environment Variables like below:

|Environment Variable|Summary|Default|
|:--|:--|:--|
|`LISTEN_PORT`|listen port of this service|5001|
|`ETCD_ENDPOINT`|endpoint url of etcd cluster|http://127.0.0.1:2379|
|`LOCK_TTL`|expire second(s) for lock key|10|
|`DATA_TTL`|expore second(s) for data|600|

## Request Payload
`Content-Type: application/json`

```json
{
  "payload": "message to check duplication"
}
```

## API specification

see [docs/swagger.yaml](/docs/swagger.yaml)


## Run as Docker container

1. Pull container [techsketch/fiware-mqtt-msgfilter](https://hub.docker.com/r/techsketch/fiware-mqtt-msgfilter/) from DockerHub.

    ```bash
    $ docker pull techsketch/fiware-mqtt-msgfilter
    ```
1. Run Container.
    * Set environment variable(s) if you want to change exposed port, etcd endpoint, and so on.

    ```bash
    $ env ETCD_ENDPOINT="http://192.168.0.3:2379" LISTEN_PORT="3000" docker run -d -p 3000:3000 techsketch/fiware-mqtt-msgfilter
    ```

## Build from source code

1. go get

    ```bash
    $ go get -u github.com/tech-sketch/fiware-mqtt-msgfilter
    $ cd ${GOPATH}/src/github.com/tech-sketch/fiware-mqtt-msgfilter
    ```
1. install dependencies

    ```bash
    $ go get -u github.com/golang/dep/cmd/dep
    $ dep ensure
    ```
1. go install

    ```bash
    $ go install github.com/tech-sketch/fiware-mqtt-msgfilter
    ```
1. run service

    ```bash
    $ env ETCD_ENDPOINT="http://192.168.0.3:2379" LISTEN_PORT="3000" ${GOPATH}/bin/fiware-mqtt-msgfilter
    ```

## License

[Apache License 2.0](/LICENSE)

## Copyright
Copyright (c) 2018 TIS Inc.

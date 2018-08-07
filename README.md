# CNI Network Controller [![Build Status](https://travis-ci.org/linkernetworks/network-controller.svg?branch=master)](https://travis-ci.org/linkernetworks/network-controller) [![codecov](https://codecov.io/gh/linkernetworks/network-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/linkernetworks/network-controller) [![Go Report Card](https://goreportcard.com/badge/github.com/linkernetworks/network-controller)](https://goreportcard.com/report/github.com/linkernetworks/network-controller)  [![Docker Build Status](https://img.shields.io/docker/build/sdnvortex/network-controller.svg)](https://hub.docker.com/r/sdnvortex/network-controller/)

![overview](./images/overview.png)

Open vSwitch, Multiple network interfaces that associate with kubernetes pods, etc.

## Table of Contents
* [Development](#development)
* [Usage](#usage)
  + [Run a Server](#run-a-server)
  + [Run a Client](#run-a-client)
    - [Bridge Name](#bridge-name)
    - [Interface Name](#interface-name)
    - [IP Address of the interface](#ip-address-of-the-interface)
    - [Add VLAN tag into the interface](#add-vlan-tag-into-the-interface)
    - [Add a IP routing table(Add route)](#add-a-ip-routing-tableadd-route)
    - [Server](#server)
    - [Example](#example)

## Development

```shell
# generate protocol buffer
make pb

# make server binary
make server

# make client binary
make client

# make test (You should run this before push codes)
make test
```

## Usage

### Run as a server
The network-controller server provide two ways to listen, TCP and Unix domain socket  
If you want to run as a UNIX domain socket server, you should specify a path to store the sock file  
and server will remove that file when server is been terminated  
```shell
./server/network-controller-server -unix=/tmp/a.sock
```
For the TCP server, just use the `IP:PORT` format, for example, `0.0.0.0:50051`
```shell
./server/network-controller-server -tcp=0.0.0.0:50051
```

### Run as a client
The clinet is used for the kubernetes pod to create the veth and you can see the example yaml in `example/kubernetes/*.yaml` to see how to use this client.  

For creating a veth for Pod, the client needs the following information
- Pod Name
- Pod Namespace
- Pod UUID
- Target Bridge Name
- The Interface Name in the container
- IP Address of the interface in the container
- Add VLAN tag into the interface
- Add a IP routing table(Add route)
- The server address

The first three variable can passed by the environemtn `POD_NAME`, `POD_NAMESPACE` and `POD_UUID`.

#### Bridge Name
`-b` or `--bridge` Target bridge name

#### Interface Name
`-n` or `--nic` The interface name in the container

#### IP Address of the interface
`-i` or `--ip` The ip address of the interface, should be a valid v4 CIDR Address

#### Add VLAN tag into the interface
`-v` or `--vlan` The Vlan Tag of the interface

#### Add a IP routing table(Add route)
`--net` The destination network for add IP routing table, like '-net target'  
`-g` or `--gateway` The gateway of the interface subnet

#### Servers
The clinet support two way to connect to the server, TCP socket and UNIX domain socket.  
In the TCP mode, use the `IP:PORT` format to connect to TCP server  
```shell
./client/network-controller-client -server=0.0.0.0:50051
```
Fot the UNIX domain socket mode, you should use the `unix://PATH` format to connect to server.  
Assume the path is `/tmp/a.sock` and you can use the following command to connect  
```shell
./client/network-controller-client -server=unix:///tmp/a.sock
```

#### Example
Following is an example of the client and you can see more use the `--help`.
```shell
./clinet/network-controller-client -s unix:///tmp/vortex/vortex.sock -b br100 -n eth100 -i 192.168.2.2/24 --net 239.0.0.0/4 -g 0.0.0.0
```

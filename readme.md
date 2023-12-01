# ps util tools

yao plugin for getting system infomation.

refer to :

- 1Panel: https://github.com/1Panel-dev/1Panel

- gopsutil: https://github.com/shirou/gopsutil

## build

```sh

cd plugins/pstuil

make build

```

## test

plugin test

```sh
# cpu
yao run plugins.psutil.cpu

# disk
yao run plugins.psutil.disk

# docker
yao run plugins.psutil.docker

# host
yao run plugins.psutil.host

# load
yao run plugins.psutil.load

# mem
yao run plugins.psutil.mem

# mem2
yao run plugins.psutil.mem2

# net every thing related to net
yao run plugins.psutil.net

# net2 only the tcp and udp connection
yao run plugins.psutil.net2

# process
yao run plugins.psutil.process

# ssh sessions
yao run plugins.psutil.ssh_session

# windows serivces only for windows system
yao run plugins.psutil.winservices

```

## LICENSE

MIT

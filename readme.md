# ps util tools

yao plugin for getting system infomation.

the default plugin folder path is `<YAO_EXTENSION_ROOT>/plugins/`, the default value for YAO_EXTENSION_ROOT is the app folder, you can change the YAO_EXTENSION_ROOT in the .env file.

refer to :

- 1Panel: https://github.com/1Panel-dev/1Panel

- gopsutil: https://github.com/shirou/gopsutil

## build

```sh

make build

mv psutil.so to yaoapp/plugins

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

# dashboard
yao run plugins.psutil.dashboard

# mem
yao run plugins.psutil.mem

# mem2 all the numbers converted to string format dispaly
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

# sxtln command (cli)

sextillion cli is called `sxtln`
The cli keeps its configuration under ${HOME}/.sxtln/config. The default config values are empty except for the output which is set to "json".
The config file is a yaml file

## Content

* [Installation](#installation)
* [Commands](#Commands)
* [Cluster](#Cluster)
* [auth](#auth)
* [Config](#Config)

## installation

Look for the [latest version at releases](https://github.com/sxtln/cli/releases/latest). 
Each release has 3 files, 
- sxtln for linux
- sxtln sxtln-darwin for macOS
- sxtln.exe for windows 64 bit

Download the file and copy it to your prefered folder. Say your using linux and your prefered folder is /usr/local/bin. 
Than you need to run the following:
``` bash
sudo cp sxtln /usr/local/bin/
sudo chmod +x /usr/local/bin/sxtln
```

That's it. you are ready to go.

## General format

The format of the command is `sxtln [global flags] <command> <subcommand> <command flags> [parameters]` 

The output can be in json or yaml format. In this document only the json format is shown as output

You can set the output of the cli by using `-o` or `--output`. possible values are "json" or "yaml"

## Commands

A list of all the command

### Cluster

Command to deploy view and destroy a cluster. Follow are sub commands of cluster

#### ls

list your allocated clusters. If you add the name of a cluster this will show only that cluster info. If the cluster does not exists the command return none zero  

##### example

```bash

sxtln cluster ls 

or 

sxtln cluster ls [cluster-name]

```

output

list of clusters:

```json

{"clusters":[{"clusterId":"67e01083","createDate":"2024-04-23T14:39:30.831Z","k8sConfig":"apiVersion: v1\nclusters:\n  - cluster:\n      certificate-authority-data: ...","name":"cli_testing_cluster","nodes":[{"clusterId":"67e01083","createDate":"2024-04-23T14:39:30.831Z","dataCenterId":"100", "name":"n67e010830","nodeForigenId":46416213,"nodeId":"2ce8a56e93","nodeIndex":0,"nodeType":"b1","privateIp":"10.0.0.2","provisionAt":"2024-04-23T14:39:33.106Z"}],"provisionAt":"2024-04-23T14:39:33.106Z","stage":"k8sReady","stageIndex":9,"upTime":9}],"userId":"acc2fe1d3b57d8650b95"}


```

single cluster:

```json

{"clusterId":"67e01083","createDate":"2024-04-23T14:39:30.831Z","k8sConfig":"apiVersion: v1\nclusters:\n  - cluster:\n      certificate-authority-data: ... ","name":"cli_testing_cluster","nodes":[{"clusterId":"67e01083","createDate":"2024-04-23T14:39:30.831Z","dataCenterId":"100","name":"n67e010830","nodeForigenId":46416213,"nodeId":"2ce8a56e93","nodeIndex":0,"nodeType":"b1","privateIp":"10.0.0.2","provisionAt":"2024-04-23T14:39:33.106Z"}],"provisionAt":"2024-04-23T14:39:33.106Z","stage":"k8sReady","stageIndex":9,"upTime":5}

```


#### create

create a new cluster. 

possible flags are:
- __name__, __n__ name of cluster. if not set sextillion will give a default name
- __type__ , __t__ optionaly, type of nodes in the cluster. default is b1
- __count__, __c__ optional number of nodes in the cluster. default 1
- __force__, __f__ optional force create even if a cluster with the same name exists
- __wait__, __w__ optional wait until the cluster is created or destroyed

##### example

```bash
sxtln cluster create [--name <optionaly, name of cluster>] [--type <optionaly, type of node of the cluster>]  [--count <number of nodes> default 1] 

```

#### delete

delete a cluster. 

possible flags are:

- __id__, __i__ id of the cluster to delete


##### example

```bash

sxtln cluster delete --id <id of cluster to delete>

or

sxtln cluster delete -i <id of cluster to delete>

```

#### config, kube-config, kubeconfig, kc

get a cluster kubeconfig file. Output is in plain text

possible flags are:

- __id__, __i__ id of the cluster to delete


##### example

```bash

sxtln cluster kc -i <id of cluster to delete>

or

sxtln cluster config --id <id of cluster to delete>

or

sxtln cluster kubeconfig -i <id of cluster to delete>

or

sxtln cluster kube-config -i <id of cluster to delete>

```

> Note: Do not confuse with `sxtln config` command


### auth

Handle authentication with sextillion api server

#### check

Checks if apiKey or token are valid

Returns `{"ok":true/false}`

if failed to test with server exit code will not be zero

##### example

command

```bash
sxtln auth check

```

output
```bash

{ok:true}

```


#### login

login with user password to server. If succeeded the token is saved as part of the config.token 

Returns `{"token": "token value"}`

if failed to test with server exit code will not be zero

##### example

command

```bash
sxtln auth login --user=<user> --password=<pwd>

or 

sxtln auth login -u=<user> -p=<pwd>

```

output

```bash
{"token":"some token"}

```

#### set-apikey

set api key on local machine. You generate your api key via the web interface under user settings

Returns sxtln config as json or yaml

##### example

command

```bash
sxtln auth set-apikey --key="api key"

or 

sxtln auth set-apikey -k "api key"

```

output

sxtln config as json



#### logout

logout by clearing the token field from the configuration

Returns sxtln config as json or yaml

##### example

command

```bash
sxtln auth logout

```

output

sxtln config as json or yaml


### Config

set, get or view sxtln configuration. Follow are sub commands of config

#### get

Get one of the configuration fields

configuration fields are:

- __token__: The token received after calling login
- __apiKey__: The apiKey
- __output__: The prefered output format of sxtln. Default is json. Can be json or yaml

##### example

get the apiKey of configuration

```bash

sxtln config get --key token

or 

sxtln config get -k token

```

output Json with the field value

```json

{"token": "token value"}

```

#### set

Set one of the configuration fields. Can only be token, apiKey or output

##### example

get the apiKey of configuration

```bash

sxtln config set --key apikey "newApiKey"

or 

sxtln config get -k apikey "newApiKey"

```

output all config fields as json or yaml


#### view

Set one of the configuration fields. Can only be token, apiKey or output

##### example

get the apiKey of configuration

```bash

sxtln config view

```

output all config fields as json or yaml


# ssh-auto
ssh-auto is a tool to manager ssh login info.
# install
```
go get github.com/yushuailiu/ssh-auto
```

# help

```

$ssh-auto --help
NAME:
   ssh-auto - A new cli application

USAGE:
   ssh-auto [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     add, a     add a server info
     delete, d  delete a serve by id or name
     edit, e    edit a serve by id or name
     list, l    list all servers
     login      login a server by id
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
Usage of subcommand of list
```
$ssh-auto list --help

NAME:
   ssh-auto list - list all servers

USAGE:
   ssh-auto list [command options] [arguments...]

OPTIONS:
   --detail, -d              show detail of server
   --filter value, -f value  filter servers

```
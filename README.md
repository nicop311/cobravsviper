# `cobravsviper`: make the most out of Golang Viper and Cobra

TL;DR: `viper.Sub()` does not maintain the right priority order (CLI > Env Vars > Config File > Default)
for flags of a cobra subcommand'.

This is an issue when unmarshaling `viper.Unmarshal`: each cobra subcommands will ignore configuration files or
cobra CLI or env var depending on how, when and in what order `viper.Unmarshal` is called.

## How the project was bootstraped

Init the go workspace:

```bash
go mod init cobravsviper
```

Then use [`cobra-cli`](https://github.com/spf13/cobra-cli) to bootstrap the root 
CLI command and add a `version` subcommand.

```bash
cobra-cli init
cobra-cli add version
```

## Build The Project Locally

```bash
go build -o cobravsviper  main.go
```

## CLI User Inputs Priority

In this project, [Cobra](https://github.com/spf13/cobra) is used to handel the CLI commands, subcommands
and all their flags.

[Viper](https://github.com/spf13/viper) is used with its features `viper.BindPFlags` and `viper.AutomaticEnv` which
allow Viper to automatically use Cobra's flags as both environment variables and configuration file.

This allow the developpers to only add and modify cobra flags, and viper will automatically adapt without the need for a
dedicated viper configuration.

### The root cobra command help message

```
./cobravsviper -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  cobravsviper [flags]
  cobravsviper [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     A SUBcommand

Flags:
      --config string      Configuration File (default "configs/myconfig.yaml")
  -h, --help               help for cobravsviper
      --rootflag1 string   root flag 1 (default "from default")
      --rootflag2 string   root flag 2 (default "from default")
      --rootflag3 string   root flag 3 (default "from default")
      --rootflag4 string   root flag 4 (default "from default")
  -t, --toggle             Help message for toggle

Use "cobravsviper [command] --help" for more information about a command.
```

### The `version` cobra subcommand help message

`version` is the name of our first subcommand.

```
./cobravsviper version -h
A Cobra Subcommand

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  cobravsviper version [flags]

Flags:
  -h, --help                  help for version
      --versionflag1 string   version flag 1 (default "fromdefault")
      --versionflag2 string   version flag 2 (default "fromdefault")
      --versionflag3 string   version flag 3 (default "fromdefault")
      --versionflag4 string   version flag 4 (default "fromdefault")

Global Flags:
      --config string      Configuration File (default "configs/myconfig.yaml")
      --rootflag1 string   root flag 1 (default "from default")
      --rootflag2 string   root flag 2 (default "from default")
      --rootflag3 string   root flag 3 (default "from default")
      --rootflag4 string   root flag 4 (default "from default")
```


### Priority In Theory: CLI > Env Vars > Config File > Default

| Priority Order & Source  | Example                                                                     | Comment                         |
|--------------------------|-----------------------------------------------------------------------------|---------------------------------|
| 1️⃣ CLI Flag             | `--rootflag1 debug`                                                         | Highest priority                |
| 2️⃣ Environment Variable | `COBRAVSVIPER_ROOTFLAG2=from envvars`                                       | Overrides config file & default |
| 3️⃣ Config File          | `rootflag3: "from config file"` in YAML/TOML/JSON                                         | Overrides default               |
| 4️⃣ Default Value        | `Flags().StringVar(&rootFlag4, "rootflag4", "from default", "root flag 4")` | Used if nothing else is set     |

In this example:
* we choose the first root flag `--rootflag1` as the one to illustrate CLI flag user input, but this is arbitrary ; we could choose any of the 4 flags `--rootflagX`.
* we choose the second root flag `--rootflag2` (aka `COBRAVSVIPER_ROOTFLAG2`) as the one to illustrate CLI environment variable user input, but this is arbitrary ; we could choose any of the 4 env vars `COBRAVSVIPER_ROOTFLAGX`.
* we choose the third root flag `--rootflag3` (aka `rootflag3`) as the one to illustrate the configuration variable; but this is arbitrary.
* we choose the third root flag `--rootflag3` to test the defaul value.

This priority works for the root cobra command. But this priority does not work for the version command.

### In Practice: 

#### Cobra Root Command: CLI Priority Success

In this situation, `viper` and `cobra` behave as expected: the priority "CLI > Env Vars > Config File > Default" is repected.

```bash
COBRAVSVIPER_ROOTFLAG2="from envvars" ./cobravsviper --rootflag1="from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] Root command called                           cobra-cmd=cobravsviper
INFO[0000] rootflag1: from cli                           cobra-cmd=cobravsviper
INFO[0000] rootflag2: from envvars                       cobra-cmd=cobravsviper
INFO[0000] rootflag3: from configuration file            cobra-cmd=cobravsviper
INFO[0000] rootflag4: from default                       cobra-cmd=cobravsviper
```

#### Cobra Sub Command: CLI Priority Failure

In this situation, we try the cobra subcommand `version` and set only the flags corresponding to this subcommand.
Notice `viper` does not handle well the subcommand's flags.

Ignore the output of the _Persistent flags from rootCmd_ for now: this output is perfectly normal, 
since the priority would be to use the values from the config file for `rootflag{1, 2, 3}` because
no other input was specified for these parameters. `rootflag4` is commented in the configuration,
and no other input was specified for this parameter, so it takes the default value.

```bash
COBRAVSVIPER_VERSIONFLAG2="from envvars" ./cobravsviper version --versionflag1="from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] version called                                cobra-cmd=version
INFO[0000] versionflag1: from cli                        cobra-cmd=version
INFO[0000] versionflag2: from envvars                    cobra-cmd=version
INFO[0000] versionflag3: from default                    cobra-cmd=version
INFO[0000] versionflag4: from default                    cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: from cli                           cobra-cmd=version
INFO[0000] rootflag2: from envvars                       cobra-cmd=version
INFO[0000] rootflag3: from configuration file            cobra-cmd=version
INFO[0000] rootflag4: from default                       cobra-cmd=version
```

Notice `versionflag3` has the value `from default` which comes from viper using cobra's default
values instead of `from configuration file`. If viper follows _CLI > Env Vars > Config File > Default_, `versionflag3`
should have taken the value `from configuration file`.

`versionflag4` is expected to be `from default`, but this is pure coincidence.


#### Using PersistenFlags from root cobra while running the sub command

In this situation, we try the cobra subcommand `version` with its 4 parameters set like previously.
And we also set the values of the _Persistent flags from rootCmd_.

Notice `versionflag{1, 2, 3, 4}` are still wrong whereas `rootflag{1, 2, 3, 4}`
follow the right priority order.

```bash
COBRAVSVIPER_ROOTFLAG2="from envvars"  COBRAVSVIPER_VERSIONFLAG2="from envvars" ./cobravsviper --rootflag1="from cli"  version --versionflag1="from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] version called                                cobra-cmd=version
INFO[0000] versionflag1: from configuration file         cobra-cmd=version
INFO[0000] versionflag2: from configuration file         cobra-cmd=version
INFO[0000] versionflag3: from configuration file         cobra-cmd=version
INFO[0000] versionflag4:                                 cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: from cli                           cobra-cmd=version
INFO[0000] rootflag2: from envvars                       cobra-cmd=version
INFO[0000] rootflag3: from configuration file            cobra-cmd=version
INFO[0000] rootflag4: from default                       cobra-cmd=version
```

We can even play with the `rootflagX` by changing which is defined through the CLI, etc...
The priority order still works for Root CMD persistent flags.

```bash
COBRAVSVIPER_ROOTFLAG3="from envvars"  COBRAVSVIPER_VERSIONFLAG2="from envvars" ./cobravsviper --rootflag2="from cli"  version --versionflag1="from cli" --config configs/cobravsviper.conf.yaml
```
```log
INFO[0000] version called                                cobra-cmd=version
INFO[0000] versionflag1: from configuration file         cobra-cmd=version
INFO[0000] versionflag2: from configuration file         cobra-cmd=version
INFO[0000] versionflag3: from configuration file         cobra-cmd=version
INFO[0000] versionflag4:                                 cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: from configuration file            cobra-cmd=version
INFO[0000] rootflag2: from cli                           cobra-cmd=version
INFO[0000] rootflag3: from envvars                       cobra-cmd=version
INFO[0000] rootflag4: from default                       cobra-cmd=version
```
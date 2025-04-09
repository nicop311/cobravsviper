# `cobravsviper`: make the most out of Golang Viper and Cobra

TL;DR: `viper.Sub()` does not maintain the right priority order (CLI > Env Vars > Config File > Default)
for flags of a cobra subcommand'.

This is an issue when unmarshaling `viper.Unmarshal`: each cobra subcommands will ignore configuration files or
cobra CLI or env var depending on how, when and in what order `viper.Unmarshal` is called.

- [1. How the project was bootstraped](#1-how-the-project-was-bootstraped)
- [2. Build The Project Locally](#2-build-the-project-locally)
- [3. CLI User Inputs Priority](#3-cli-user-inputs-priority)
  - [3.1. The root cobra command help message](#31-the-root-cobra-command-help-message)
  - [3.2. The `version` cobra subcommand help message](#32-the-version-cobra-subcommand-help-message)
  - [3.3. Priority In Theory: CLI \> Env Vars \> Config File \> Default](#33-priority-in-theory-cli--env-vars--config-file--default)
  - [3.4. In Practice:](#34-in-practice)
    - [3.4.1. Cobra Root Command: CLI Priority Success](#341-cobra-root-command-cli-priority-success)
    - [3.4.2. Cobra Sub Command: CLI Priority Failure](#342-cobra-sub-command-cli-priority-failure)
      - [v0.1.0](#v010)
      - [v0.2.0](#v020)
    - [3.4.3. Using PersistenFlags from root cobra while running the sub command](#343-using-persistenflags-from-root-cobra-while-running-the-sub-command)
      - [v0.1.0](#v010-1)
      - [v0.2.0](#v020-1)


## 1. How the project was bootstraped

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

## 2. Build The Project Locally

```bash
go build -o cobravsviper  main.go
```

## 3. CLI User Inputs Priority

In this project, [Cobra](https://github.com/spf13/cobra) is used to handel the CLI commands, subcommands
and all their flags.

[Viper](https://github.com/spf13/viper) is used with its features `viper.BindPFlags` and `viper.AutomaticEnv` which
allow Viper to automatically use Cobra's flags as both environment variables and configuration file.

This allow the developpers to only add and modify cobra flags, and viper will automatically adapt without the need for a
dedicated viper configuration.

### 3.1. The root cobra command help message

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
      --rootflag1 string   root flag 1 (default "value from default")
      --rootflag2 string   root flag 2 (default "value from default")
      --rootflag3 string   root flag 3 (default "value from default")
      --rootflag4 string   root flag 4 (default "value from default")
  -t, --toggle             Help message for toggle

Use "cobravsviper [command] --help" for more information about a command.
```

### 3.2. The `version` cobra subcommand help message

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
      --versionflag1 string   version flag 1 (default "value from default")
      --versionflag2 string   version flag 2 (default "value from default")
      --versionflag3 string   version flag 3 (default "value from default")
      --versionflag4 string   version flag 4 (default "value from default")

Global Flags:
      --config string      Configuration File (default "configs/myconfig.yaml")
      --rootflag1 string   root flag 1 (default "value from default")
      --rootflag2 string   root flag 2 (default "value from default")
      --rootflag3 string   root flag 3 (default "value from default")
      --rootflag4 string   root flag 4 (default "value from default")
```


### 3.3. Priority In Theory: CLI > Env Vars > Config File > Default

| Priority Order & Source  | Example                                                                     | Comment                         |
|--------------------------|-----------------------------------------------------------------------------|---------------------------------|
| 1️⃣ CLI Flag             | `--rootflag1 "value from cli"`                                                         | Highest priority                |
| 2️⃣ Environment Variable | `COBRAVSVIPER_ROOTFLAG2="value from envvars"`                                       | Overrides config file & default |
| 3️⃣ Config File          | `rootflag3: "value from config file"` in YAML/TOML/JSON                                         | Overrides default               |
| 4️⃣ Default Value        | `Flags().StringVar(&rootFlag4, "rootflag4", "value from default", "root flag 4")` | Used if nothing else is set     |

In this example:
* we choose the first root flag `--rootflag1` as the one to illustrate CLI flag user input, but this is arbitrary ; we could choose any of the 4 flags `--rootflagX`.
* we choose the second root flag `--rootflag2` (aka `COBRAVSVIPER_ROOTFLAG2`) as the one to illustrate CLI environment variable user input, but this is arbitrary ; we could choose any of the 4 env vars `COBRAVSVIPER_ROOTFLAGX`.
* we choose the third root flag `--rootflag3` (aka `rootflag3`) as the one to illustrate the configuration variable; but this is arbitrary.
* we choose the third root flag `--rootflag3` to test the defaul value.

This priority works for the root cobra command. But this priority does not work for the version command.

### 3.4. In Practice: 

#### 3.4.1. Cobra Root Command: CLI Priority Success

In this situation, `viper` and `cobra` behave as expected: the priority "CLI > Env Vars > Config File > Default" is repected.

> Note: in the config file, `rootflag4` is commented, which means the value is not set in the config file.
> This is to illustrate fallback to default if no user input is set for a flag.

```bash
COBRAVSVIPER_ROOTFLAG2="value from envvars" ./cobravsviper --rootflag1="value from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] Root command called                           cobra-cmd=cobravsviper
INFO[0000] rootflag1: value from cli                     cobra-cmd=cobravsviper
INFO[0000] rootflag2: value from envvars                 cobra-cmd=cobravsviper
INFO[0000] rootflag3: value from configuration file      cobra-cmd=cobravsviper
INFO[0000] rootflag4: value from default                 cobra-cmd=cobravsviper
```

Or run this without any arguments:output is taken from cobra default values.

```
./cobravsviper
```

```log
INFO[0000] Case config file from default location       
INFO[0000] Search config in .config directory /home/thedetective with name cobravsviper.conf.yaml (without extension). 
INFO[0000] Search config in home directory /home/thedetective/.config/cobravsviper with name cobravsviper.conf.yaml (without extension). 

INFO[0000] Root command called                           cobra-cmd=cobravsviper
INFO[0000] rootflag1: value from default                 cobra-cmd=cobravsviper
INFO[0000] rootflag2: value from default                 cobra-cmd=cobravsviper
INFO[0000] rootflag3: value from default                 cobra-cmd=cobravsviper
INFO[0000] rootflag4: value from default                 cobra-cmd=cobravsviper
```

Or witness the priority of cli flags over env var: `--rootflag4` overrides `COBRAVSVIPER_ROOTFLAG4`.

```bash
COBRAVSVIPER_ROOTFLAG4="value from envvar" ./cobravsviper --rootflag3 "value from cli"  --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] Root command called                           cobra-cmd=cobravsviper
INFO[0000] rootflag1: value from configuration file      cobra-cmd=cobravsviper
INFO[0000] rootflag2: value from configuration file      cobra-cmd=cobravsviper
INFO[0000] rootflag3: value from cli                     cobra-cmd=cobravsviper
INFO[0000] rootflag4: value from envvar                  cobra-cmd=cobravsviper
```

```bash
COBRAVSVIPER_ROOTFLAG4="value from envvar" ./cobravsviper --rootflag3 "value from cli" --rootflag4 "value from cli" --config configs/cobravsviper.conf.yaml 
```

```log
INFO[0000] Root command called                           cobra-cmd=cobravsviper
INFO[0000] rootflag1: value from configuration file      cobra-cmd=cobravsviper
INFO[0000] rootflag2: value from configuration file      cobra-cmd=cobravsviper
INFO[0000] rootflag3: value from cli                     cobra-cmd=cobravsviper
INFO[0000] rootflag4: value from cli                     cobra-cmd=cobravsviper
```

Or you can test any variation of user input, you will still have priority CLI > Env Vars > Config File > Default, which is nice.

#### 3.4.2. Cobra Sub Command: CLI Priority Failure

##### v0.1.0

In this situation `v0.1.0`, we try the cobra subcommand `version` and set only the flags corresponding to this subcommand.
Notice `viper` does not handle well the subcommand's flags.

Notice all `versionflag{1,2,3}` are set to configuration file, `versionflag4` is empty because the value is commente in the config file.

Ignore the output of the _Persistent flags from rootCmd_ for now: this output is perfectly normal, 
since the priority would be to use the values from the config file for `rootflag{1, 2, 3}` because
no other input was specified for these parameters. `rootflag4` is commented in the configuration,
and no other input was specified for this parameter, so it takes the default value.

```bash
COBRAVSVIPER_VERSIONFLAG2="value from envvars" ./cobravsviper version --versionflag1="value from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] version subcommand called                     cobra-cmd=version
INFO[0000] versionflag1: value from configuration file   cobra-cmd=version
INFO[0000] versionflag2: value from configuration file   cobra-cmd=version
INFO[0000] versionflag3: value from configuration file   cobra-cmd=version
INFO[0000] versionflag4:                                 cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: value from configuration file      cobra-cmd=version
INFO[0000] rootflag2: value from configuration file      cobra-cmd=version
INFO[0000] rootflag3: value from configuration file      cobra-cmd=version
INFO[0000] rootflag4: value from default                 cobra-cmd=version
```

Notice `versionflag3` has the value `from default` which comes from viper using cobra's default
values instead of `from configuration file`. If viper follows _CLI > Env Vars > Config File > Default_, `versionflag3`
should have taken the value `from configuration file`.

`versionflag4` is expected to be `from default`, but this is pure coincidence.

##### v0.2.0

With tag v0.2.0, we modify `cmd/version.go` to first unmarshall from config file (`viper.Sub("version").Unmarshal(&vprFlgsVersion)`), and then unmarshal from cobra autobindenv (`viper.Unmarshal(&vprFlgsVersion)`).
Notice the `versionflag3` does not get values from configuration file, since the second unmarshal (without `Sub`) overrides the first unmarshal (the one with `.Sub`).


```bash
COBRAVSVIPER_VERSIONFLAG2="value from envvars" ./cobravsviper version --versionflag1="value from cli" --config configs/cobravsviper.conf.yaml
```
```log
INFO[0000] version subcommand called                     cobra-cmd=version
INFO[0000] versionflag1: value from cli                  cobra-cmd=version
INFO[0000] versionflag2: value from envvars              cobra-cmd=version
INFO[0000] versionflag3: value from default              cobra-cmd=version
INFO[0000] versionflag4: value from default              cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: value from configuration file      cobra-cmd=version
INFO[0000] rootflag2: value from configuration file      cobra-cmd=version
INFO[0000] rootflag3: value from configuration file      cobra-cmd=version
INFO[0000] rootflag4: value from default                 cobra-cmd=version
```

#### 3.4.3. Using PersistenFlags from root cobra while running the sub command

##### v0.1.0

In this situation, we try the cobra subcommand `version` with its 4 parameters set like previously.
And we also set the values of the _Persistent flags from rootCmd_.

Notice `versionflag{1, 2, 3, 4}` are still wrong whereas `rootflag{1, 2, 3, 4}`
follow the right priority order.

```bash
COBRAVSVIPER_ROOTFLAG2="value from envvars"  COBRAVSVIPER_VERSIONFLAG2="value from envvars" ./cobravsviper --rootflag1="value from cli"  version --versionflag1="value from cli" --config configs/cobravsviper.conf.yaml
```

```log
INFO[0000] version subcommand called                     cobra-cmd=version
INFO[0000] versionflag1: value from configuration file   cobra-cmd=version
INFO[0000] versionflag2: value from configuration file   cobra-cmd=version
INFO[0000] versionflag3: value from configuration file   cobra-cmd=version
INFO[0000] versionflag4:                                 cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: value from cli                     cobra-cmd=version
INFO[0000] rootflag2: value from envvars                 cobra-cmd=version
INFO[0000] rootflag3: value from configuration file      cobra-cmd=version
INFO[0000] rootflag4: value from default                 cobra-cmd=version
```

We can even play with the `rootflagX` by changing which is defined through the CLI, etc...
The priority order still works for Root CMD persistent flags. But the subcommand does not use values from 
config file, since the

```bash
COBRAVSVIPER_ROOTFLAG3="value from envvars"  COBRAVSVIPER_VERSIONFLAG3="value from envvars" ./cobravsviper --rootflag2="value from cli"  version --versionflag3="value from cli" --config configs/cobravsviper.conf.yaml
```
```log
INFO[0000] version subcommand called                     cobra-cmd=version
INFO[0000] versionflag1: value from configuration file   cobra-cmd=version
INFO[0000] versionflag2: value from configuration file   cobra-cmd=version
INFO[0000] versionflag3: value from configuration file   cobra-cmd=version
INFO[0000] versionflag4:                                 cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: value from configuration file      cobra-cmd=version
INFO[0000] rootflag2: value from cli                     cobra-cmd=version
INFO[0000] rootflag3: value from envvars                 cobra-cmd=version
INFO[0000] rootflag4: value from default                 cobra-cmd=version
```

This problem comes from `viper.Sub("version").Unmarshal(&vprFlgsVersion)` that do not keep the priority order.

##### v0.2.0

With tag v0.2.0, we modify `cmd/version.go` to first unmarshall from config file (`viper.Sub("version").Unmarshal(&vprFlgsVersion)`), and then unmarshal from cobra autobindenv (`viper.Unmarshal(&vprFlgsVersion)`).
Notice the `versionflag3` does not get values from configuration file, since the second unmarshal (without `Sub`) overrides the first unmarshal (the one with `.Sub`).

But root flags are fine.

```bash
COBRAVSVIPER_ROOTFLAG2="value from envvars"  COBRAVSVIPER_VERSIONFLAG2="value from envvars" ./cobravsviper --rootflag1="value from cli"  version --versionflag1="value from cli" --config configs/cobravsviper.conf.yaml
```

```
COBRAVSVIPER_ROOTFLAG2="value from envvars"  COBRAVSVIPER_VERSIONFLAG2="value from envvars" ./cobravsviper --rootflag1="value from cli"  version --versionflag1="value from cli" --config configs/cobravsviper.conf.yaml
INFO[0000] version subcommand called                     cobra-cmd=version
INFO[0000] versionflag1: value from cli                  cobra-cmd=version
INFO[0000] versionflag2: value from envvars              cobra-cmd=version
INFO[0000] versionflag3: value from default              cobra-cmd=version
INFO[0000] versionflag4: value from default              cobra-cmd=version

INFO[0000] Persistent flags from rootCmd                 cobra-cmd=version
INFO[0000] rootflag1: value from cli                     cobra-cmd=version
INFO[0000] rootflag2: value from envvars                 cobra-cmd=version
INFO[0000] rootflag3: value from configuration file      cobra-cmd=version
INFO[0000] rootflag4: value from default                 cobra-cmd=version
```

# `cobravsviper`: make the most out of Golang Viper and Cobra

TL;DR: `viper.Sub()` does not maintain the right priority order (CLI > Env Vars > Config File > Default)
for flags of a cobra subcommand'.

This is an issue when unmarshaling `viper.Unmarshal`: each cobra subcommands will ignore configuration files or
cobra CLI or env var depending on how, when and in what order `viper.Unmarshal` is called.

- [1. How the project was bootstraped](#1-how-the-project-was-bootstraped)
- [2. Build The Project Locally](#2-build-the-project-locally)
  - [2.1. Build with `make` + `go`](#21-build-with-make--go)
  - [2.2. Build with `go`](#22-build-with-go)
- [3. CLI User Inputs Priority: the Theory and Issue in Practice](#3-cli-user-inputs-priority-the-theory-and-issue-in-practice)
  - [3.1. The root cobra command help message](#31-the-root-cobra-command-help-message)
  - [3.2. The `version` cobra subcommand help message](#32-the-version-cobra-subcommand-help-message)
  - [3.3. Priority In Theory: CLI \> Env Vars \> Config File \> Default](#33-priority-in-theory-cli--env-vars--config-file--default)
  - [3.4. In Practice:](#34-in-practice)
    - [3.4.1. Cobra Root Command: CLI Priority Success](#341-cobra-root-command-cli-priority-success)
    - [3.4.2. Cobra Sub Command: CLI Priority Failure](#342-cobra-sub-command-cli-priority-failure)
      - [3.4.2.1. v0.1.0](#3421-v010)
      - [3.4.2.2. v0.2.0](#3422-v020)
    - [3.4.3. Using PersistenFlags from root cobra while running the sub command](#343-using-persistenflags-from-root-cobra-while-running-the-sub-command)
      - [3.4.3.1. v0.1.0](#3431-v010)
      - [3.4.3.2. v0.2.0](#3432-v020)
- [4. CLI User Inputs Priority: my workaround for Viper and Cobra](#4-cli-user-inputs-priority-my-workaround-for-viper-and-cobra)
  - [4.1. Results of workarround](#41-results-of-workarround)


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

### 2.1. Build with `make` + `go`

```
make build-debug
```

or

```
make build
```

### 2.2. Build with `go`

```bash
go build -o cobravsviper  main.go
```

With debug:

```bash
go build -gcflags="all=-N -l" -o cobravsviper  main.go
```

## 3. CLI User Inputs Priority: the Theory and Issue in Practice

In this project, [Cobra](https://github.com/spf13/cobra) is used to handel the CLI commands, subcommands
and all their flags.

[Viper](https://github.com/spf13/viper) is used with its features `viper.BindPFlags` and `viper.AutomaticEnv` which
allow Viper to automatically use Cobra's flags as both environment variables and configuration file.

This allow the developpers to only add and modify cobra flags, and viper will automatically adapt without the need for a
dedicated viper configuration.

### 3.1. The root cobra command help message

> Note: This below corresponds to tag [v0.1.0](https://github.com/nicop311/cobravsviper/tree/v0.1.0). Check latest tag for workarround.

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

> Note: This below corresponds to tag [v0.1.0](https://github.com/nicop311/cobravsviper/tree/v0.1.0). Check latest tag for workarround.

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

> Note: This below corresponds to tag [v0.1.0](https://github.com/nicop311/cobravsviper/tree/v0.1.0). Check latest tag for workarround.

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

##### 3.4.2.1. v0.1.0

> Note: This below corresponds to tag [v0.1.0](https://github.com/nicop311/cobravsviper/tree/v0.1.0). Check latest tag for workarround.

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

##### 3.4.2.2. v0.2.0

> Note: This below corresponds to tag [v0.2.0](https://github.com/nicop311/cobravsviper/tree/v0.2.0). Check latest tag for workarround.

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

##### 3.4.3.1. v0.1.0

> Note: This below corresponds to tag [v0.1.0](https://github.com/nicop311/cobravsviper/tree/v0.1.0). Check latest tag for workarround.

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

##### 3.4.3.2. v0.2.0

> Note: This below corresponds to tag [v0.2.0](https://github.com/nicop311/cobravsviper/tree/v0.2.0). Check latest tag for workarround.

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

## 4. CLI User Inputs Priority: my workaround for Viper and Cobra

I found a workarround without modifying cobra nor viper. 

Look at file [`cmd/viper-patch-sub.go`](https://github.com/nicop311/cobravsviper/blob/48485b7e16c9a2837dbb45e743371d129fbab426/cmd/viper-patch-sub.go) with my patch & replacement for `viper.Sub` and `viper.Unmarshal`. I define the function [`UnmarshalSubMergedE`](https://github.com/nicop311/cobravsviper/blob/6111c0a815a3caf2616787f9989e50b0724ed20d/cmd/viper-patch-sub.go#L72), which is a replacement for `viper.Unmarshal` which supports input priority `flags > env > merged config > defaults`. And I define function [`InitViperSubCmdE`](https://github.com/nicop311/cobravsviper/blob/6111c0a815a3caf2616787f9989e50b0724ed20d/cmd/viper-patch-sub.go#L114) which does the binding between Viper and Cobra using my custom `UnmarshalSubMergedE` mathod and taking into account the YAML/TOML paths.

Last but not least, 

For a Level 1 cobra command : I need to call my function within a `cobra.PersistentPreRunE` like this:

```golang
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitViperSubCmdE(viper.GetViper(), cmd, &vprFlgsVersion); err != nil {
			logrus.WithField("cobra-cmd", cmd.Use).WithError(err).Error("Error initializing Viper")
			return err
		}
		return nil
	},
```

For a Level 2 cobra command : same idea, but I need to manually call the `cobra.PersistentPreRunE` command from the parent command, otherwise the persistentFlags from the parent commands are not taken into account.

```golang
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Manually call parent’s PersistentPreRunE
		if cmd.Parent() != nil && cmd.Parent().PersistentPreRunE != nil {
			if err := cmd.Parent().PersistentPreRunE(cmd.Parent(), args); err != nil {
				return err
			}
		}

		InitViperSubCmdE(viper.GetViper(), cmd, &vprFlgsSub221)
		return nil
	},
```

TODO: contribute to viper?
* https://github.com/spf13/viper/discussions/1756?sort=new#discussioncomment-12981228
* https://github.com/spf13/viper/issues/368#issuecomment-2838865972


### 4.1. Results of workarround

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

Group 1 Commands:
  grp1cmd1    A brief description of your command
  grp1cmd2    A brief description of your command

Group 2 Commands:
  grp2cmd1    A brief description of your command
  grp2cmd2    Test Nested Command of 1st level

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version information

Flags:
      --config string                Configuration File
      --debug                        Set logrus.SetLevel to "debug". This is equivalent to using --log-level=debug. Flags --log-level and --debug flag are mutually exclusive. Corresponding environment variable: K8S_KMS_PLUGIN_DEBUG.
  -h, --help                         help for cobravsviper
      --log-format string            Logrus log output format. Possible values: text, json. Corresponding environment variable: K8S_KMS_PLUGIN_LOG_FORMAT (default "text")
      --log-level string             Set logrus.SetLevel. Possible values: trace, debug, info, warning, error, fatal and panic. Flags --log-level and --debug flag are mutually exclusive. Corresponding environment variable: K8S_KMS_PLUGIN_LOG_LEVEL. (default "info")
      --rootflag1 string             root flag 1 (default "value from default")
      --rootflag2 string             root flag 2 (default "value from default")
      --rootflag3 string             root flag 3 (default "value from default")
      --rootflag4 string             root flag 4 (default "value from default")
      --rootpersistentflag1 string   persistent root flag 1 (default "value from default")
      --rootpersistentflag2 string   persistent root flag 2 (default "value from default")
      --rootpersistentflag3 string   persistent root flag 3 (default "value from default")
      --rootpersistentflag4 string   persistent root flag 4 (default "value from default")
  -t, --toggle                       Help message for toggle

Use "cobravsviper [command] --help" for more information about a command.
```

Run this example command which uses the 4 inputs (CLI > Env Vars > Config File > Default) for a level 2 cobra nested subcommand: 

```bash
COBRAVSVIPER_ROOTPERSISTENTFLAG2="value from envvars" \
COBRAVSVIPER_GRP2CMD2_GRP2CMD2PERSISTENTFLAG2="value from envvars" \
COBRAVSVIPER_GRP2CMD2_SUB221_SUB221FLAG2="value from envvars" \
COBRAVSVIPER_GRP2CMD2_SUB221_SUB221FLAGNOVAR2="value from envvars" \
cobravsviper \
  --rootpersistentflag1  "value from cli" \
  grp2cmd2 \
  --grp2cmd2persistentflag1  "value from cli" \
  sub221 \
  --sub221flag1  "value from cli" \
  --sub221flagnovar1  "value from cli" \
  --config  "configs/cobravsviper.conf.yaml"
```

Or run the vscode debug scenario `dlv vscode: sub221 env vars and config file` from [`.vscode/launch.json`](.vscode/launch.json).

Results below: every input is following the proper priority (CLI > Env Vars > Config File > Default).
```log
DEBU logrus log-level is set to: debug            

INFO flags from subcommand sub221                  cobra-cmd=sub221
INFO sub221flag1: value from cli                   cobra-cmd=sub221
INFO sub221flag2: value from envvars               cobra-cmd=sub221
INFO sub221flag3: value from YAML configuration file sub221 3  cobra-cmd=sub221
INFO sub221flag4: value from default               cobra-cmd=sub221

INFO sub221flagnovar1: value from cli              cobra-cmd=sub221
INFO sub221flagnovar2: value from envvars          cobra-cmd=sub221
INFO sub221flagnovar3: value from YAML configuration file sub221 3  cobra-cmd=sub221
INFO sub221flagnovar4: value from default 0.0.0.4  cobra-cmd=sub221

INFO Persistent flags from subcommand grp2cmd2     cobra-cmd=sub221
INFO grp2cmd2persistentflag1: value from cli       cobra-cmd=sub221
INFO grp2cmd2persistentflag2: value from envvars   cobra-cmd=sub221
INFO grp2cmd2persistentflag3: value from YAML configuration file grp2cmd2 3  cobra-cmd=sub221
INFO grp2cmd2persistentflag4: value from default   cobra-cmd=sub221

INFO Persistent flags from rootCmd                 cobra-cmd=sub221
INFO rootpersistentflag1: value from cli           cobra-cmd=sub221
INFO rootpersistentflag2: value from envvars       cobra-cmd=sub221
INFO rootpersistentflag3: value from YAML configuration file root 3  cobra-cmd=sub221
INFO rootpersistentflag4: value from default       cobra-cmd=sub221
```


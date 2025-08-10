# etime

The `time` command, with energy awareness.

**Outline**:

- [Getting started](#getting-started)
- [Measurement information](#measurement-information)
- [Program information](#program-information)
- [Repository information](#repository-information)

## Getting started

1. Purchase an [Emporia Smart Plug][plug] and set your device up.
2. Download the [latest released version][releases] or
   [build from source][source].
3. Optionally, create a symbolic link to or move the compiled binary into your
   `/bin` to run the command globally:

```sh
$ ln -s ~/path/to/emporia-time/etime /usr/local/bin

$ mv etime /usr/local/bin
```

4. Use the binary with your favorite command or script:

```sh
$ etime sleep 12
       12.00 real         0.00 user         0.00 sys
      922.63 joules      76.87 watts      100.0% sure
```

The first time you run `etime`, you will be prompted to login with your Emporia
credentials and select a device. Credentials are only used to gather API tokens,
and tokens are stored in `~/.config/etime/settings.json`.

## Measurement information

### Time

The duration of the input command is measured with the built-in `time` command.

Meanings of these measurements are as follows:

- `real`: the actual execution time from start to finish
- `user`: CPU time spent executing user-mode code for the process
- `sys`: CPU time spent making system calls in kernel mode

A more detailed explanation can be found in [this StackOverflow answer][time].

### Energy

Measurements of electricity used while executing the input command are collected
from the Smart Plug.

This usage is shown in the following units:

- `joules`: the total [energy][energy] used during the command duration
- `watts`: the average [power][power] output over the command duration
- `sure`: a confidence score for the above values

Results from [the Emporia API][docs] may not always be complete, so missing
usage is estimated by scaling the average measured energy over the total elapsed
time.

The ratio of observed-to-expected measurements is shown in the `sure` score.
Lookups are repeated until a sureness greater than 80.0% is achieved.

## Program information

### Usage guide

This program can be configured using positional flags and arguments to produce
certain behaviors:

```sh
$ etime [flags] <command> [args]
```

- `flags`: optional flags to provide this program
- `command`: the program to execute and measure
- `args`: optional arguments for the command

#### Flags

Configurations to this program can be made using `flags` before the `command`:

- `-h`, `--help`: display a hopefully helpful message
- `-p`, `--portable`: output measurements on separate lines
- `--device <string>`: name or ID of the smart plug to measure
- `--username <string>`: account username for Emporia
- `--password <string>`: account password for Emporia
- `--version`: print the current version of the build

#### Command

The provided `command` can be either a program or a path to an executable.
Pretty much anything that can be invoked from the command line.

#### Args

Any additional arguments for the `command` should follow the `command`. These
might include subcommands, positional values, or other flags.

### Program environment variables

Environment variables can be used as another way to configure the program:

- `EMPORIA_DEVICE`: name or ID of the smart plug to measure
- `EMPORIA_USERNAME`: account username for Emporia
- `EMPORIA_PASSWORD`: account password for Emporia
- `XDG_CONFIG_HOME`: the directory to store configurations

### Manual pages

Program documentation can be downloaded from root or the [releases][releases]
and added to reference:

```sh
$ cp etime.1 /usr/local/share/man/man1/
$ mandb
$ man etime
```

Additional permissions or created paths might be needed to complete the process.

## Repository information

This project is licensed under the MIT license and is not affiliated with or
endorsed by Emporia Energy.

Documentation for the Emporia API was graciously gathered from the
[`magico13/PyEmVue`][docs] project.

Notes on submitting contributions of any type are taken in
[`.github/CONTRIBUTING.md`][contributing].

Details on the processes around code for this repository are shared in the
[`.github/MAINTAINERS_GUIDE.md`][maintainers].

[contributing]: ./.github/CONTRIBUTING.md
[dashboard]: https://web.emporiaenergy.com/#/home
[docs]: https://github.com/magico13/PyEmVue/blob/master/api_docs.md
[energy]: https://en.wikipedia.org/wiki/Energy
[golang]: https://go.dev/dl
[maintainers]: ./.github/MAINTAINERS_GUIDE.md
[plug]: https://www.emporiaenergy.com/emporia-smart-plug
[power]: https://en.wikipedia.org/wiki/Power_(physics)
[releases]: https://github.com/zimeg/emporia-time/releases
[source]: ./.github/MAINTAINERS_GUIDE.md#project-setup
[time]: https://stackoverflow.com/a/556411

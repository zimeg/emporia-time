# etime

The `time` command, with energy awareness.

## Getting started

1. Purchase an [Emporia Smart Plug][plug] and set your device up.
2. Install the latest version of [Go][golang].
3. From a directory for development, download the source and compile `etime`:

```sh
$ git clone https://github.com/zimeg/emporia-time.git
$ cd emporia-time
$ make build
```

4. Optionally, create a symbolic link to or move the compiled binary into your
`/bin` to run the command globally:

```sh
$ ln -s ~/path/to/emporia-time/etime /usr/local/bin

$ mv etime /usr/local/bin
```

5. Use the binary with your favorite command or script:

```sh
$ ./etime sleep 12
       12.00 real         0.00 user         0.00 sys
        9.35 watt       100.0% sure
```

The first time you run `etime`, you will be prompted to login with your
Emporia credentials and select a device. Credentials are only used to gather
API tokens, and tokens are stored in `~/.config/etime/settings.json`.

## Measurement information

### Time

The duration of the input command is measured with the built-in `time` command.

Meanings of these measurements are as follows:

- `real`: The actual execution time from start to finish
- `user`: CPU time spent executing user-mode code for the process
- `sys`: CPU time spent making system calls in kernel mode

A more detailed explanation can be found in [this StackOverflow answer][time].

### Energy

The amount of electricity used during the execution of the input command is
collected from the Smart Plug and displayed in `watt`.

Results from [the Emporia API][docs] may not always be complete, so missing
usage is estimated by scaling the average measured energy over the total elapsed
time.

The ratio of observed-to-expected measurements is shown in the `sure` score.
Lookups are repeated until a sureness greater than 80.0% is achieved.

## Command information 

### Usage guide

This program can be configured using positional flags and arguments to produce
certain behaviors. Example usage follows:

```sh
$ ./etime [flags] <command> [args]
```

- `flags`: optional flags to provide this program
- `command`: the program to execute and measure
- `args`: optional arguments for the command

#### Flags

Configurations to this program can be made using `flags` before the `command`:

- `-h`, `--help`: output a hopefully helpful message
- `--device <string>`: name or ID of the smart plug to measure
- `--username <string>`: account username for Emporia
- `--password <string>`: account password for Emporia

#### Command

The provided command can be either a program or a path to an executable.
Pretty much anything that runs when invoked from the command line.

#### Args

Any additional arguments for the `command` should follow the `command`. These
might include subcommands, positional values, and other flags.

### Program environment variables

Environment variables can be used as another way to configure the program:

- `EMPORIA_DEVICE`: name or ID of the smart plug to measure
- `EMPORIA_USERNAME`: account username for Emporia
- `EMPORIA_PASSWORD`: account password for Emporia
- `XDG_CONFIG_HOME`: the directory to store configurations

<!-- links -->
[plug]: https://www.emporiaenergy.com/emporia-smart-plug
[golang]: https://go.dev/dl
[dashboard]: https://web.emporiaenergy.com/#/home
[time]: https://stackoverflow.com/a/556411
[docs]: https://github.com/magico13/PyEmVue/blob/master/api_docs.md

# etime

The `time` command, with energy awareness.

## Getting started

1. Purchase an [Emporia Smart Plug][plug] and set your device up.

2. From a directory for development, download the source and compile `etime`:

```sh
$ git clone https://github.com/e-zim/emporia-time.git
$ cd emporia-time
$ make build
```

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
        9.53 watt        61.5% sure
```

The first time you run `etime`, you will be prompted to login with your
Emporia credentials and select a device. Credentials are only used to gather
API tokens, and tokens are stored in `~/.config/etime/settings.json`.

## Measurement info

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

<!-- links -->
[plug]: https://www.emporiaenergy.com/emporia-smart-plug
[dashboard]: https://web.emporiaenergy.com/#/home
[time]: https://stackoverflow.com/a/556411
[docs]: https://github.com/magico13/PyEmVue/blob/master/api_docs.md

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

<!-- links -->
[plug]: https://www.emporiaenergy.com/emporia-smart-plug
[dashboard]: https://web.emporiaenergy.com/#/home

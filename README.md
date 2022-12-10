# etime

The `time` command, with energy awareness.

## Getting started

1. Purchase an [Emporia Smart Plug][plug].

2. Login to your [Emporia energy dashboard][dashboard] and naviage to the graphs
tab.

3. Open "Inspect Element" then open the "Network" tab. Click on an `AppAPI`
request and copy the `authToken`.

4. In your terminal, run the following commands with your own token and device:

```sh
$ export EMPORIA_TOKEN=your-authToken
$ export EMPORIA_DEVICE=your-deviceGid
```

5. From a directory for development, download the source and compile `etime`:

```sh
$ git clone https://github.com/e-zim/emporia-time.git
$ cd emporia-time
$ make build
```

6. Optionally, create a symbolic link to or move the compiled binary into your
`/bin` to run the command globally:

```sh
$ ln -s ~/path/to/emporia-time/etime /usr/local/bin

$ mv etime /usr/local/bin
```

7. Use the binary with your favorite command or script:

```sh
$ etime sleep 12
       12.00 real         0.00 user         0.00 sys
        9.53 watt        61.5% sure
```

<!-- links -->
[plug]: https://www.emporiaenergy.com/emporia-smart-plug
[dashboard]: https://web.emporiaenergy.com/#/home

# etime

The `time` command, with energy awareness.

## Getting started

1. Purchase an [Emporia Smart Plug][plug].

2. Login to your [Emporia energy dashboard][dashboard] and naviage to the graphs
tab.

3. Open "Inspect Element" then open the "Network" tab. Click on an `AppAPI`
request and locate the `deviceGid` you want to use. Save these for later.

4. From a directory for development, download the source and compile `etime`:

```sh
$ git clone https://github.com/e-zim/emporia-time.git
$ cd emporia-time
$ make build
```

5. Optionally, create a symbolic link to or move the compiled binary into your
`/bin` to run the command globally:

```sh
$ ln -s ~/path/to/emporia-time/etime /usr/local/bin

$ mv etime /usr/local/bin
```

6. Update your `~/.config/etime/settings.json` file with your `deviceGid` from
step 3 and Emporia credentials:

```json
{
  "EmporiaDevice": "012345",
  "EmporiaUsername": "you@email.com",
  "EmporiaPassword": "password123"
}
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

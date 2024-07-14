# Runners

A self-hosted runner is used to verify valid measurements are made when
monitoring energy usage during the remote integration tests.

To bring the runner online, [add a **New self-hosted runner**][runner] using a
device connected to a smart plug:

```sh
$ nix develop .#gh
$ cd .
$ config.sh # https://github.com/zimeg/emporia-time/settings/actions/runners/new
$ run.sh
```

Then set values for `EMPORIA_DEVICE`, `EMPORIA_USERNAME`, and `EMPORIA_PASSWORD`
in your action repository secrets using your Emporia information. Also add these
for Dependabot to configure this workflow.

[runner]: https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/adding-self-hosted-runners

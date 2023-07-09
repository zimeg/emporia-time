# Maintainers guide

Hey there! It's about time... Watt have you been jouling!?

**Outline**:

- [Project setup](#project-setup)
- [Testing](#testing)
- [Merging pull requests](#merging-pull-requests)
- [Runner setup](#runner-setup)

## Project setup

After setting up the project for normal usage, you're ready for development!

An [understanding of Go][learn_go] is a likely prerequisite for any programming
and can be an enjoyable language to learn!

### Project structure

This project hopes to use different directories to separate various concerns,
currently using the following structure: 

- `/` – project files
- `.github/` – information for collaboration and continuous integrations
- `internal/` – helpful utilities needed to create the program
- `pkg/` – various concerns that are pieced together to form the program

### Makefile commands

For ease of development, some commands are added in a `Makefile`:

- `make build` – build the program binary
- `make test` – perform the written code tests
- `make clean` – remove all program artifacts

The name of the binary can be changed with `make build BIN=timer`.

## Testing

All tests should aim to cover common cases, both happy and erroneous, to build
confidence in any changes.

### Unit tests

Written tests should reside in a file adjacent to the functionality being tested
and suffixed with `_test.go`.

All tests can be run with `make test` and example test cases can be found
throughout the repo.

While coverage isn't critical, various permutations of input are often used to
check edge cases. There's some balance.

### Integration tests

Assurance that the program works as expected with the Emporia API can be gained
by running the program with a command:

```sh
make build
./etime sleep 4
```

A smart plug and Emporia credentials are needed for this to be successful.

### On the remote

When changes are proposed or made to the remote repository, the full test suite
is performed to verify stability in any changes.

Additionally, some change to the `CHANGELOG.md` is checked for on pull requests.

## Merging pull requests

Confidence in the tests should cover edge cases well enough to trust the suite.
A green status signals nothing broke as a result of changes, and an example run
can be seen in the actions output.

On any change, the following should be verified before merging:

- Documentation is correct and updated everywhere necessary
- Code changes move the project in a positive direction

If that all looks good and the change is solid, the **Squash and merge** awaits.

## Runner setup

A self-hosted runner is used to verify valid measurements are made when
monitoring energy usage during the remote integration tests.

To bring runner online, [add a **New self-hosted runner**][runner] using a
device connected to a smart plug.

Then set values for `EMPORIA_DEVICE`, `EMPORIA_USERNAME`, and `EMPORIA_PASSWORD`
in your action repository secrets using your Emporia information. Also add these
for Dependabot to configure this workflow.

<!-- a collection of links -->
[learn_go]: https://go.dev/learn/
[runner]: https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/adding-self-hosted-runners

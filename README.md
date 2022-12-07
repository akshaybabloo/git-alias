# Prettify Git Alias

This is literally the most over-engineered thing I've ever made. It prints out a list of git aliases in a pretty way. It's a bit of a joke, but it's also useful.

## Usage

Add the binary file to your path, and then run `git-alias`. It will print out a list of aliases, and the commands they run.

```bash
$ git-alias
# or, add it to your alias
$ git config --global alias.la '!git-alias'
$ git la
```

Search by using the `--search` or `-s` flag.

```bash
$ git-alias -s "log"
```

Override the default Git configuration file by using the `--config` flag.

```bash
$ git-alias --config "~/.gitconfig"
```

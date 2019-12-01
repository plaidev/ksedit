# ksedit

kubernetes secret edit

## usage

pipe or filepath

```bash
$ cat ./secret.yml | ksedit > ./scret-edited.yml
# or
$ ksedit ./secret.yml > ./scret-edited.yml
```

edit and write

```bash
$ ksedit -w ./secret.yml
```

manually decode and encode

```bash
$ ksedit -d ./secret.yml > ./secret.dec.yml
$ vim ./secret.dec.yml
$ ksedit -e ./secret.dec.yml > ./secret.yml
```

## options

```bash
NAME:
   ksedit - kubernetest secret resource edit

USAGE:
   ksedit [global options] filepath

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --write, -w     write secret
   --encode, -e    encode secret
   --decode, -d    decode secret
   --editor value  editor (default: "vim")
   --help, -h      show help
   --version, -v   print the version
```


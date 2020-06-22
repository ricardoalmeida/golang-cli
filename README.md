CLI to process big JSON files
====

Install and set up
----

Please clone the project in the directory:

```sh
docker build -t golang-cli .
docker create -t -i golang-cli bash
```

The command will return a HASH, needed for the next step.

```sh
docker start -a -i <HASH>
```

And inside the machine, the command line is already available for using without arguments (defaults)

```sh
golang-cli stats
```

For help please use:

```sh
golang-cli stats -h
```

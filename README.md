# suez

A command line tool that executes a command periodically/reactively.

## Usage

### Run the command periodically according to the `interval-sec`

```shell
$ suez --interval-sec 60 something-command you-want
```

### Run the command periodically if the watched file is changed while the duration of `interval-sec`

```shell
$ suez --interval-sec 60 --watched-file /path/to/file something-command you-want
```

### Run the command immediately if the watched file is changed

```shell
$ suez --watched-file /path/to/file something-command you-want
```

## Author

moznion (<moznion@gmail.com>)


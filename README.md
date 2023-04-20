# Go sample project

Sample repo to quickstart golang applications. Plz note that this is designed to demonstrate common used features, not best practices. E.g. the custom flagset usage here is mostly superficial, in a program like this it's better to just use default flagset.

## Functionality

Sample functionality includes:

- Reading data from stdin
- Making some api requests
- Outputting data to stdout

## Environment setup

Use [gvm](https://github.com/moovweb/gvm)

For first-time setup use binary install. Note that `--default` is necessary so that the version is available in PATH:

```zsh
gvm install go1.20.3 --prefer-binary --default
```

After doing all of this, install the go extension for vscode and configure the following in the vscode settings (json):

```json
{
    //...
    "go.goroot": "~/.gvm/gos/go1.20.3",
}
```

Unfortunately this is needed as there is no way to currently marry VSCode with gvm. If you want to switch your default go version, you will need to update this setting also.

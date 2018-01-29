# webd

A website server that can serve static content.

The content can be plain files or HTML templates and YAML content.
It is served via HTTP or HTTPS using Let's Encrypt with automatic certificate renewal.

# Installation

## Using `go get`.
If you have the go language installed, run `go get` to download the source code and build the `webd` server.

Once you have a built binary you don't need the go runtime on your prod server, just the webd or webd.exe file and your site templates.

```bash
$ go get github.com/robert-wallis/webd
```

## Running a set of sites.

In the [example/](example/) folder you will find a [sites.yml](example/sites.yaml) file.

```yaml
-
  host: localhost
  aliases: [www.example.com, local.example.com]
  email: test@example.com
  path: .
  bind:
    http: :80
  letsencrypt: true
-
  host: files.example.com
  static: true
  path: ../test_data/files.example.com
  bind:
    http: :80
```

Each site is a seperate section (yaml document).

If you want to serve https then make sure you use `https: :443` in the bind section.

`static` means don't serve templates, but just any files in the `path` folder.

`path` is relative to the configuration yaml file's location.

To run a site using the example sites.yml file run:

```
webd example/sites.yaml
```

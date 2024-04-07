# Download Progress
This example demonstrates how to download a file from a URL and show its progress with a [Progress Bubble][progress].

In this case we're getting download progress with an [`io.TeeReader`][tee] and sending progress `Msg`s to the `Program` with `Program.Send()`.

[progress]: https://github.com/rprtr258/tea/components/
[tee]: https://pkg.go.dev/io#TeeReader

## How to Run
Build the application with `go build .`, then run with a `--url` argument specifying the URL of the file to download. For example:

```bash
./progress-download --url="https://download.blender.org/demo/color_vortex.blend"
```

Note that in this example a TUI will not be shown for URLs that do not respond with a ContentLength header.

* * *
This example originally came from [this discussion][discussion].

[discussion]: https://github.com/charmbracelet/bubbles/discussions/127

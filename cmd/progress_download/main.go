package progress_download //nolint:revive,stylecheck

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/progress"
)

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64) tea.Msg
	send       func(tea.Msg)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.send(pw.onProgress(float64(pw.downloaded) / float64(pw.total)))
	}
	return len(p), nil
}

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	if _, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw)); err != nil {
		pw.send(msgProgressErr{err})
	}
}

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) // nolint:gosec,noctx
	if err != nil {
		log.Fatalln(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func Main(ctx context.Context) error {
	url := flag.String("url", "https://download.blender.org/demo/color_vortex.blend", "url for the file to download")
	flag.Parse()

	resp, err := getResponse(*url)
	if err != nil {
		log.Fatalln("could not get response", err.Error())
	}
	defer resp.Body.Close() // nolint:errcheck

	// Don't add TUI if the header doesn't include content size
	// it's impossible see progress without total
	if resp.ContentLength <= 0 {
		return errors.New("can't parse content length, aborting download")
	}

	filename := filepath.Base(*url)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close() // nolint:errcheck

	var p *tea.Program[*tea.AdapterModel[*model]]

	pw := &progressWriter{
		total:  int(resp.ContentLength),
		file:   file,
		reader: resp.Body,
		send: func(m tea.Msg) {
			p.Send(m)
		},
		onProgress: func(ratio float64) tea.Msg {
			return msgProgress(ratio)
		},
	}

	m := &model{
		pw:       pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p = tea.NewProgram2(ctx, m)

	// Start the download
	go pw.Start()

	_, err = p.Run()
	return err
}

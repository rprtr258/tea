package progress_download

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
	"github.com/rprtr258/tea/bubbles/progress"
)

var p *tea.Program[*model]

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		p.Send(msgProgressErr{err})
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url) // nolint:gosec
	if err != nil {
		log.Fatalln(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func Main(ctx context.Context) error {
	url := flag.String("url", "", "url for the file to download")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		return errors.New("url is required")
	}

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

	pw := &progressWriter{
		total:  int(resp.ContentLength),
		file:   file,
		reader: resp.Body,
		onProgress: func(ratio float64) {
			p.Send(msgProgress(ratio))
		},
	}

	// Start the download
	go pw.Start()

	m := &model{
		pw:       pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}

	_, err = tea.NewProgram(ctx, m).Run()
	return err
}

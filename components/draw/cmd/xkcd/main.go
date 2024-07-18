package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	drawille "github.com/rprtr258/tea/components/draw"
	"golang.org/x/image/draw"
)

func usage() {
	println("Usage: %s <url/id>")
}

var re = regexp.MustCompile(`src="(//imgs.xkcd.com/comics/.*\.)(jpg|png)"`)

func getTerminalSize() (w, h int) {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	var ws winsize
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col), int(ws.Row)
}

func luminocity(r, g, b uint32) uint64 {
	return (uint64(r)*299 + uint64(g)*587 + uint64(b)*114) / 1000
}

func loadImage(url string) image.Image {
	img, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer img.Body.Close()

	i, _, err := image.Decode(img.Body)
	if err != nil {
		panic(err)
	}
	return i
}

func xkcd(id string) image.Image {
	var url string
	if id == "" {
		url = "http://xkcd.com/"
	} else {
		url = fmt.Sprintf("http://xkcd.com/%s/", id)
	}

	c, _ := http.Get(url)
	defer c.Body.Close()
	cb, _ := io.ReadAll(c.Body)

	img_url := "https:" + strings.Join(re.FindStringSubmatch(string(cb))[1:], "")
	return loadImage(img_url)
}

func main() {
	argv := os.Args
	var i image.Image
	switch {
	case len(argv) < 2:
		i = xkcd("")
	case argv[1] == "-h" || argv[1] == "--help":
		usage()
		return
	case strings.HasPrefix(argv[1], "http"):
		i = loadImage(argv[1])
	default:
		if _, err := strconv.ParseInt(argv[1], 10, 64); err == nil {
			i = xkcd(argv[1])
		} else {
			f, err := os.Open(argv[1])
			if err != nil {
				panic(err)
			}
			defer f.Close()

			i, _, err = image.Decode(f)
			if err != nil {
				panic(err)
			}
		}
	}

	tw, th := getTerminalSize()

	var w, h int
	if b := i.Bounds(); b.Dx() > b.Dy() {
		w, h = 2*tw, 4*2*th*b.Dy()/b.Dx()
	} else {
		w, h = 2*tw*b.Dx()/b.Dy(), 4*th
	}

	// Resize
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dst, dst.Rect, i, i.Bounds(), draw.Over, nil)

	can := drawille.NewCanvas()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, _ := dst.At(x, y).RGBA()
			if luminocity(r, g, b) < 65535/2 {
				can.SetN(x, y)
			}
		}
	}
	fmt.Println(can)
}

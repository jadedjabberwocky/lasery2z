package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jadedjabberwocky/lasery2z/coordmap"
	"github.com/jadedjabberwocky/lasery2z/coordmap/imagemap"
	"github.com/jadedjabberwocky/lasery2z/gcode"
	"github.com/jadedjabberwocky/lasery2z/noopio"
)

// Errors
var (
	ErrUnsupportedCommand = errors.New("Unsupported command")
)

func openReader(filename string) (io.ReadCloser, error) {
	if filename == "" {
		return &noopio.ReadCloser{}, nil
	}

	if filename == "-" {
		return os.Stdin, nil
	}

	return os.Open(filename)
}

func openWriter(filename string) (io.WriteCloser, error) {
	if filename == "" {
		return &noopio.WriteCloser{}, nil
	}

	if filename == "-" {
		return os.Stdin, nil
	}

	if filename == "!" {
		return os.Stdout, nil
	}

	return os.Create(filename)
}

func run() error {
	var (
		cfg = &Config{}
		cm  coordmap.CoordMap
	)

	ctx := kong.Parse(cfg, kong.UsageOnError())

	cmd := ctx.Command()
	switch cmd {
	case "image-map":
		o, err := imagemap.DefaultOptions().
			WithImageWidth(cfg.ImageMap.ImageWidth).
			WithImageHeight(cfg.ImageMap.ImageHeight).
			Check()
		if err != nil {
			return err
		}

		cm, err = imagemap.New(cfg.ImageMap.ImageInputFilename, o)
		if err != nil {
			return err
		}

	default:
		return ErrUnsupportedCommand
	}

	cmw, err := openWriter(cfg.MapOutputFilename)
	if err != nil {
		return err
	}
	defer cmw.Close()

	_, _ = cmw.Write([]byte(cm.String()))

	gin, err := openReader(cfg.InputFilename)
	if err != nil {
		return err
	}
	defer gin.Close()

	gout, err := openWriter(cfg.OutputFilename)
	if err != nil {
		return err
	}
	defer gout.Close()

	err = gcode.Process(gin, gout, cm)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

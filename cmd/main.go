package main

import (
	"log"
	"os"
	"runtime"

	"github.com/burnb/ankifiller/internal/anki"
	"github.com/burnb/ankifiller/internal/configs"
	"github.com/burnb/ankifiller/internal/filler"
	"github.com/burnb/ankifiller/pkg/image"
	"github.com/burnb/ankifiller/pkg/phonemic"
)

var exitCode int

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		os.Exit(exitCode)
	}()

	cfg := &configs.App{}
	if err := cfg.Prepare(); err != nil {
		exitWithError(err)
	}

	phonemicSrv := phonemic.NewService(cfg.Phonemic)
	if phonemicSrv != nil {
		if err := phonemicSrv.Init(); err != nil {
			exitWithError(err)
		}
	}

	imageSrv := image.NewService(cfg.GoogleCustomSearch)
	ankiSrv := anki.NewService(cfg.Anki)

	fillerSrv := filler.NewService(ankiSrv, imageSrv, phonemicSrv)
	if err := fillerSrv.Run(); err != nil {
		exitWithError(err)
	}

	log.Print("Done.")
}

func exitWithError(err error) {
	exitCode = 1
	log.Println(err)
	runtime.Goexit()
}

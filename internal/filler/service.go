package filler

import (
	"fmt"
	"log"

	"github.com/burnb/ankifiller/internal/anki"
	"github.com/burnb/ankifiller/pkg/image"
	"github.com/burnb/ankifiller/pkg/phonemic"
)

type Service struct {
	ankiSrv     *anki.Service
	imageSrv    *image.Service
	phonemicSrv *phonemic.Service
}

func NewService(ankiSrv *anki.Service, imageSrv *image.Service, phonemicSrv *phonemic.Service) *Service {
	return &Service{ankiSrv: ankiSrv, imageSrv: imageSrv, phonemicSrv: phonemicSrv}
}

func (s *Service) Run() error {
	notes, err := s.ankiSrv.Notes()
	if err != nil {
		return err
	}
	for _, note := range notes {
		if s.imageSrv != nil && !note.SkipImageUpdate {
			img, imgErr := s.imageSrv.ImageByDefinition(note.ImageDefinition)
			if imgErr != nil {
				return imgErr
			}
			note.Image = img

			log.Println(fmt.Sprintf("Added card '%s' image", note.ImageDefinition))
		}

		if s.phonemicSrv != nil && !note.SkipPhonemicUpdate {
			note.Phonemic = s.phonemicSrv.Transcript(note.PhonemicDefinition)

			log.Println(fmt.Sprintf("Added card '%s' phonemic '%s'", note.PhonemicDefinition, note.Phonemic))
		}

		if sendErr := s.ankiSrv.UpdateNote(note); sendErr != nil {
			return sendErr
		}
	}

	return nil
}

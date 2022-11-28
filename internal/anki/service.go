package anki

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atselvan/ankiconnect"

	"github.com/burnb/ankifiller/internal/configs"
)

type Service struct {
	cfg    configs.Anki
	client *ankiconnect.Client
}

func NewService(cfg configs.Anki) *Service {
	return &Service{cfg: cfg, client: ankiconnect.NewClient()}
}

func (s *Service) Notes() (notes []*Note, err error) {
	allNotesPtr, restErr := s.client.Notes.Get(fmt.Sprintf("deck:%s", s.cfg.Deck))
	if restErr != nil {
		return nil, errors.New(restErr.Message)
	}

	allNotes := *allNotesPtr
	if len(allNotes) > 0 {
		firstNote := allNotes[0]

		if s.cfg.ImageField != nil {
			if _, ok := firstNote.Fields[*s.cfg.ImageField]; !ok {
				return nil, fmt.Errorf("unable to find image field %s", *s.cfg.ImageField)
			}
		}

		if s.cfg.ImageDefinitionField != nil {
			if _, ok := firstNote.Fields[*s.cfg.ImageDefinitionField]; !ok {
				return nil, fmt.Errorf("unable to find image definition field %s", *s.cfg.ImageDefinitionField)
			}
		}

		if s.cfg.PhonemicField != nil {
			if _, ok := firstNote.Fields[*s.cfg.PhonemicField]; !ok {
				return nil, fmt.Errorf("unable to find phonemic field %s", *s.cfg.PhonemicField)
			}
		}

		if s.cfg.PhonemicDefinitionField != nil {
			if _, ok := firstNote.Fields[*s.cfg.PhonemicDefinitionField]; !ok {
				return nil, fmt.Errorf("unable to find phonemic definition field %s", *s.cfg.PhonemicDefinitionField)
			}
		}
	}

	for _, note := range allNotes {
		emptyNote := &Note{
			Id:                 note.NoteId,
			SkipImageUpdate:    s.cfg.ImageDefinitionField == nil || s.cfg.ImageField == nil || note.Fields[*s.cfg.ImageField].Value != "",
			SkipPhonemicUpdate: s.cfg.PhonemicDefinitionField == nil || s.cfg.PhonemicField == nil || note.Fields[*s.cfg.PhonemicField].Value != "",
		}

		if emptyNote.SkipImageUpdate && emptyNote.SkipPhonemicUpdate {
			continue
		}

		if !emptyNote.SkipImageUpdate {
			imgDfnField, _ := note.Fields[*s.cfg.ImageDefinitionField]
			imgDfnParts := strings.Split(imgDfnField.Value, " [sound")
			emptyNote.ImageDefinition = imgDfnParts[0]
		}

		if !emptyNote.SkipPhonemicUpdate {
			phmcDfnField, _ := note.Fields[*s.cfg.PhonemicDefinitionField]
			phmcDfnParts := strings.Split(phmcDfnField.Value, " [sound")
			emptyNote.PhonemicDefinition = phmcDfnParts[0]
		}

		notes = append(notes, emptyNote)
	}

	return notes, nil
}

func (s *Service) UpdateNote(note *Note) error {
	updateNote := ankiconnect.UpdateNote{Id: note.Id, Fields: make(map[string]string)}

	if !note.SkipImageUpdate {
		fieldName := *s.cfg.ImageField
		updateNote.Fields[fieldName] = ""
		updateNote.Picture = append(
			updateNote.Picture,
			ankiconnect.Picture{URL: note.Image.Url, Filename: note.Image.FileName, Fields: []string{fieldName}},
		)
	}

	if !note.SkipPhonemicUpdate {
		fieldName := *s.cfg.PhonemicField
		updateNote.Fields[fieldName] = ""
		updateNote.Fields[fieldName] = note.Phonemic
	}

	if err := s.client.Notes.Update(updateNote); err != nil {
		return errors.New(err.Message)
	}

	return nil
}

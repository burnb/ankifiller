package anki

import (
	"github.com/burnb/ankifiller/pkg/models"
)

type Note struct {
	Id                 int64
	ImageDefinition    string
	Image              *models.Image
	PhonemicDefinition string
	Phonemic           string
	SkipImageUpdate    bool
	SkipPhonemicUpdate bool
}

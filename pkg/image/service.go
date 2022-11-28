package image

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/burnb/ankifiller/pkg/models"
)

type Service struct {
	cfg    Config
	client *http.Client
}

func NewService(cfg Config) *Service {
	if cfg == nil {
		return nil
	}

	return &Service{cfg: cfg, client: http.DefaultClient}
}

func (s *Service) ImageByDefinition(definition string) (*models.Image, error) {
	query := url.Values{}
	query.Add("q", definition)
	query.Add("key", s.cfg.GetAPIKey())
	query.Add("cx", s.cfg.GetCx())
	query.Add("searchType", "image")
	query.Add("imgSize", "xlarge")
	if s.cfg.GetGl() != nil {
		query.Add("gl", *s.cfg.GetGl())
	}

	requestUrl := &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "customsearch/v1",
		RawQuery: query.Encode(),
	}

	httpResp, err := s.client.Get(requestUrl.String())
	if err != nil {
		return nil, fmt.Errorf("unable to do request %w", err)
	}

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body %w", err)
	}

	resp := &Response{}
	if err := resp.UnmarshalJSON(respBody); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body %w", err)
	}

	if resp.Items == nil {
		return nil, errors.New("url not fond")
	}

	var img *models.Image
	for _, item := range resp.Items {
		fileExtList, err := mime.ExtensionsByType(item.Mime)
		if err == nil && len(fileExtList) != 0 {
			img = &models.Image{FileName: definition + fileExtList[0], Url: item.Link}
			break
		}
	}

	return img, nil
}

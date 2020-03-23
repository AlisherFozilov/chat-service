package fileconnector

import (
	"bytes"
	"github.com/AlisherFozilov/chat-service/pkg/rest"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type RemoteURL string

type Service struct {
	url string
}

func NewFileSvcConnector(url RemoteURL) *Service {
	if url == "" {
		panic("url can't be empty")
	}
	return &Service{url: string(url)}
}

func (s *Service) SaveOnFileServiceAndGetUrls(message []byte) ([]FileURL, error) {

	client := &http.Client{}

	values := map[string]io.Reader{
		"file": bytes.NewReader(message),
	}
	response, err := s.upload(client, s.url, values)
	if err != nil {
		return nil, err
	}

	fileURL := []FileURL{}
	err = rest.ReadJSONBody(response, &fileURL)
	if err != nil {
		log.Print(err)
	}

	return fileURL, nil
}

func (s *Service) upload(client *http.Client, url string,
	values map[string]io.Reader) (response *http.Response, err error) {

	var buffer bytes.Buffer
	writerBuffer := multipart.NewWriter(&buffer)
	for key, reader := range values {
		err := formMulipartFile(key, reader, writerBuffer)
		if err != nil {
			log.Println(err)
		}
	}

	err = writerBuffer.Close()
	if err != nil {
		log.Print(err)
	}

	request, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", writerBuffer.FormDataContentType())

	response, err = client.Do(request)
	if err != nil {
		return
	}

	return response, nil
}

func formMulipartFile(key string, reader io.Reader, writerBuffer *multipart.Writer) (err error) {
	var fileWriter io.Writer
	if closer, ok := reader.(io.Closer); ok {
		defer func() {
			errdefer := closer.Close()
			if errdefer != nil {
				log.Println("can't close", errdefer)
			}
		}()
	}
	if fileWriter, err = writerBuffer.CreateFormFile(key, "-"); err != nil {
		return nil
	}

	_, err = io.Copy(fileWriter, reader)
	if err != nil {
		return err
	}
	return nil
}

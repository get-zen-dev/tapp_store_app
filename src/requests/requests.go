package requests

import (
	"context"
	env "environment"
	"errors"
	"os"

	"github.com/google/go-github/v51/github"
)

// Returns a Response object and a nil error. If the request failed, then the error is not nil,
// but the response is nil
func GetListAddons() (*Response, error) {
	client := github.NewClient(nil)
	_, directoryContent, res, err := client.Repositories.GetContents(
		context.Background(),
		env.GetOwner(), env.GetRepository(), env.GetPath(),
		&github.RepositoryContentGetOptions{Ref: env.GetRef()})
	if res.Rate.Remaining == 0 && err != nil {
		e := errors.New("number of requests available early 0")
		errors.Join(err, e)
	}
	if err != nil {
		return nil, err
	}
	resp := Response{}
	for i := 0; i < len(directoryContent); i++ {
		model := newModel(directoryContent[i].GetName(), directoryContent[i].GetPath(), directoryContent[i].GetURL())
		resp.append(model)
	}
	return &resp, nil
}

// Download information about addons. Returns an error on failure or nil on success
func DownloadInfoAddons() error {
	client := github.NewClient(nil)
	fileContent, _, res, err := client.Repositories.GetContents(
		context.Background(),
		env.GetOwner(), env.GetRepository(), "addons.yaml",
		&github.RepositoryContentGetOptions{Ref: env.GetRef()})
	if res.Rate.Remaining == 0 && err != nil {
		e := errors.New("number of requests available early 0")
		errors.Join(err, e)
	}
	if err != nil {
		return err
	}
	content, err := fileContent.GetContent()
	os.WriteFile("./../configs/addons.yaml", []byte(content), os.FileMode(0666))
	if err != nil {
		return err
	}
	return nil
}

package requests

import (
	"context"
	env "environment"
	"errors"

	"github.com/google/go-github/v51/github"
)

// Returns a Response object and a nil error. If the request failed, then the error is not nil,
// but the response is nil
func GetListAddons() (*Response, error) {
	client := github.NewClient(nil)
	_, directoryContent, res, err := client.Repositories.GetContents(
		context.Background(),
		env.GetOwner(), env.GetRepository(), env.GetPath(),
		&github.RepositoryContentGetOptions{})
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

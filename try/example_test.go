package try_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/martinohmann/exp/try"
)

type Release struct {
	Name string `json:"name"`
}

func listReleasesIdiomatic(owner, repo string) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []Release

	if err := json.Unmarshal(buf, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func listReleasesTryRun(owner, repo string) ([]Release, error) {
	var releases []Release

	err := try.Run(func() {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		try.Check(err)

		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := http.DefaultClient.Do(req)
		try.Check(err)
		defer resp.Body.Close()

		buf, err := io.ReadAll(resp.Body)
		try.Check(err)
		try.Check(json.Unmarshal(buf, &releases))
	})

	return releases, err
}

func Example() {
	releases, err := listReleasesTryRun("martinohmann", "keyring")
	if err != nil {
		panic(err)
	}

	fmt.Println(releases)

	releases, err = listReleasesIdiomatic("martinohmann", "keyring")
	if err != nil {
		panic(err)
	}

	fmt.Println(releases)
}

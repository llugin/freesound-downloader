package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/llugin/freesound-downloader/config"
)

type Client struct {
	Config config.APIConfig
}

type SearchRequest struct {
	Query string
}

type SearchResult struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Results  []Result `json:"results"`
	Previous string   `json:"previous"`
}

type Result struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	License  string   `json:"license"`
	Rating   float64  `json:"avg_rating"`
	Duration float64  `json:"duration"`
	Download string   `json:"download"`
	Tags     []string `json:"tags"`
	Created  string   `json:"created"`
}

type Query struct {
	Query    string
	MaxLen   int
	PageSize int
	Page     int
}

type AccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) Authorize() error {

	u, err := url.Parse("https://freesound.org/apiv2/oauth2/authorize/")
	if err != nil {
		return err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	q.Add("client_id", c.Config.ClientID)
	q.Add("response_type", "code")
	u.RawQuery = q.Encode()
	cmd := exec.Command("open", u.String())
	cmd.Run()
	return nil
}

func (c *Client) GetNewest(query Query) (*SearchResult, error) {

	u, err := url.Parse("https://freesound.org/apiv2/search/text/")
	if err != nil {
		return nil, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	q.Add("query", query.Query)
	q.Add("page", strconv.Itoa(query.Page))
	q.Add("token", c.Config.ApiKey)
	q.Add("sort", "created_desc")
	q.Add("filter", fmt.Sprintf("duration:[0 TO %d]", query.MaxLen))
	q.Add("license", "Creative Commons 0")
	q.Add("fields", "type,created,name,download,duration,avg_rating,license,tags")
	if query.PageSize != 0 {
		q.Add("page_size", strconv.Itoa(query.PageSize))
	}

	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	var res SearchResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetAccessToken(authCode string) (*AccessResponse, error) {
	resp, err := http.PostForm(
		"https://freesound.org/apiv2/oauth2/access_token/",
		url.Values{"client_id": {c.Config.ClientID},
			"client_secret": {c.Config.ClientSecret},
			"grant_type":    {"authorization_code"},
			"code":          {authCode}},
	)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(body))
	}

	var ar AccessResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, err
	}
	return &ar, nil

}

func (c *Client) DownloadOne(result *Result, dir, accessToken string) error {
	cl := &http.Client{}

	req, err := http.NewRequest("GET", result.Download, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}

	filename := strings.ReplaceAll(result.Name, " ", "_")

	fullpath := path.Join(dir, addExtension(filename, result.Type))
	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		for i := 0; ; i++ {
			fullpath = path.Join(
				dir,
				addExtension(fmt.Sprintf("%v_%v", filename, i), result.Type),
			)
			if _, err := os.Stat(fullpath); os.IsNotExist(err) {
				break
			}
		}
	}

	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Downloaded: %s", fullpath)

	return nil
}

func (c *Client) Download(result *SearchResult, accessToken string) error {

	var downloadDir string
	for i := 0; ; i++ {
		downloadDir = path.Join(c.Config.Dir, fmt.Sprintf("fs_%d", i))
		if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
			os.Mkdir(downloadDir, 0777)
			break
		}
	}
	for _, r := range result.Results {
		err := c.DownloadOne(&r, downloadDir, accessToken)
		if err != nil {
			return err
		}
	}
	return nil

}

func addExtension(name, ftype string) string {
	if !strings.HasSuffix(name, fmt.Sprintf(".%s", ftype)) {
		return fmt.Sprintf("%s.%s", name, ftype)
	}
	return name
}

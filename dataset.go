package domoapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type RequestHandlerService interface {
	Handler(req *http.Request) (*http.Response, error)
}

type RequestHandler struct{}

type DomoAPI struct {
	requestHandlerService RequestHandlerService
}

func NewDomoAPI() *DomoAPI {
	return &DomoAPI{
		requestHandlerService: &RequestHandler{},
	}
}

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	ExpiresAt   time.Time
}

//DomoDataset
type DomoDataset struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Rows        int        `json:"rows,omitempty"`
	Columns     int        `json:"columns,omitempty"`
	Schema      *Schema    `json:"schema,omitempty"`
	Owner       *Owner     `json:"owner,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

//Owner is domo dataset's owner
type Owner struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

//Schema is data's schema
type Schema struct {
	Columns []Column `json:"columns,omitempty"`
}

//Column is schema's column
type Column struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

//GetDataByDatasetID fetch data given domo's datasetID. Use header=true to include header in the response
func (d *DomoAPI) GetDataByDatasetID(token string, datasetID string, header bool) (string, error) {
	includeHeader := ""
	if header {
		includeHeader = "?includeHeader=true"
	}
	apiURL := os.Getenv("DOMO_API_URL") + "/v1/datasets/" + datasetID + "/data" + includeHeader
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "text/csv")
	req.Header.Add("Authorization", "bearer "+token)

	resp, err := d.requestHandlerService.Handler(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Api request failed : %v", string(body))
	}

	return string(body), nil
}

//GetDatasetIDByName get domo datasetID using domo dataset name
func (d *DomoAPI) GetDatasetIDByName(datasetName string, token string) ([]string, error) {

	var datasetIDs []string
	datasets, err := d.ListDatasets(token)
	if err != nil {
		return nil, err
	}
	for _, d := range datasets {
		if d.Name == datasetName {
			datasetIDs = append(datasetIDs, d.ID)
		}
	}
	if len(datasetIDs) > 0 {
		return datasetIDs, nil
	}
	return nil, nil
}

//ListDatasets list all domo datasets in the belonging domo instance
func (d *DomoAPI) ListDatasets(token string) ([]DomoDataset, error) {
	var dataSets []DomoDataset

	var tmpSets []DomoDataset
	limit := 50
	counter := 1

	for counter >= 1 {
		strOffset := fmt.Sprintf("&offset=%d", (counter-1)*limit)
		strLimit := fmt.Sprintf("&limit=%d", limit)
		apiURL := os.Getenv("DOMO_API_URL") + "/v1/datasets?sort=name" + strLimit + strOffset
		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "bearer "+token)

		resp, err := d.requestHandlerService.Handler(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {

			return nil, fmt.Errorf("Domo api resoonseded with erorr: %d", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &tmpSets); err != nil {
			return nil, fmt.Errorf("Cannot parse json : %v", err)
		}
		dataSets = append(dataSets, tmpSets...)
		counter++
		if len(tmpSets) == 0 {
			break
		}
	}

	return dataSets, nil
}

//AddDataToDataset adds data to the given dataset. Use replace=true to reset dataset's data with the given data.
func (d *DomoAPI) AddDataToDataset(datasetID string, data string, replace bool, token string) error {
	if datasetID == "" {
		return fmt.Errorf(" error : issing datasetID")
	}
	method := "APPEND"
	if replace {
		method = "REPLACE"
	}
	apiURL := os.Getenv("DOMO_API_URL") + "/v1/datasets/" + datasetID + "/data?updateMethod=" + method
	req, err := http.NewRequest(http.MethodPut, apiURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/csv")
	req.Header.Add("Authorization", "bearer "+token)

	resp, err := d.requestHandlerService.Handler(req)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Domo api resonseded with erorr: %d \n URL: %s", resp.StatusCode, apiURL)
	}
	return err
}

//CreateDataset create dataset on domo instance
func (d *DomoAPI) CreateDataset(dds DomoDataset, token string) (*DomoDataset, error) {
	apiURL := os.Getenv("DOMO_API_URL") + "/v1/datasets"

	sDataset, err := json.Marshal(dds)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(sDataset))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "bearer "+token)

	resp, err := d.requestHandlerService.Handler(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating dataset: %s", string(sDataset))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var s *DomoDataset

	if err := json.Unmarshal(body, &s); err != nil {
		return nil, fmt.Errorf("Error deserializing dataset schema - %v", err)
	}

	return s, err
}

//CreateAccessToken create domo accessToken using key clientKey and clientSecrete in .env file.
func (d *DomoAPI) CreateAccessToken() (*Token, error) {
	scopes := os.Getenv("DOMO_AUTH_SCOPE")
	if scopes == "" {
		scopes = "data"
	}
	apiURL := os.Getenv("DOMO_API_URL") + "/oauth/token?grant_type=client_credentials&scope=" + scopes
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	clientID := os.Getenv("DOMO_CLIENT_ID")
	clientSecret := os.Getenv("DOMO_CLIENT_SECRET")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := d.requestHandlerService.Handler(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body request - %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected status 2XX getting oauth access_token, got %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	var token *Token
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("Error deserializing access_token - %v", err)
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("error: invalid accesstoken")
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn-1) * time.Second)
	return token, nil
}

//Handler handles http client request.
func (r *RequestHandler) Handler(req *http.Request) (*http.Response, error) {
	var client *http.Client
	proxyURL := os.Getenv("DOMO_PROXY_URL")
	if proxyURL != "" {
		proxy, _ := url.Parse(proxyURL)
		transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
		client = &http.Client{
			Transport: transport,
			Timeout:   3 * time.Minute,
		}
	} else {
		client = &http.Client{
			Timeout: 3 * time.Minute,
		}
	}
	return client.Do(req)
}

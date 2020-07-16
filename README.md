# go-domo-api

[![CircleCI](https://circleci.com/gh/rakutentech/go-domo-api.svg?style=svg)](https://circleci.com/gh/rakutentech/go-domo-api)

This library suppor-ts Domo api endpoints to administer your data and users, giving you the power and flexibility to get the most out of Domo. It is built with `go version go1.14`

Reference: <https://developer.domo.com/docs/api-overview/api-overview>

## Installing

### \*go get

```bash
    go get -u github.com/rakutentech/go-domo-api
```

## Configurations

- This package use golang environmrnt variable as setting. It uses `os.Getenv` to get the configuration values. You can use any enivronment setting package. One of the common package is `godotenv` from `https://github.com/joho/godotenv`.

### General Configs

| No  | Environment Variable | default | Required | Explanation                                                                                                                                       |
| --- | -------------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | DOMO_API_URL         | ""      | Yes      | Domo api url                                                                                                                                      |
| 2   | DOMO_CLIENT_ID       | ""      | Yes      | Domo ClientID                                                                                                                                     |
| 3   | DOMO_CLIENT_SECRET   | ""      | Yes      | Domo Client Secret                                                                                                                                |
| 4   | DOMO_PROXY_URL       | ""      | No       | Proxy URL to access to DOMO from a proxied environment                                                                                            |
| 5   | DOMO_AUTH_SCOPE      | "data"  | No       | Domo Auth token scopes. (data, user, workflow, datasboard, account, audit, buzz) It can be specified with multiple values. Separated by comma(,). |

## Usage

```golang
 //import
 import domoapi "github.com/rakutentech/go-domo-api"

 //Create DomoAPI
d := domoapi.NewDomoAPI()

//Create accessToken
 tk, _ := d.CreateAccessToken()

// Create Domo dataset
dataset := &DomoDataset{
		Name:        "Leonhard Euler Party",
		Description: "Mathematician Guest List",
		Rows:        0,
		Schema: &Schema{
			Columns: []Column{
				{
					Type: "STRING",
					Name: "Friend",
				}, {
					Type: "STRING",
					Name: "Attending",
				},
			},
		},
		Owner: &Owner{
			Name: "DomoSupport",
			ID:   27,
		},
	}
ds, _ := d.CreateDataset(dataset, tk)

//Get DatasetID
dID, _ := d.GetDatasetID("dataset_name", tk)

// Get Data from dataset
data, _ :=d.GetDataByDatasetID(tk, "dataset_id", true)

//List all datasets
datasetList, _ := d.ListDatasets(tk)

```

### Sample Configuration

- Create a `.env` file and add the setting value

```markdown
DOMO_API_URL=https://rakuten-training.domo.com
DOMO_CLIENT_ID=dummy_id
DOMO_CLIENT_SECRET=dummy_secret
DOMO_PROXY_URL=https://proxy_dummy.example.com:8080
```

package domoapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	mocks "github.com/rakutentech/go-domo-api/mocks"
	"github.com/golang/mock/gomock"
)

var (
	tokenAPIRespJSON = `{
		"access_token": "eyJhbGciOiJSUzI1NiJ9.eyJyb2xlIjoiQWRtaW4iLCJzY29wZSI6WyJkYXRhIiwidXNlciJdLCJleHAiOjE1MDE3ODM1MTksImVudiI6InByb2QxIiwidXNlcklkIjo5NjQzODI1OTAsImp0aSI6IjgxZGVjZTRjLWRmNzMtNDU2OS04NTNjLTJkMWEzMjg4OTdmZCIsImNsaWVudF9pZCI6IjQ0MWUzMDdhLWIyYTEtNGE5OS04NTYxLTE3NGU1YjE5M2VlYSIsImN1c3RvbWVyIjoidHJhaW5pbmcifQ.O0H29RFMjruZjA7cXvIctNf8IyuKSyoLggRB5ps7r8LwJL7FM-BWu8oyEPnCopvED-2Sy6emg__8PKQ6r4FPRnsyp4A1uwLwyii8fUmp8NQX3QlprL72_Xc-5ghDtWILPX_Pg77giFLbH-nAys5nsI9S2MDL-xwIFnS3p0iAqSmio5yk6F61Gi_fXVPbDARQ6c2Ci2q2wxFG-xk1lqdlKoPnQwmcQcIJhfgvQ6-15SxNuwZhoQZCZt01wTJ65phFsqoYFHPgznlMxISxDJz2moCjKIY8O9qobj9kIbegpusKUzLeBO9SVc7V3KNZal9f_m-8eGJ3j2bbMFEnaHz8jMJQ",
		"customer": "acmecompany",
		"expires_in": 3599,
		"jti": "81dece4c-df73-4569-853c-2d1a328897fd",
		"role": "Admin",
		"scope": "data user",
		"token_type": "bearer",
		"userId": 964382593
	}`
	createDatasetOKJson = `{
		"id": "4405ff58-1957-45f0-82bd-914d989a3ea3",
		"name": "Leonhard Euler Party",
		"description": "Mathematician Guest List",
		"rows": 0,
		"columns": 0,
		"schema": {
		  "columns": [ {
			"type": "STRING",
			"name": "Friend"
		  }, {
			"type": "STRING",
			"name": "Attending"
		  } ]
		},
		"owner": {
		  "id": 27,
		  "name": "DomoSupport"
		},
		"createdAt": "2016-06-21T17:20:36Z",
		"updatedAt": "2016-06-21T17:20:36Z"
	  }
	`
	listDatasetsJSON = `
	[ {
		"id": "08a061e2-12a2-4646-b4bc-20beddb403e3",
		"name": "Questions regarding Euclid's Fundamental Theorem of Arithmetic",
		"rows": 1,
		"columns": 6,
		"createdAt": "2015-12-10T07:06:14Z",
		"updatedAt": "2016-02-29T20:56:20.567Z"
	  }, {
		"id": "317970a1-6a6e-4f70-8e09-44cf5f34cf44",
		"name": "Ideas Regarding Physics",
		"description": "Notes",
		"rows": 1289280,
		"columns": 9,
		"createdAt": "2013-09-24T20:51:48Z",
		"updatedAt": "2016-02-29T20:56:07.619Z"
	  }, {
		"id": "cc22901d-c856-47c5-89a3-5228a4fa5663",
		"name": "Rene Descartes Mentions",
		"description": "",
		"rows": 194723231,
		"columns": 12,
		"createdAt": "2014-05-01T22:01:17Z",
		"updatedAt": "2016-02-29T20:56:05.034Z"
	  }, {
		"id": "36ea3481-5b90-4181-a4a8-4d9388e85d9e",
		"name": "Symbolic Logic",
		"description": "",
		"rows": 349660,
		"columns": 12,
		"createdAt": "2014-11-26T18:29:09Z",
		"updatedAt": "2016-02-29T20:55:50.337Z"
	  } ]
	`
	csvResponse = `
	first,second,third
1,2,3
4,5,6
7,8,92
	`
	emptyJSON   = ""
	errorJSON   = "{'error': 500}"
	sampleToken = Token{
		AccessToken: "eyJhbGciOiJSUzI1NiJ9.eyJyb2xlIjoiQWRtaW4iLCJzY29wZSI6WyJkYXRhIiwidXNlciJdLCJleHAiOjE1MDE3ODM1MTksImVudiI6InByb2QxIiwidXNlcklkIjo5NjQzODI1OTAsImp0aSI6IjgxZGVjZTRjLWRmNzMtNDU2OS04NTNjLTJkMWEzMjg4OTdmZCIsImNsaWVudF9pZCI6IjQ0MWUzMDdhLWIyYTEtNGE5OS04NTYxLTE3NGU1YjE5M2VlYSIsImN1c3RvbWVyIjoidHJhaW5pbmcifQ.O0H29RFMjruZjA7cXvIctNf8IyuKSyoLggRB5ps7r8LwJL7FM-BWu8oyEPnCopvED-2Sy6emg__8PKQ6r4FPRnsyp4A1uwLwyii8fUmp8NQX3QlprL72_Xc-5ghDtWILPX_Pg77giFLbH-nAys5nsI9S2MDL-xwIFnS3p0iAqSmio5yk6F61Gi_fXVPbDARQ6c2Ci2q2wxFG-xk1lqdlKoPnQwmcQcIJhfgvQ6-15SxNuwZhoQZCZt01wTJ65phFsqoYFHPgznlMxISxDJz2moCjKIY8O9qobj9kIbegpusKUzLeBO9SVc7V3KNZal9f_m-8eGJ3j2bbMFEnaHz8jMJQ",
		ExpiresIn:   3599,
		ExpiresAt:   (time.Now().Add(time.Duration(3599-1) * time.Second)),
	}
)

func getMockResponse(jsonData string, statusCode int) *http.Response {
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonData)))
	return &http.Response{
		StatusCode: statusCode,
		Body:       r,
	}
}
func TestDomoAPI_CreateAccessToken(t *testing.T) {

	tests := []struct {
		name    string
		want    *Token
		wantErr bool
		api     func(ctrl *gomock.Controller) *DomoAPI
	}{
		{
			name:    "success and received token",
			want:    &sampleToken,
			wantErr: false,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(tokenAPIRespJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "empty response",
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(emptyJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
			wantErr: true,
		},
		{
			name: "500 response",
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(errorJSON, 500), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)

			got, err := domoAPI.CreateAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.CreateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.AccessToken, tt.want.AccessToken) {
					t.Errorf("Domo.CreateAccessToken() AccessToken = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got.ExpiresIn, tt.want.ExpiresIn) {
					t.Errorf("Domo.CreateAccessToken() ExpiresIn = %v, want %v", got.ExpiresIn, tt.want.ExpiresIn)
				}
			}
		})
	}
}

func TestDomoAPI_CreateDataSet(t *testing.T) {

	type args struct {
		dds   DomoDataset
		token Token
	}
	tests := []struct {
		name    string
		api     func(ctrl *gomock.Controller) *DomoAPI
		args    args
		want    *DomoDataset
		wantErr bool
	}{
		{
			name: "success and return dataset",
			args: args{
				token: sampleToken,
			},
			wantErr: false,
			want: &DomoDataset{
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
			},
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(createDatasetOKJson, 201), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "error and return nil",
			args: args{
				token: sampleToken,
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(errorJSON, 500), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "domo api return empty",
			args: args{
				token: sampleToken,
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(emptyJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)

			got, err := domoAPI.CreateDataset(tt.args.dds, tt.args.token.AccessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.CreateDataSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.Name, tt.want.Name) {
					t.Errorf("DomoAPI.CreateDataSet() = %v, want %v", got.Name, tt.want.Name)
				}
				if !reflect.DeepEqual(got.Description, tt.want.Description) {
					t.Errorf("DomoAPI.CreateDataSet() = %v, want %v", got.Description, tt.want.Description)
				}
				if !reflect.DeepEqual(got.Columns, tt.want.Columns) {
					t.Errorf("DomoAPI.CreateDataSet() = %v, want %v", got.Columns, tt.want.Columns)
				}
				if !reflect.DeepEqual(got.Owner, tt.want.Owner) {
					t.Errorf("DomoAPI.CreateDataSet() = %v, want %v", got.Owner, tt.want.Owner)
				}
			}

		})
	}
}

func TestDomoAPI_AddDataToDataset(t *testing.T) {
	type args struct {
		datasetID string
		data      string
		replace   bool
		token     Token
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		api     func(ctrl *gomock.Controller) *DomoAPI
	}{
		{
			name: "error and return nil",
			args: args{
				token:     sampleToken,
				datasetID: "ds_id001",
				replace:   false,
				data:      "1,1,1,1",
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(errorJSON, 500), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "wrong status code",
			args: args{
				token:     sampleToken,
				datasetID: "ds_id001",
				replace:   false,
				data:      "1,1,1,1",
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(emptyJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "missing dataset id",
			args: args{
				token: sampleToken,
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)

				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)

			if err := domoAPI.AddDataToDataset(tt.args.datasetID, tt.args.data, tt.args.replace, tt.args.token.AccessToken); (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.AddDataToDataset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDomoAPI_ListDatasets(t *testing.T) {
	var ld []DomoDataset
	_ = json.Unmarshal([]byte(listDatasetsJSON), &ld)

	type args struct {
		token Token
	}
	tests := []struct {
		name    string
		args    args
		want    []DomoDataset
		wantErr bool
		api     func(ctrl *gomock.Controller) *DomoAPI
	}{
		{
			name: "success and get list of datasets",
			args: args{
				token: sampleToken,
			},
			want: ld,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(listDatasetsJSON, 200), nil)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse("[]", 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "500 response",
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(errorJSON, 500), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
			wantErr: true,
		},
		{
			name: "domo api return empty",
			args: args{
				token: sampleToken,
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(emptyJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)
			got, err := domoAPI.ListDatasets(tt.args.token.AccessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.ListDatasets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomoAPI.ListDatasets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomoAPI_GetDatasetIDByName(t *testing.T) {
	var ld []DomoDataset
	_ = json.Unmarshal([]byte(listDatasetsJSON), &ld)

	type args struct {
		datasetName string
		token       Token
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
		api     func(ctrl *gomock.Controller) *DomoAPI
	}{
		{
			name: "success and get dataset's key",
			args: args{
				token:       sampleToken,
				datasetName: "Rene Descartes Mentions",
			},
			want: []string{"cc22901d-c856-47c5-89a3-5228a4fa5663"},
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(listDatasetsJSON, 200), nil)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse("[]", 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "success ,but key not found",
			args: args{
				token:       sampleToken,
				datasetName: "Not found dataset",
			},
			want: nil,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(listDatasetsJSON, 200), nil)
				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse("[]", 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "500 response",
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)

				r := ioutil.NopCloser(bytes.NewReader([]byte(errorJSON)))
				mresp := &http.Response{
					StatusCode: 500,
					Body:       r,
				}
				rmock.EXPECT().Handler(gomock.Any()).Return(mresp, nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)
			got, err := domoAPI.GetDatasetIDByName(tt.args.datasetName, tt.args.token.AccessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.GetDatasetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomoAPI.GetDatasetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomoAPI_GetDataByDatasetID(t *testing.T) {
	type args struct {
		token     Token
		datasetID string
		header    bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		api     func(ctrl *gomock.Controller) *DomoAPI
	}{
		{
			name: "success and get dataset's key",
			args: args{
				token:     sampleToken,
				datasetID: "dummy_id",
			},
			want: csvResponse,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)

				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(csvResponse, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "success but get empty response",
			args: args{
				token:     sampleToken,
				datasetID: "dummy_id",
			},
			want: "",
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)

				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(emptyJSON, 200), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
		{
			name: "error occured",
			args: args{
				token:     sampleToken,
				datasetID: "dummy_id",
			},
			wantErr: true,
			api: func(ctrl *gomock.Controller) *DomoAPI {
				rmock := mocks.NewMockRequestHandlerService(ctrl)

				rmock.EXPECT().Handler(gomock.Any()).Return(getMockResponse(errorJSON, 500), nil)
				return &DomoAPI{
					requestHandlerService: rmock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			domoAPI := tt.api(ctrl)
			got, err := domoAPI.GetDataByDatasetID(tt.args.token.AccessToken, tt.args.datasetID, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("DomoAPI.GetDataByDatasetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomoAPI.GetDataByDatasetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

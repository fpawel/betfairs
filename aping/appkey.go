package aping

import (
	"encoding/json"
	"fmt"

	"github.com/nu7hatch/gouuid"
)

type developerApp struct {
	AppName     string                `json:"appName"`
	AppId       int                   `json:"appId"`
	AppVersions []developerAppVersion `json:"appVersions"`
}

type developerAppVersion struct {
	Owner                string `json:"owner"`
	VersionId            int    `json:"versionId"`
	Version              string `json:"version"`
	ApplicationKey       string `json:"applicationKey"`
	DelayData            bool   `json:"delayData"`
	SubscriptionRequired bool   `json:"subscriptionRequired"`
	OwnerManaged         bool   `json:"ownerManaged"`
	Active               bool   `json:"active"`
	VendorId             string `json:"vendorId,omitempty"`
	VendorSecret         string `json:"vendorSecret,omitempty"`
}


func fetchAppKey(sessionToken string) (appKey string, err error) {

	responseBody, err := AccauntAPIEndpoint("getDeveloperAppKeys").getResponse(sessionToken, nil, nil)
	if err != nil {
		return
	}


	var tmp struct {
		Jsonrpc string         `json:"appName"`
		Result  []developerApp `json:"result"`
	}

	if json.Unmarshal(responseBody, &tmp) == nil && len(tmp.Result) > 0 && len(tmp.Result[0].AppVersions) > 0 {
		appKey = tmp.Result[0].AppVersions[0].ApplicationKey
		return
	}


	var u4 *uuid.UUID
	if u4, err = uuid.NewV4(); err != nil {
		return
	}


	responseBody, err = AccauntAPIEndpoint("createDeveloperAppKeys").getResponse(sessionToken, nil, struct {
		AppName string `json:"appName"`
	}{u4.String()})

	if err != nil {
		return
	}

	var x struct {
		Jsonrpc string       `json:"appName"`
		Result  developerApp `json:"result"`
	}
	if json.Unmarshal(responseBody, &x) == nil && len(x.Result.AppVersions) > 0 {
		appKey = x.Result.AppVersions[0].ApplicationKey
	} else {
		err = fmt.Errorf("createDeveloperAppKeys: required fields missing: %v", string(responseBody))
	}

	return
}



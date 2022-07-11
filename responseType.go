package main

import "applemint-go/crud"

type MessageModel struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type GroupInfoResponse struct {
	TotalCount int64            `json:"totalCount"`
	GroupInfos []crud.GroupInfo `json:"groupInfos"`
}

type AnalyzeImgurResponse struct {
	Images []string `json:"images"`
}

type DropboxResponse struct {
	AsyncJobId string `json:"asyncJobId"`
}

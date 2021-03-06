package crud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"go.mongodb.org/mongo-driver/bson"
)

const raindropEndPoint = "https://api.raindrop.io"
const collectionAPI = "/rest/v1/collections"

//const raindropAPI = "/rest/v1/raindrop"

func SendToDropbox(path string, url string) (string, error) {
	// connect to dropbox
	access_token := os.Getenv("ENV_DROPBOX_ACCESS_TOKEN")
	if access_token == "" {
		return "", errors.New("env_dropbox_access_token not set")
	}

	config := dropbox.Config{
		Token:    access_token,
		LogLevel: dropbox.LogInfo,
	}
	dbx := files.New(config)
	arg := files.NewSaveUrlArg(path, url)
	result, err := dbx.SaveUrl(arg)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(result)

	return result.AsyncJobId, err
}

func GetCollectionFromRaindrop() ([]map[string]interface{}, error) {
	// connect to raindrop
	log.Print("GetCollectionFromRaindrop")
	req, err := http.NewRequest("GET", raindropEndPoint+collectionAPI, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ENV_RAINDROP_ACCESS_TOKEN"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawData map[string]interface{}
	log.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	collections := make([]map[string]interface{}, 0, 0)
	items := rawData["items"].([]interface{})
	for _, item := range items {
		itemRawMap := item.(map[string]interface{})
		idString := fmt.Sprintf("%.0f", itemRawMap["_id"])

		resultItem := make(map[string]interface{})
		resultItem["id"] = idString
		resultItem["title"] = itemRawMap["title"]
		collections = append(collections, resultItem)
	}
	return collections, nil
}

func SendToRaindrop(item Item, collection string) ([]byte, error) {
	// connect to raindrop
	log.Print("SendToRaindrop")
	jsonData := bson.M{}
	jsonData["title"] = item.TextContent
	jsonData["link"] = item.Url
	jsonData["tags"] = item.Tags
	collectionJson := bson.M{}
	collectionJson["$id"] = collection
	jsonData["collection"] = collectionJson

	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.raindrop.io/rest/v1/raindrop", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ENV_RAINDROP_ACCESS_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("response Body:", string(body))

	return body, nil
}

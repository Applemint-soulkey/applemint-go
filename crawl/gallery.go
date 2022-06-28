package crawl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const imgurAlbumApi = "https://api.imgur.com/3/album/"
const imgurImageApi = "https://api.imgur.com/3/image/"
const imgurGalleryApi = "https://api.imgur.com/3/gallery/album/"

func HandleImgurLink(link string) ([]string, error) {
	imgurAuthHeader := "Client-ID " + os.Getenv("ENV_IMGUR_CLIENT_ID")
	albumKeyword := "/a/"
	galleryKeyword := "/gallery/"
	imgurHash := link[strings.LastIndex(link, "/")+1:]
	log.Println("HandleImgurLink:", link)
	log.Println("imgurHash:", imgurHash)
	log.Println("imgurAuthHeader:", imgurAuthHeader)
	var images []string
	var err error

	if strings.Contains(link, albumKeyword) {
		// Album Process
		images, err = getImagesFromAlbum(imgurHash, imgurAuthHeader)
		if err != nil {
			log.Println("HandleImgurLink: GetImagesFromAlbum error:", err)
			return nil, err
		}
		log.Println("HandleImgurLink: images:", images)
	} else if strings.Contains(link, galleryKeyword) {
		// Gallery Process
		images, err = getImagesFromGallery(imgurHash, imgurAuthHeader)
		if err != nil {
			log.Println("HandleImgurLink: GetAlbumsFromGallery error:", err)
			return nil, err
		}
		log.Println("HandleImgurLink: albums:", images)
	} else {
		// Default Process
		images, err = getImageByHash(imgurHash, imgurAuthHeader)
		if err != nil {
			log.Println("HandleImgurLink: GetImageByHash error:", err)
			return nil, err
		}
	}
	return images, nil
}

func getImagesFromGallery(imgurHash string, imurAuthHeader string) ([]string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", imgurGalleryApi+imgurHash, nil)
	if err != nil {
		log.Println("GetImagesFromGallery: NewRequest error:", err)
		return nil, err
	}
	req.Header.Set("Authorization", imurAuthHeader)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("GetImagesFromGallery: Do error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("GetImagesFromGallery: Gallery response:", resp.StatusCode)
	if resp.StatusCode == 200 {
		// Gallery OK
		log.Println("GetImagesFromGallery: Gallery OK")
		var rawData map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &rawData)
		if err != nil {
			log.Println("GetImagesFromGallery: Unmarshal error:", err)
			return nil, err
		}
		items := rawData["data"].(map[string]interface{})["images"].([]interface{})
		var images []string
		for _, item := range items {
			itemRawMap := item.(map[string]interface{})
			log.Println(itemRawMap["link"].(string))
			images = append(images, itemRawMap["link"].(string))
		}

		return images, err
	} else {
		// Gallery Error
		log.Println("GetAlbumsFromGallery: Gallery Error")
		return nil, err
	}
}

func getImagesFromAlbum(imgurHash string, imgurAuthHeader string) ([]string, error) {
	var images []string

	client := &http.Client{}
	req, err := http.NewRequest("GET", imgurAlbumApi+imgurHash+"/images", nil)
	if err != nil {
		log.Println("HandleImgurLink: NewRequest error:", err)
		return nil, err
	}
	req.Header.Set("Authorization", imgurAuthHeader)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HandleImgurLink: Do error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("HandleImgurLink: Album response:", resp.StatusCode)

	if resp.StatusCode == 200 {
		// Album OK
		log.Println("HandleImgurLink: Album OK")
		var rawData map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &rawData)
		if err != nil {
			log.Println("HandleImgurLink: Unmarshal error:", err)
			return nil, err
		}
		items := rawData["data"].([]interface{})
		for _, item := range items {
			itemRawMap := item.(map[string]interface{})
			log.Println(itemRawMap["link"].(string))
			images = append(images, itemRawMap["link"].(string))
		}
		return images, nil
	} else {
		// Album Error
		log.Println("HandleImgurLink: Album Error")
		return nil, err
	}
}

func getImageByHash(imgurHash string, imgurAuthHeader string) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", imgurImageApi+imgurHash, nil)
	if err != nil {
		log.Println("HandleImgurLink: NewRequest error:", err)
		return nil, err
	}
	req.Header.Set("Authorization", imgurAuthHeader)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("HandleImgurLink: Do error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("HandleImgurLink: Image response:", resp.StatusCode)
	if resp.StatusCode == 200 {
		// Image OK
		log.Println("HandleImgurLink: Image OK")
		var rawData map[string]interface{}
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &rawData)
		if err != nil {
			log.Println("HandleImgurLink: Unmarshal error:", err)
			return nil, err
		}
		log.Println("HandleImgurLink: Image link:", rawData["data"].(map[string]interface{})["link"].(string))
		return []string{rawData["data"].(map[string]interface{})["link"].(string)}, nil
	} else {
		// Image Error
		log.Println("HandleImgurLink: Image Error")
		return nil, err
	}
}

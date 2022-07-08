package main

import (
	"applemint-go/crawl"
	"applemint-go/crud"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageModel struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// moveItemRequestHandler godoc
// @Summary Move Item
// @Description Move Item via Collections
// @name Move Item
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param target query string true "Target Collection"
// @Param origin query string true "Origin Collection"
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /item/move/{id} [get]
func moveItemRequestHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	target := ctx.Query("target")
	origin := ctx.Query("origin")
	log.Printf("move item %s from %s to %s", id, origin, target)
	if target == "" || origin == "" {
		ctx.JSON(http.StatusBadRequest, MessageModel{Type: "error", Message: "missing target or origin"})
		return
	}
	err := crud.MoveItem(id, origin, target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: "item " + id + " moved from " + origin + " to " + target})
}

// getItemRequestHandler godoc
// @Summary Get Item
// @Description Get Item
// @name Get Item
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param collection path string true "Collection"
// @Success 200 {object} crud.Item
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /item/{id}/{collection} [get]
func getItemRequestHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	collection := ctx.Param("collection")
	item, err := crud.GetItem(id, collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

//updateItemRequestHandler godoc
// @Summary Update Item
// @Description Update Item
// @name Update Item
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param collection path string true "Collection"
// @Param item body crud.Item true "Item"
// @Success 200 {object} crud.Item
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /item/{id}/{collection} [put]
func updateItemReqeustHandler(ctx *gin.Context) {
	collection := ctx.Param("collection")
	id := ctx.Param("id")
	item := crud.Item{}
	err := ctx.BindJSON(&item)
	if collection == "" || id == "" {
		ctx.JSON(http.StatusBadRequest, MessageModel{Type: "error", Message: "missing collection or id"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	updateCnt := crud.UpdateItem(id, collection, item)
	if updateCnt > 0 {
		ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: "item " + id + " updated in " + collection})
	} else {
		ctx.JSON(http.StatusOK, MessageModel{Type: "error", Message: "item " + id + " not found in " + collection})
	}
}

// deleteItemRequestHandler godoc
// @Summary Delete Item
// @Description Delete Item
// @name Delete Item
// @Tags Item
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param collection path string true "Collection"
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /item/{id}/{collection} [delete]
func deleteItemRequestHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	collection := ctx.Param("collection")
	deleteCnt := crud.DeleteItem(id, collection)
	log.Printf("Deleted %d items from collection %s", deleteCnt, collection)
	if deleteCnt > 0 {
		ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: "item " + id + " deleted from " + collection})
	} else {
		ctx.JSON(http.StatusOK, MessageModel{Type: "error", Message: "item " + id + " not found in " + collection})
	}
}

// getBookmarkListRequestHandler godoc
// @Summary Get Bookmark List
// @Description Get Bookmark List
// @name Get Bookmark List
// @Tags Bookmark
// @Accept json
// @Produce json
// @Success 200 {object} []crud.BookmarkInfo
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /bookmark [get]
func getBookmarkListRequestHandler(ctx *gin.Context) {
	log.Print("get bookmark list")
	bookmarkList, err := crud.GetBookmarkList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, bookmarkList)
}

func makeBookmarkRequestHander(ctx *gin.Context) {
	log.Print("make bookmark")
	item := crud.Item{}
	origin := ctx.Query("from")
	path := ctx.Query("path")
	err := ctx.BindJSON(&item)
	if origin == "" || path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing origin or path"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := crud.SendToBookmark(item, origin, path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func getItemListRequestHandler(ctx *gin.Context) {
	collection := ctx.Param("target")
	cursor, err := strconv.Atoi(ctx.Query("cursor"))
	if err != nil {
		cursor = 0
	}
	domain := ctx.Query("domain")
	path := ctx.Query("path")

	items, err := crud.GetItems(collection, int64(cursor), domain, path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func getCollectionListRequestHandler(ctx *gin.Context) {
	collections, err := crud.GetCollectionList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, collections)
}

func cleanOldItemsRequestHandler(ctx *gin.Context) {
	crawl.ClearOldData("trash")
	crawl.ClearOldData("image-queue")
	ctx.JSON(http.StatusOK, gin.H{"message": "old time cleaner launched"})
}

func getCollectionInfoRequestHandler(ctx *gin.Context) {
	collection := ctx.Param("target")
	totalCount, GroupInfos, err := crud.GetCollectionInfo(collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"totalCount": totalCount,
		"groupInfos": GroupInfos,
	})
}

func clearCollectionRequestHandler(ctx *gin.Context) {
	collection := ctx.Param("target")
	delCount := crud.ClearCollection(collection)
	ctx.JSON(http.StatusOK, gin.H{"deleted": delCount})
}

func crawlRequestHandler(ctx *gin.Context) {
	target := ctx.Param("target")
	resultCount := 0
	switch target {
	case "bp":
		resultCount = crawl.CrawlBP()
	case "isg":
		resultCount = crawl.CrawlISG()
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "crawl launched", "resultCount": resultCount})
}

func getGalleryItemsRequestHandler(ctx *gin.Context) {
	cursor, err := strconv.Atoi(ctx.Query("cursor"))
	if err != nil {
		cursor = 0
	}
	items, err := crud.GetGalleryItems(int64(cursor))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func analyzeImgurRequestHandler(ctx *gin.Context) {
	imgurLink := ctx.Query("link")
	if imgurLink == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing link"})
		return
	}
	log.Printf("analyze imgur to get image info: %s", imgurLink)
	imageInfoList, err := crawl.HandleImgurLink(imgurLink)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, imageInfoList)
}

func dropboxRequestHandler(ctx *gin.Context) {
	path := ctx.Query("path")
	url := ctx.Query("url")
	if path == "" || url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing path or url"})
		return
	}
	log.Printf("dropbox request: %s", path)
	asyncJobId, err := crud.SendToDropbox(path, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"asyncJobId": asyncJobId})
}

func getRaindropCollectionListRequestHandler(ctx *gin.Context) {
	collections, err := crud.GetCollectionFromRaindrop()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, collections)
}

func sendToRaindropRequestHandler(ctx *gin.Context) {
	collectionId := ctx.Query("collectionId")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing collectionId"})
		return
	}
	log.Printf("send to raindrop: %s", collectionId)
	item := crud.Item{}
	err := ctx.BindJSON(&item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	raindropResp, err := crud.SendToRaindrop(item, collectionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, raindropResp)
}

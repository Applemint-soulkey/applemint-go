package main

import (
	"applemint-go/crawl"
	"applemint-go/crud"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
// @Router /item/{collection}/{id} [get]
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
// @Router /item/{collection}/{id} [put]
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
// @Router /item/{collection}/{id} [delete]
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
// @Router /item/bookmark/list [get]
func getBookmarkListRequestHandler(ctx *gin.Context) {
	log.Print("get bookmark list")
	bookmarkList, err := crud.GetBookmarkList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, bookmarkList)
}

// makeBookmarkRequestHandler godoc
// @Summary Make Bookmark
// @Description Make Bookmark
// @name Make Bookmark
// @Tags Bookmark
// @Accept json
// @Produce json
// @Param Item body crud.Item true "Bookmark"
// @Param from query string true "Origin Collection of Item"
// @Param path query string true "Bookmark Path"
// @Success 200 {object} crud.Item
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /item/bookmark [post]
func makeBookmarkRequestHander(ctx *gin.Context) {
	log.Print("make bookmark")
	item := crud.Item{}
	origin := ctx.Query("from")
	path := ctx.Query("path")
	err := ctx.BindJSON(&item)
	if origin == "" || path == "" {
		ctx.JSON(http.StatusBadRequest, MessageModel{Type: "error", Message: "missing origin or path"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	result, err := crud.SendToBookmark(item, origin, path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// getItemListRequestHandler godoc
// @Summary Get Item List
// @Description Get Item List
// @name Get Item List
// @Tags Item
// @Accept json
// @Produce json
// @Param target path string true "Collection"
// @Param cursor query number false "Cursor"
// @Param domain query string false "Domain"
// @Param path query string false "Path"
// @Success 200 {object} []crud.Item
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /collection/{target} [get]
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
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// getCollectionListRequestHandler godoc
// @Summary Get Collection List
// @Description Get Collection List
// @name Get Collection List
// @Tags Collection
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /collection/list [get]
func getCollectionListRequestHandler(ctx *gin.Context) {
	collections, err := crud.GetCollectionList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, collections)
}

// cleanOldItemsRequestHandler godoc
// @Summary Clean Old Items
// @Description Clean Old Items
// @name Clean Old Items
// @Tags Collection
// @Accept json
// @Produce json
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
func cleanOldItemsRequestHandler(ctx *gin.Context) {
	crawl.ClearOldData("trash")
	crawl.ClearOldData("image-queue")
	ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: "old items cleaned"})
}

// getCollectionInfoRequestHandler godoc
// @Summary Get Collection Info
// @Description Get Collection Info
// @name Get Collection Info
// @Tags Collection
// @Accept json
// @Produce json
// @Param target path string true "Collection"
// @Success 200 {object} GroupInfoResponse
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /collection/{target}/info [get]
func getCollectionInfoRequestHandler(ctx *gin.Context) {
	collection := ctx.Param("target")
	totalCount, GroupInfos, err := crud.GetCollectionInfo(collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, GroupInfoResponse{TotalCount: totalCount, GroupInfos: GroupInfos})
}

// clearCollectionRequestHandler godoc
// @Summary Clear Collection
// @Description Clear Collection
// @name Clear Collection
// @Tags Collection
// @Accept json
// @Produce json
// @Param target path string true "Collection"
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /collection/{target}/clear [delete]
func clearCollectionRequestHandler(ctx *gin.Context) {
	collection := ctx.Param("target")
	delCount := crud.ClearCollection(collection)
	ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: fmt.Sprintf("%d items deleted", delCount)})
}

// crawlRequestHandler godoc
// @Summary Crawl Collection
// @Description Crawl Target Site
// @name Crawl Collection
// @Tags Common
// @Accept json
// @Produce json
// @Param target path string true "Collection"
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /crawl/{target} [get]
func crawlRequestHandler(ctx *gin.Context) {
	target := ctx.Param("target")
	resultCount := 0
	switch target {
	case "bp":
		resultCount = crawl.CrawlBP()
	case "isg":
		resultCount = crawl.CrawlISG()
	}
	ctx.JSON(http.StatusOK, MessageModel{Type: "success", Message: fmt.Sprintf("%d items crawled", resultCount)})
}

// getGalleryRequestHandler godoc
// @Summary Get Gallery Items
// @Description Get Gallery Items
// @name Get Gallery Items
// @Tags Gallery
// @Accept json
// @Produce json
// @Param cursor query number false "Cursor"
// @Success 200 {object} []crud.GalleryResponse
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /gallery [get]
func getGalleryItemsRequestHandler(ctx *gin.Context) {
	cursor, err := strconv.Atoi(ctx.Query("cursor"))
	if err != nil {
		cursor = 0
	}
	itemData, err := crud.GetGalleryItems(int64(cursor))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, itemData)
}

// analyzeRequestHandler godoc
// @Summary Analyze Gallery Item
// @Description Analyze Gallery Item
// @name Analyze Gallery Item
// @Tags Gallery
// @Accept json
// @Produce json
// @Param link query string true "Imgur Link"
// @Success 200 {object} AnalyzeImgurResponse
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /gallery/imgur [get]
func analyzeImgurRequestHandler(ctx *gin.Context) {
	imgurLink := ctx.Query("link")
	if imgurLink == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing link"})
		return
	}
	log.Printf("analyze imgur to get image info: %s", imgurLink)
	imageInfoList, err := crawl.HandleImgurLink(imgurLink)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, AnalyzeImgurResponse{Images: imageInfoList})
}

// dropboxRequestHandler godoc
// @Summary Send to Dropbox
// @Description Send to Dropbox
// @name Send to Dropbox
// @Tags External
// @Accept json
// @Produce json
// @Param path query string true "Path"
// @Param url query string true "URL"
// @Success 200 {object} MessageModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /dropbox [get]
func dropboxRequestHandler(ctx *gin.Context) {
	path := ctx.Query("path")
	url := ctx.Query("url")
	if path == "" || url == "" {
		ctx.JSON(http.StatusBadRequest, MessageModel{Type: "error", Message: "missing path or url"})
		return
	}
	log.Printf("dropbox request: %s", path)
	asyncJobId, err := crud.SendToDropbox(path, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, DropboxResponse{AsyncJobId: asyncJobId})
}

// getRaindropCollectionListRequestHandler godoc
// @Summary Get Raindrop Collection List
// @Description Get Raindrop Collection List
// @name Get Raindrop Collection List
// @Tags External
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /raindrop/list [get]
func getRaindropCollectionListRequestHandler(ctx *gin.Context) {
	collections, err := crud.GetCollectionFromRaindrop()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, collections)
}

// sendToRaindropRequestHandler godoc
// @Summary Send to Raindrop
// @Description Send to Raindrop
// @name Send to Raindrop
// @Tags External
// @Accept json
// @Produce json
// @Param collectionId path string true "Raindrop Collection ID"
// @Param Item body crud.Item true "Raindrop Request"
// @Success 200 string CustomModel
// @Failure 400 {object} MessageModel
// @Failure 500 {object} MessageModel
// @Router /raindrop/{collectionId} [put]
func sendToRaindropRequestHandler(ctx *gin.Context) {
	collectionId := ctx.Param("collectionId")
	if collectionId == "" {
		ctx.JSON(http.StatusBadRequest, MessageModel{Type: "error", Message: "missing collectionId"})
		return
	}
	log.Printf("send to raindrop: %s", collectionId)
	item := crud.Item{}
	err := ctx.BindJSON(&item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}

	raindropResp, err := crud.SendToRaindrop(item, collectionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, MessageModel{Type: "error", Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, raindropResp)
}

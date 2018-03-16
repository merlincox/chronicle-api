package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"encoding/json"
	"fmt"
	"errors"
	"projects/chronicle-api/models"
    "projects/chronicle-api/utils"
)

const AWS = true

type RawVideo struct {
	Title       string `json:"title"`
	Summary string `json:"summary,omitempty"`
	Id string `json:"id"`
	HoldingImageURL string `json:"holdingImageURL"`
	Topic string `json:"topic"`
	Items []models.MediaItem `json:"items"`
}

var (
     DummyVideosMap map[string]models.Video
     DummyVideos []models.Video
     DummyHeros []models.Hero
     DummyPlaylistsMap map[string]models.Playlist
     DummyPlaylists []models.Playlist
     DummyHeroCollection models.HeroCollection
     DummyVideoCollection models.VideoCollection
     DummyPlaylistCollection models.PlaylistCollection
)

func getDummyVideos() []byte {

	dummyData := []byte(`[{
  "title": "Footage of first polar bear cub born in UK in 25 years",
  "id": "p061b9c4",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p061bbfx.jpg",
  "guidance": "",
  "embedRights": "allowed",
  "summary": "The polar bear cub is described as a \"confident and curious\" character",
  "liveRewind": false,
  "simulcast": false,
  "items": [
    {
      "vpid": "p061b9c8",
      "live": false,
      "duration": 29,
      "kind": "programme"
    }
  ]
},{
  "embedRights": "allowed",
  "id": "p05ss2r1",
  "summary": "A female polar bear at a Scottish animal park has given birth to a cub, says the Royal Zoological Society of Scotland.",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ss338.jpg",
  "items": [
    {
      "vpid": "p05ss2r3",
      "live": false,
      "kind": "programme"
    }
  ],
  "title": "UK's first polar bear cub in 25 years born in Highlands",
  "simulcast": false
},{
  "embedRights": "blocked",
  "id": "p05ncdqc",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ncg65.jpg",
  "items": [
    {
      "vpid": "p05ncdqf",
      "live": false,
      "kind": "programme"
    }
  ],
  "title": "Doug Allan: A life capturing the natural world on camera",
  "summary":"Doug Allan is one of the world's best nature cameramen and has filmed some of the most memorable scenes ever broadcast, with some close scrapes with animals along the way.",
  "simulcast": false
},{
  "embedRights": "blocked",
  "id": "p059scjv",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p059sdxl.jpg",
  "items": [
    {
      "vpid": "p059sdm4",
      "live": false,
      "kind": "programme"
    }
  ],
  "title": "The snow might not last long with July temperatures reaching 25 degrees",
  "summary": "Christmas has come early for Lapland zoo polar bears with snow in July.",
  "simulcast": false
}]`)
	return dummyData
}

func getDummyPlaylists() []byte {

	dummyData := []byte(`{

	"uri" : "/playlist/abcde",
	"title" : "Sample playlist",
	"fileid" : "abcde"
	"videos": []
	}`)

	return dummyData
}

func getRawVideos() ([]RawVideo, error) {

	var rva []RawVideo

	err := json.Unmarshal(getDummyVideos(), &rva)

	return rva, err
}

func makeResponse(data interface{}, statusCode int, err error) events.APIGatewayProxyResponse {

        var body string

	if err != nil {

		body = utils.JsonStringify(models.ErrorMessage{Message: fmt.Sprintf("Error: %v", err)})

		// if statusCode is 200 but there is an unhandled error
		// that is a system failure

		if statusCode == 200 {
			statusCode = 500
		}
	}

	if statusCode == 200 {

		body = utils.JsonStringify(data)

	}

	return events.APIGatewayProxyResponse{Body:body, StatusCode: statusCode}
}

func outerHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	handler := getHandler(request)

	data, statusCode, err := handler(request)

	return makeResponse(data, statusCode, err), nil
}

func getHandler(request events.APIGatewayProxyRequest) InnerHandler {

	route := request.RequestContext.HTTPMethod + request.RequestContext.ResourcePath

	switch route {

    case "GET/heros":
        return heroHandler

    case "GET/video/{fileid}":
		return mostWatchedHandler

	case "GET/most-watched/{period}":
		return mostWatchedHandler

	case "GET/playlists":
		return playlistsHandler

	}

	return unknownRouteHandler
}

type InnerHandler func(events.APIGatewayProxyRequest) (interface{}, int, error)

func unknownRouteHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

	var data interface{}
	return data, 404, errors.New("Unknown route")
}

func getPlaylists() (models.PlaylistCollection, error) {

    return DummyPlaylistCollection, error(nil)
}

func getVideos() (models.VideoCollection, error) {

    return DummyVideoCollection, error(nil)
}

func getHeros() (models.HeroCollection, error) {

    return DummyHeroCollection, error(nil)
}

func heroHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

	var data interface{}
	code := 200

	data, err := getHeros()

    if err == nil {
        return data, 200, nil
    }

	return data, code, err
}

func mostWatchedHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

	var data interface{}
	var err error

	if models.ValidPeriods[request.PathParameters["period"]] {

		data, err = getVideos()

		if err == nil {
			return data, 200, nil
		}

		return data, 500, err
	}

	return data, 404, errors.New("Undefined period")
}

func playlistsHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

	var data interface{}
	code := 200

	data, err := getPlaylists()

	if err != nil {
		code = 500
	}

	return data, code, err
}

func populateFixtures() {

    rva, _ := getRawVideos()
    DummyVideosMap = make(map[string]models.Video, len(rva))

	for _, rv := range rva {

        DummyVideosMap[rv.Id] = models.Video{
            Smp: &models.SMP{
                Items: rv.Items,
                Title: rv.Title,
                Summary: rv.Summary,
                HoldingImageURL: rv.HoldingImageURL,
            },
            Id: rv.Id,
            Topic: "Earth",
            LinkUri: "videos/" + rv.Id + "/" + utils.Slug(rv.Title),
        }

        DummyVideos = append(DummyVideos, DummyVideosMap[rv.Id])
        DummyHeros = append(DummyHeros, models.Hero{
            Smp: DummyVideosMap[rv.Id].Smp,
            LinkUri: DummyVideosMap[rv.Id].LinkUri,
            Topic: DummyVideosMap[rv.Id].Topic,

        })
    }

    DummyHeroCollection = models.HeroCollection{
        Items: DummyHeros,
        Total: len(DummyHeros),
        Offset: 0,
    }

}

func main() {


    populateFixtures()

    if AWS {
		lambda.Start(outerHandler)
	} else {
        fmt.Println(utils.JsonStringify(DummyVideos))
        fmt.Println(utils.JsonStringify(DummyHeros))

	}
}


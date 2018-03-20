package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"fmt"
	"projects/chronicle-api/models"
    "projects/chronicle-api/utils"
    "projects/chronicle-api/db"
    "projects/chronicle-api/handlers"
    "os"
)

func frontHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    handler := getInnerHandler(request)

    data, statusCode, err := handler(request)

    return buildResponse(data, statusCode, err), nil
}

func getInnerHandler(request events.APIGatewayProxyRequest) handlers.InnerHandler {

	route := request.RequestContext.HTTPMethod + request.RequestContext.ResourcePath

	switch route {

    case "GET/heros":
        return handlers.HerosHandler

    case "GET/video/{id}":
		return handlers.VideoHandler

    case "GET/playlist/{id}":
        return handlers.PlaylistHandler

    case "GET/most-watched/{period}":
		return handlers.MostWatchedHandler

	case "GET/playlists":
		return handlers.PlaylistsHandler

    case "GET/videos":
        return handlers.VideosHandler

	}

	return handlers.UnknownRouteHandler
}

func buildResponse(data interface{}, statusCode int, err error) events.APIGatewayProxyResponse {

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

func main() {

    db.Init()

    if os.Getenv("AWS") == "true" {
        lambda.Start(frontHandler)
    } else {

        videos, _ := db.GetVideos(0, 0)
        heros, _ := db.GetHeros(0, 0)
        playlists, _ := db.GetPlaylists(0, 0)
        playlist, _ := db.GetPlaylist("abcdef")

        fmt.Println(utils.JsonStringify(videos))
        fmt.Println(utils.JsonStringify(heros))
        fmt.Println(utils.JsonStringify(playlists))
        fmt.Println(utils.JsonStringify(playlist))
    }
}


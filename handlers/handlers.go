package handlers

import (
    "github.com/aws/aws-lambda-go/events"
    "projects/chronicle-api/db"
    "projects/chronicle-api/models"
    "errors"
    "strconv"
)

type InnerHandler func(events.APIGatewayProxyRequest) (interface{}, int, error)

func UnknownRouteHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    var empty interface{}

    return empty, 404, errors.New("Unknown route")
}

func HerosHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    code := 500

    offset, limit := extractOffsetAndLimit(request)

    data, err := db.GetHeros(offset, limit)

    if err == nil {
        return data, 200, nil
    }

    return data, code, err
}

func VideoHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    data, err := db.GetVideoPackage(request.PathParameters["id"])

    if err == nil {
        return data, 200, nil
    }

    return data, 404, err
}

func MostWatchedHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    if models.ValidPeriods[request.PathParameters["period"]] {

        offset, limit := extractOffsetAndLimit(request)
        data, err := db.GetVideos(offset, limit)

        if err == nil {

            data.Uri = request.RequestContext.ResourcePath + "/" + request.PathParameters["period"]
            return data, 200, nil
        }

        return data, 500, err
    }

    var data interface{}

    return data, 404, errors.New("Undefined period")
}

func PlaylistsHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    code := 200

    offset, limit := extractOffsetAndLimit(request)

    data, err := db.GetPlaylists(offset, limit)

    if err != nil {
        code = 500
    }

    return data, code, err
}

func TopicsHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    code := 200
    offset, limit := extractOffsetAndLimit(request)
    width := extractIntegerQuery(request,"width")

    data, err := db.GetTopics(offset, limit, width)

    if err != nil {
        code = 500
    }

    return data, code, err
}


func PicksHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

    code := 200
    offset, limit := extractOffsetAndLimit(request)

    videos, err := db.GetVideos(offset, limit)

    if err != nil {
        code = 500
    }

    videos.Uri = request.RequestContext.ResourcePath

    return videos, code, err
}

func extractIntegerQuery(request events.APIGatewayProxyRequest, param string) int {

    intval, err :=  strconv.Atoi(request.QueryStringParameters[param])

    if err == nil {
        return intval
    }

    return 0
}

func extractOffsetAndLimit(request events.APIGatewayProxyRequest) (int, int) {

    return extractIntegerQuery(request, "offset"), extractIntegerQuery(request, "limit")
}


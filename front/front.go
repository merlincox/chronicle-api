package front

import (

    "fmt"

    "github.com/aws/aws-lambda-go/events"

    "projects/chronicle-api/models"
    "projects/chronicle-api/utils"
    "projects/chronicle-api/handlers"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    handler := getInnerHandler(request)

    data, statusCode, err := handler(request)

    return buildResponse(data, statusCode, err), nil
}

func getInnerHandler(request events.APIGatewayProxyRequest) handlers.InnerHandler {

    route := request.RequestContext.HTTPMethod + request.RequestContext.ResourcePath

    switch route {

    case "GET/heros":
        return handlers.HerosHandler

    case "GET/topics":
        return handlers.TopicsHandler

    case "GET/video/{id}":
        return handlers.VideoHandler

    case "GET/most-watched/{period}":
        return handlers.MostWatchedHandler

    case "GET/playlists":
        return handlers.PlaylistsHandler

    case "GET/picks":
        return handlers.PicksHandler

    }

    return handlers.UnknownRouteHandler
}

func buildResponse(data interface{}, statusCode int, err error) events.APIGatewayProxyResponse {

    var body string

    if err != nil {

        body = utils.JsonStringify(models.ErrorMessage{Message: fmt.Sprintf("Error: %v", err)})

        // if statusCode is 200 but there is an unhandled error
        // that indicates a system failure

        if statusCode == 200 {
            statusCode = 500
        }
    }

    if statusCode == 200 {
        body = utils.JsonStringify(data)
    }

    return events.APIGatewayProxyResponse{Body:body, StatusCode: statusCode, Headers: map[string]string{"Access-Control-Allow-Origin": "*"}}
}

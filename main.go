package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"fmt"
)

var (
	validPeriods = map[string]bool{
		"month": true,
		"week": true,
	}
)

type VideosResponse struct {
	Metadata ResponseMeta `json:"metadata"`
	Results []Video `json:"results"`
}

type ResponseMeta struct {
	Count int `json:"count"`
	Offset int `json:"offset"`
	Total int `json:"total"`
}

type Video struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	HoldingImageURL string `json:"holdingImageURL"`
	Topic string `json:"topic"`
	Items []VideoItem `json:"items"`
}

type VideoItem struct {
	VersionID string `json:"versionID"`
	Duration int `json:"duration"`
	Kind string `json:"kind"`
}

func getDummyData() []byte {

	dummyData := []byte(`[
    {
      "title": "Severe weather: UK buffeted by blizzards 1",
      "topic": "PLANET EARTH",
      "description":
        "Heavy snow around the UK forced thousands of schools to close, the loss of power to many homes and continued travel disruption.",
      "holdingImageURL":
        "http://news.bbcimg.co.uk/media/images/65379000/jpg/_65379468_65379467.jpg",
      "items": [
        {
          "versionID": "p013z20w",
          "duration": 142,
          "kind": "programme"
        }
      ]
    },
    {
      "title": "Severe weather: UK buffeted by blizzards 2",
      "topic": "MODERN WORLD",
      "description":
        "Heavy snow around the UK forced thousands of schools to close, the loss of power to many homes and continued travel disruption.",
      "holdingImageURL":
        "http://news.bbcimg.co.uk/media/images/65379000/jpg/_65379468_65379467.jpg",
      "items": [
        {
          "versionID": "p013z20w",
          "duration": 142,
          "kind": "programme"
        }
      ]
    }
  ]`)
	return dummyData
}

func getData() (VideosResponse, error) {

	var videoResponse VideosResponse

	err := json.Unmarshal(getDummyData(), &videoResponse.Results)

	videoResponse.Metadata = ResponseMeta{
		Offset: 0,
		Count: len(videoResponse.Results),
		Total: len(videoResponse.Results),
	}

	return videoResponse, err
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {


	data, statusCode, err := innerHandler(request)

	if statusCode == 200 {

		return events.APIGatewayProxyResponse{Body:jsonStringify(data), StatusCode: 200}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: statusCode}, err
}

func innerHandler(request events.APIGatewayProxyRequest) (interface{}, int, error) {

	var data interface{}
	var err error

	if validPeriods[request.PathParameters["period"]] {

		data, err = getData()

		if err == nil {
			return data, 200, nil
		}

		return data, 500, err
	}

	return data, 404, nil
}

func jsonStringify(data interface{}) string {

	raw, err := json.Marshal(data)

	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return string(raw[:])
}

func main() {
	lambda.Start(Handler)
	//videos, err :=  getData()
	//
	//if err != nil {
	//	fmt.Errorf("Something went wrong:" , err)
	//}
	//
	//fmt.Printf("Videos: %v\n", videos)
}

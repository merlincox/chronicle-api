package main

import (
	"fmt"
    "os"

    "github.com/aws/aws-lambda-go/lambda"

    //"projects/chronicle-api/utils"
    "projects/chronicle-api/db"
    "projects/chronicle-api/front"
    //"projects/chronicle-api/s3"
    //"log"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
    "projects/chronicle-api/utils"
)

func getAwsOptions() session.Options {

    var options session.Options

    if os.Getenv("AWS") == "true" {

        options = session.Options{
            Config: aws.Config{Region: aws.String("eu-west-2")},
        }

    } else {
        options = session.Options{
            Config: aws.Config{Region: aws.String("eu-west-2")},
            Profile: "denma-digital",
        }
    }

    return options
}

func main() {

    db.Init(getAwsOptions(), os.Getenv("BUCKET"), os.Getenv("FILENAME"))

    /*
    NB: AWS needs to be defined as an environmental variable for the Lambda
     */
    if os.Getenv("AWS") == "true" {

        lambda.Start(front.Handler)

    } else {

        heros, _ := db.GetHeros(0, 0)
        fmt.Println(utils.JsonStringify(heros))

        playlists, _ := db.GetPlaylists(0, 0)
        fmt.Println(utils.JsonStringify(playlists))

        topics, _ := db.GetTopics(0, 0)
        fmt.Println(utils.JsonStringify(topics))
    }
}


package main

import (
	"fmt"
    "os"

    "github.com/aws/aws-lambda-go/lambda"

    "projects/chronicle-api/utils"
    "projects/chronicle-api/db"
    "projects/chronicle-api/front"
)

func main() {

    db.Init()

    /*
    NB: AWS needs to be defined as an environmental variable for the Lambda
     */
    if os.Getenv("AWS") == "true" {
        lambda.Start(front.Handler)
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


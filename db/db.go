package db

import (
    "projects/chronicle-api/models"
    "encoding/json"
    "projects/chronicle-api/s3"

    "github.com/aws/aws-sdk-go/aws/session"
    "fmt"
)

var (

    videoPackageMap map[string]models.VideoPackage
    videoMap map[string]models.Video
    playlistMap map[string]models.Playlist
    topicMap map[string][]models.Video

    allVideos []models.Video
    allHeros []models.Hero
    allPlaylists []models.Playlist
    allTopics []models.VideoCollection

    s3Reader s3.Reader
    s3Bucket string
    s3Filename string
)

func GetPlaylists(offset, limit int) (models.PlaylistCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(allPlaylists))

    dummyPlaylistCollection := models.PlaylistCollection{
        Items: allPlaylists[offset:offset + limit],
        Total: len(allPlaylists),
        Offset: offset,
    }

    return dummyPlaylistCollection, error(nil)
}

func GetVideos(offset, limit int) (models.VideoCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(allVideos))

    dummyVideoCollection := models.VideoCollection{
        Items: allVideos[offset:offset + limit],
        Total: len(allVideos),
        Offset: offset,
    }

    return dummyVideoCollection, error(nil)
}

func GetTopics(offset, limit int) (models.VideoCollectionList, error) {

    offset, limit = sanitizeParams(offset, limit, len(allTopics))

    topicCollection := models.VideoCollectionList {
        Items: allTopics[offset:offset + limit],
        Total: len(allTopics),
        Offset: offset,
    }

    return topicCollection, error(nil)
}

func GetHeros(offset, limit int) (models.HeroCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(allHeros))

    herosCollection := models.HeroCollection{
        Items: allHeros[offset:offset + limit],
        Total: len(allHeros),
        Offset: offset,
    }

    return herosCollection, error(nil)
}

func GetVideoPackage(id string) (models.VideoPackage, error) {

    var video models.VideoPackage

    if video, ok := videoPackageMap[id]; ok {
        return video, nil
    }

    return video, fmt.Errorf("No such video as %v", id)
}

func Init(options session.Options, bucket string, filename string) {

    s3Reader = s3.NewReader(options)
    s3Bucket = bucket
    s3Filename = filename

    skel, _ := getSkeletonCollectionFromS3()

    ProcessSkeletons(skel)
}

func ProcessSkeletons(skels models.SkeletonCollection) {

    videoMap = make(map[string]models.Video, len(skels.Videos))
    videoPackageMap = make(map[string]models.VideoPackage, len(skels.Videos))
    playlistMap = make(map[string]models.Playlist, len(skels.Playlists))
    topicMap = make(map[string][]models.Video)

    for _, video := range skels.Videos {
        videoMap[video.Id] = video
        allVideos = append(allVideos, video)

        topicMap[video.Topic] = append(topicMap[video.Topic], video)
    }

    for topic, videos := range topicMap {

        topicCollection := models.VideoCollection{
            Items: videos,
            Total: len(videos),
            Offset: 0,
            Title: topic,
        }

        allTopics = append(allTopics, topicCollection)
    }

    for _, playlistSkel := range skels.Playlists {

        playlist := hydratePlaylist(playlistSkel)
        playlistMap[playlistSkel.Id] = playlist
        allPlaylists = append(allPlaylists, playlist)

        for _, id := range playlistSkel.Items {
            videoPackageMap[id] = models.VideoPackage{
                Items: reorderedVideos(playlist.Items, id),
            }
        }
    }

    for _, heroSkel := range skels.Heros {

        hero := hydrateHero(heroSkel)
        allHeros = append(allHeros, hero)
    }

}

func getVideo(id string) models.Video {

    return videoMap[id]
}

func getPlaylist(id string) models.Playlist {

    return playlistMap[id]
}

func getTopicVideos(topic string) []models.Video {

    return topicMap[topic]
}

func hydrateHero(skel models.HeroSkeleton) models.Hero {

    var linkSmpData models.SmpData
    var linkUri string

    if skel.LinkType == "video" {
        video := getVideo(skel.LinkId)
        linkSmpData = *video.SmpData
        linkUri = video.LinkUri
    } else {
        playlist := getPlaylist(skel.LinkId)
        video := playlist.Items[0]
        linkSmpData.Title = playlist.Title
        linkSmpData.Summary = playlist.Summary
        linkSmpData.HoldingImageURL = playlist.CoverImageUrl
        linkUri = video.LinkUri
    }

    hero := models.Hero {
        Topic: skel.Topic,
        BackgroundImage: skel.BackgroundImage,
        SponsorId: skel.SponsorId,
        SmpData: &linkSmpData,
        LinkUri: linkUri,
    }

    if skel.PreviewId != "" {
        video := getVideo(skel.PreviewId)
        hero.PreviewSmpData = video.SmpData
    }

    return hero
}

func hydratePlaylist(skel models.PlaylistSkeleton) models.Playlist {

    var videos []models.Video

    for _, id := range skel.Items {
        videos = append(videos, getVideo(id))
    }

    return models.Playlist{
        Items: videos,
        Topic: skel.Topic,
        Title: skel.Title,
        Summary: skel.Summary,
        SponsorID: skel.SponsorID,
        CoverImageUrl: skel.CoverImageUrl,
        Curator: skel.Curator,
    }
}

func sanitizeParams(offset, limit, size int) (int, int) {

    if offset < 0 {
        offset = 0
    }

    if offset > size {
        offset = size
    }

    if ( limit < 1 ) || (offset + limit > size) {
        limit = size - offset
    }

    return offset, limit
}

func reorderedVideos(input []models.Video, id string) []models.Video {

    var output []models.Video

    for _, v := range input {
        if v.Id == id {
            output = append(output, v)
        }
    }

    for _, v := range input {
        if v.Id != id {
            output = append(output, v)
        }
    }

    return output
}

func getSkeletonCollectionFromS3() (models.SkeletonCollection, error) {

    var skel models.SkeletonCollection

    raw, err := s3Reader.ReadBytes(s3Bucket, s3Filename)

    if err != nil {
        return skel, err
    }

    err = json.Unmarshal(raw, &skel)

    return skel, err
}

func getSkeletonCollectionInline() (models.SkeletonCollection, error) {

    var skel models.SkeletonCollection

    err := json.Unmarshal(getInlineSkeletonData(), &skel)

    if err != nil {
        fmt.Println(err)
    }

    return skel, err
}

func getInlineSkeletonData() []byte {

    skeletonBytes := []byte(`
{
  "videos": [
        {
            "id": "p061b9c4",
            "linkUri": "/video/p061b9c4/footage-of-first-polar-bear-cub-born-in-uk-in-25-years",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p061bbfx.jpg",
                "items": [
                    {
                        "duration": 29,
                        "kind": "programme",
                        "vpid": "p061b9c8"
                    }
                ],
                "summary": "The polar bear cub is described as a \"confident and curious\" character",
                "title": "Footage of first polar bear cub born in UK in 25 years"
            },
            "topic": "Earth",
            "uri": "/video/p061b9c4"
        },
        {
            "id": "p05ss2r1",
            "linkUri": "/video/p05ss2r1/uk-s-first-polar-bear-cub-in-25-years-born-in-highlands",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ss338.jpg",
                "items": [
                    {
                        "kind": "programme",
                        "vpid": "p05ss2r3"
                    }
                ],
                "summary": "A female polar bear at a Scottish animal park has given birth to a cub, says the Royal Zoological Society of Scotland.",
                "title": "UK's first polar bear cub in 25 years born in Highlands"
            },
            "topic": "Earth",
            "uri": "/video/p05ss2r1"
        },
        {
            "id": "p05ncdqc",
            "linkUri": "/video/p05ncdqc/doug-allan-a-life-capturing-the-natural-world-on-camera",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ncg65.jpg",
                "items": [
                    {
                        "kind": "programme",
                        "vpid": "p05ncdqf"
                    }
                ],
                "summary": "Doug Allan is one of the world's best nature cameramen and has filmed some of the most memorable scenes ever broadcast, with some close scrapes with animals along the way.",
                "title": "Doug Allan: A life capturing the natural world on camera"
            },
            "topic": "Earth",
            "uri": "/video/p05ncdqc"
        },
        {
            "id": "p059scjv",
            "linkUri": "/video/p059scjv/the-snow-might-not-last-long-with-july-temperatures-reaching-25-degrees",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p059sdxl.jpg",
                "items": [
                    {
                        "kind": "programme",
                        "vpid": "p059sdm4"
                    }
                ],
                "summary": "Christmas has come early for Lapland zoo polar bears with snow in July.",
                "title": "The snow might not last long with July temperatures reaching 25 degrees"
            },
            "topic": "Africa",
            "uri": "/video/p059scjv"
        },
        {
            "id": "p061dd3c",
            "linkUri": "/video/p061dd3c/animal-imitator-justice-osei-hopes-his-talent-will-earn-him-a-place-in-the-guinness-book-of-world-records",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p061dkdj.jpg",
                "items": [
                    {
                        "duration": 47,
                        "kind": "programme",
                        "vpid": "p061dd3h"
                    }
                ],
                "summary": "As a boy, Justice Osei from Ghana discovered he had an unusual talent: imitating the sounds of his sheep, goats and other local wildlife. Since then he has taught himself many more, and now has more than 50 species in his vocal menagerie. He performed some of them for BBC Pidgin.",
                "title": "Animal imitator Justice Osei hopes his talent will earn him a place in the Guinness Book of World Records"
            },
            "topic": "Africa",
            "uri": "/video/p061dd3c"
        },
        {
            "id": "p0529lrx",
            "linkUri": "/video/p0529lrx/the-deer-hunters-of-ghana",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p0529nd6.jpg",
                "items": [
                    {
                        "duration": 178,
                        "kind": "programme",
                        "vpid": "p0529ls0"
                    }
                ],
                "summary": "Every May the people of the Ghanaian coastal town of Winneba celebrate their migration from Timbuktu with a traditional hunt, known as the Aboakyer festival.",
                "title": "The deer hunters of Ghana"
            },
            "topic": "Africa",
            "uri": "/video/p0529lrx"
        },
        {
            "id": "p0518nt7",
            "linkUri": "/video/p0518nt7/the-villages-that-have-become-a-sanctuary-for-monkeys",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p051r5v5.jpg",
                "items": [
                    {
                        "duration": 162,
                        "kind": "programme",
                        "vpid": "p0518nt9"
                    }
                ],
                "summary": "Local people believe the monkeys are sacred and give them proper burials when they die.",
                "title": "The villages that have become a sanctuary for monkeys"
            },
            "topic": "Africa",
            "uri": "/video/p0518nt7"
        },
        {
            "id": "p04yfy06",
            "linkUri": "/video/p04yfy06/ghana-s-record-breaking-dinner-table-seats-3-600-guests",
            "smpData": {
                "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p04yjc0p.jpg",
                "items": [
                    {
                        "duration": 70,
                        "kind": "programme",
                        "vpid": "p04yfy0d"
                    }
                ],
                "summary": "A food seasoning manufacturer in Ghana has built a table which has been recognised by Guinness World Records as the longest in the world.",
                "title": "Ghana's record-breaking dinner table seats 3,600 guests"
            },
            "topic": "Africa",
            "uri": "/video/p04yfy06"
        }
    ],
  "playlists": [
    {
      "title": "Cuddly polar bears",
      "description": "Polar bears are very cute - unless you happen to be a seal.",
      "coverImageUrl": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ss338.jpg",
      "topic": "Earth",
      "id": "abcdef",
      "items": [
        "p061b9c4",
        "p05ss2r1",
        "p05ncdqc",
        "p059scjv"
      ]
    },
    {
      "title": "Funny things in Ghana",
      "description": "Funny things seem to be happening in Ghana for some reason.",
      "coverImageUrl": "https://ichef.bbci.co.uk/images/ic/$recipe/p0529nd6.jpg",
      "topic": "Africa",
      "id": "ghijk",
      "items": [
        "p061dd3c",
        "p0529lrx",
        "p0518nt7",
        "p04yfy06"
      ]
    }
  ],
  "heros": [
    {
      "linkId": "p061b9c4",
      "linkType": "video",
      "topic": "Earth"
    },
    {
      "linkId": "ghijk",
      "linkType": "playlist",
      "topic": "Africa"
    },
    {
      "linkId": "p05ss2r1",
      "linkType": "video",
      "topic": "Earth"
    },
    {
      "linkId": "abcdef",
      "linkType": "playlist",
      "topic": "Earth"
    }
  ]
}
  `)

    return skeletonBytes
}

package db

import (
    "projects/chronicle-api/models"
    "projects/chronicle-api/utils"
    "encoding/json"
    "fmt"
)

var (
    dummyVideoPackagesMap map[string]models.VideoPackage
    dummyVideoPackages []models.VideoPackage
    dummyHeros []models.Hero
    dummyPlaylistsMap map[string]models.Playlist
    dummyPlaylists []models.Playlist
)

type RawDummyVideo struct {
    Title       string `json:"title"`
    Summary string `json:"summary,omitempty"`
    Id string `json:"id"`
    HoldingImageURL string `json:"holdingImageURL"`
    Topic string `json:"topic"`
    Items []models.MediaItem `json:"items"`
}

func getDummyVideos() []byte {

    dummyData1 := []byte(`[{
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
},
{
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
},{
  "title": "Animal imitator Justice Osei hopes his talent will earn him a place in the Guinness Book of World Records",
  "id":"p061dd3c",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p061dkdj.jpg",
  "guidance": "",
  "embedRights": "allowed",
  "summary": "As a boy, Justice Osei from Ghana discovered he had an unusual talent: imitating the sounds of his sheep, goats and other local wildlife. Since then he has taught himself many more, and now has more than 50 species in his vocal menagerie. He performed some of them for BBC Pidgin.",
  "liveRewind": false,
  "simulcast": false,
  "items": [
    {
      "vpid": "p061dd3h",
      "live": false,
      "duration": 47,
      "kind": "programme"
    }
  ]
},
{
  "title": "The deer hunters of Ghana",
  "id":"p0529lrx",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p0529nd6.jpg",
  "guidance": "",
  "embedRights": "blocked",
  "summary": "Every May the people of the Ghanaian coastal town of Winneba celebrate their migration from Timbuktu with a traditional hunt, known as the Aboakyer festival.",
  "liveRewind": false,
  "simulcast": false,
  "items": [
    {
      "vpid": "p0529ls0",
      "live": false,
      "duration": 178,
      "kind": "programme"
    }
  ]
},
{
  "title": "The villages that have become a sanctuary for monkeys",
  "id":"p0518nt7",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p051r5v5.jpg",
  "guidance": "",
  "embedRights": "blocked",
  "summary": "Local people believe the monkeys are sacred and give them proper burials when they die.",
  "liveRewind": false,
  "simulcast": false,
  "items": [
    {
      "vpid": "p0518nt9",
      "live": false,
      "duration": 162,
      "kind": "programme"
    }
  ]
},
{
  "title": "Ghana's record-breaking dinner table seats 3,600 guests",
  "id":"p04yfy06",
  "holdingImageURL": "https://ichef.bbci.co.uk/images/ic/$recipe/p04yjc0p.jpg",
  "guidance": "",
  "embedRights": "blocked",
  "summary": "A food seasoning manufacturer in Ghana has built a table which has been recognised by Guinness World Records as the longest in the world.",
  "liveRewind": false,
  "simulcast": false,
  "items": [
    {
      "vpid": "p04yfy0d",
      "live": false,
      "duration": 70,
      "kind": "programme"
    }
  ]
}
]`)
    return dummyData1
}

type RawDummyPlaylist struct {
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    Id string `json:"Id"`
    Topic string `json:"topic"`
    CoverImageUrl string `json:"coverImageUrl"`
}

func getDummyPlaylists() []byte {

    dummyData := []byte(`[{

	"title" : "Cuddly polar bears",
	"description": "Polar bears are very cute - unless you happen to be a seal.",
	"coverImageUrl": "https://ichef.bbci.co.uk/images/ic/$recipe/p05ss338.jpg",
	"topic": "Earth",
	"Id" : "abcdef"
	},{

	"title" : "Funny things in Ghana",
	"description": "Funny things seem to be happening in Ghana for some reason.",
	"coverImageUrl": "https://ichef.bbci.co.uk/images/ic/$recipe/p0529nd6.jpg",
	"topic": "Africa",
	"Id" : "ghijk"
	}]`)

    return dummyData
}

func getRawPlaylists() ([]RawDummyPlaylist, error) {

    var rpa []RawDummyPlaylist

    err := json.Unmarshal(getDummyPlaylists(), &rpa)

    if err != nil {
        fmt.Println(err)
    }

    return rpa, err
}

func getRawVideos() ([]RawDummyVideo, error) {

    var rva []RawDummyVideo

    err := json.Unmarshal(getDummyVideos(), &rva)

    if err != nil {
        fmt.Println(err)
    }

    return rva, err
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

func GetPlaylists(offset, limit int) (models.PlaylistCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(dummyPlaylists))

    dummyPlaylistCollection := models.PlaylistCollection{
        Items: dummyPlaylists[offset:offset + limit],
        Total: len(dummyPlaylists),
        Offset: offset,
    }

    return dummyPlaylistCollection, error(nil)
}

func GetVideos(offset, limit int) (models.VideoPackageCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(dummyVideoPackages))

    dummyVideoPackageCollection := models.VideoPackageCollection{
        Items: dummyVideoPackages[offset:offset + limit],
        Total: len(dummyVideoPackages),
        Offset: offset,
    }

    return dummyVideoPackageCollection, error(nil)
}

func GetHeros(offset, limit int) (models.HeroCollection, error) {

    offset, limit = sanitizeParams(offset, limit, len(dummyHeros))

    dummyHeroCollection := models.HeroCollection{
        Items: dummyHeros[offset:offset + limit],
        Total: len(dummyHeros),
        Offset: offset,
    }

    return dummyHeroCollection, error(nil)
}

func GetPlaylist(id string) (models.Playlist, error) {

    var playlist models.Playlist

    if playlist, ok := dummyPlaylistsMap[id]; ok {
        return playlist, nil
    }

    return playlist, fmt.Errorf("No such playlist as %v", id)
}

func GetVideoPackage(id string) (models.VideoPackage, error) {

    var video models.VideoPackage

    if video, ok := dummyVideoPackagesMap[id]; ok {
        return video, nil
    }

    return video, fmt.Errorf("No such video as %v", id)
}

func Init() {

    var (
        dummyVideosMap map[string]models.Video
        dummyVideos []models.Video
    )

    rva, _ := getRawVideos()
    rpa, _ := getRawPlaylists()

    dummyVideosMap = make(map[string]models.Video, len(rva))
    dummyVideoPackagesMap = make(map[string]models.VideoPackage, len(rva))
    dummyPlaylistsMap = make(map[string]models.Playlist, len(rpa))

    for _, rp := range rpa {

        dummyPlaylists = append(dummyPlaylists, models.Playlist{

            Title: rp.Title,
            Topic: rp.Topic,
            Id: rp.Id,
            Description: rp.Description,
            CoverImageUrl: rp.CoverImageUrl,
            Videos: []models.Video{},
            LinkUri: "/playlist/" + rp.Id + "/" + utils.Slug(rp.Title),
        })
    }

    j := 0
    topic := "Earth"

    for i, rv := range rva {

        if i == 3 {
            topic = "Africa"
            j = 1
        }

        dummyVideosMap[rv.Id] = models.Video{
            SmpData: &models.SmpData{
                Items: rv.Items,
                Title: rv.Title,
                Summary: rv.Summary,
                HoldingImageURL: rv.HoldingImageURL,
            },
            Id: rv.Id,
            Topic: topic,
            LinkUri: "/video/" + rv.Id + "/" + utils.Slug(rv.Title),
        }

        dummyPlaylists[j].Videos = append(dummyPlaylists[j].Videos, dummyVideosMap[rv.Id])

        dummyVideos = append(dummyVideos, dummyVideosMap[rv.Id])

        var smpData models.SmpData
        if len(dummyHeros) < 2 {
            useVideo := (len(dummyHeros) % 2) == 0
            smpData = makeSmpFromVideo(dummyVideosMap[rv.Id], useVideo)
            dummyHeros = append(dummyHeros, models.Hero{
                SmpData: &smpData,
                LinkUri: dummyVideosMap[rv.Id].LinkUri,
                Topic: dummyVideosMap[rv.Id].Topic,
            })
        }
    }

    for _, pl := range dummyPlaylists {
        for _, v := range pl.Videos {

            dummyVideoPackagesMap[v.Id] = models.VideoPackage{
                Primary: &v,
                Siblings: filteredVideos(pl.Videos, v.Id),
            }

            dummyVideoPackages = append(dummyVideoPackages, dummyVideoPackagesMap[v.Id])
        }
        dummyPlaylistsMap[pl.Id] = pl
        var smpData models.SmpData
        if len(dummyHeros) < 4 {
            useVideo := (len(dummyHeros) % 2) == 0
            smpData = makeSmpFromPlaylist(pl, useVideo)
            dummyHeros = append(dummyHeros, models.Hero{
                SmpData: &smpData,
                LinkUri: pl.LinkUri,
                Topic: pl.Topic,
            })
        }

    }

}

func filteredVideos(input []models.Video, id string) []models.Video {

    var output []models.Video

    for _, v := range input {
        if v.Id != id {
            output = append(output, v)
        }
    }

    return output
}

func makeSmpFromPlaylist(pl models.Playlist, withVideo bool) models.SmpData {

    var items []models.MediaItem

    if withVideo {
        items = pl.Videos[0].SmpData.Items
    }

    return models.SmpData{
        Title: pl.Title,
        Summary: pl.Description,
        HoldingImageURL: pl.CoverImageUrl,
        Items: items,
    }
}

func makeSmpFromVideo(v models.Video, withVideo bool) models.SmpData {

    var items []models.MediaItem

    if withVideo {
        items = v.SmpData.Items
    }

    return models.SmpData{
        Title: v.SmpData.Title,
        Summary: v.SmpData.Summary,
        HoldingImageURL: v.SmpData.HoldingImageURL,
        Items: items,
    }
}

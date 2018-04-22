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
	videoMap        map[string]models.Video
	playlistMap     map[string]models.Playlist
	topicVideoMap   map[string][]models.Video
	topicMap		map[string]models.Topic

	allVideos      []models.Video
	allHeros       []models.Hero
	allPlaylists   []models.Playlist
	allTopicVideos []models.VideoCollection

	s3Reader   s3.Reader
	s3Bucket   string
	s3Filename string
)

func GetPlaylists(offset, limit int) (models.PlaylistCollection, error) {

	offset, limit = sanitizeParams(offset, limit, len(allPlaylists))

	dummyPlaylistCollection := models.PlaylistCollection{
		Items:  allPlaylists[offset : offset+limit],
		Total:  len(allPlaylists),
		Offset: offset,
	}

	return dummyPlaylistCollection, error(nil)
}

func GetVideos(offset, limit int) (models.VideoCollection, error) {

	offset, limit = sanitizeParams(offset, limit, len(allVideos))

	dummyVideoCollection := models.VideoCollection{
		Items:  allVideos[offset : offset+limit],
		Total:  len(allVideos),
		Offset: offset,
	}

	return dummyVideoCollection, error(nil)
}

func GetTopics(offset, limit, width int) (models.VideoCollectionList, error) {

	offset, limit = sanitizeParams(offset, limit, len(allTopicVideos))

	var topicCollections []models.VideoCollection

	if width > 0 {

		for _, videoCollection := range allTopicVideos[offset : offset+limit] {

			topicCollection := models.VideoCollection{
				Items:  videoCollection.Items[:width],
				Total:  len(videoCollection.Items),
				Offset: 0,
				Title:  videoCollection.Title,
				Summary: videoCollection.Summary,
			}

			topicCollections = append(topicCollections, topicCollection)
		}

	} else {
		topicCollections = allTopicVideos[offset : offset+limit]
	}

	topicCollection := models.VideoCollectionList{
		Items:  topicCollections,
		Total:  len(allTopicVideos),
		Offset: offset,
	}

	return topicCollection, error(nil)
}

func GetHeros(offset, limit int) (models.HeroCollection, error) {

	offset, limit = sanitizeParams(offset, limit, len(allHeros))

	herosCollection := models.HeroCollection{
		Items:  allHeros[offset : offset+limit],
		Total:  len(allHeros),
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
	topicVideoMap = make(map[string][]models.Video)
	topicMap = make(map[string]models.Topic)

	for _, topic := range skels.Topics {
		topicMap[topic.Id] = topic
	}

	for _, video := range skels.Videos {
		videoMap[video.Id] = video
		allVideos = append(allVideos, video)

		topicVideoMap[video.Topic] = append(topicVideoMap[video.Topic], video)
	}

	for id, videos := range topicVideoMap {

		topic, err := getTopic(id)

		if err == nil {
			topicCollection := models.VideoCollection{
				Items:  videos,
				Total:  len(videos),
				Offset: 0,
				Title:  topic.Title,
				Summary: topic.Summary,
				Uri: "/" + id,
			}

			allTopicVideos = append(allTopicVideos, topicCollection)

		}

	}

	for _, playlistSkel := range skels.Playlists {

		playlist, err := hydratePlaylist(playlistSkel)

		if err == nil {
			playlistMap[playlistSkel.Id] = playlist
			allPlaylists = append(allPlaylists, playlist)

			for _, id := range playlistSkel.Items {
				videoPackageMap[id] = models.VideoPackage{
					Items: reorderedVideos(playlist.Items, id),
				}
			}
		}
	}

	for _, heroSkel := range skels.Heros {

		hero, err := hydrateHero(heroSkel)
		if err == nil {
			allHeros = append(allHeros, hero)
		}
	}

}

func getVideo(id string) (models.Video, error) {

	video, ok := videoMap[id]

	if ok {
		return video, nil
	}

	return video, fmt.Errorf("No such video as %v", id)
}

func getTopic(id string) (models.Topic, error) {

	topic, ok := topicMap[id]

	if ok {
		return topic, nil
	}

	return topic, fmt.Errorf("No such topic as %v", id)
}

func getPlaylist(id string) (models.Playlist, error) {

	playlist, ok := playlistMap[id]

	if ok {
		return playlist, nil
	}

	return playlist, fmt.Errorf("No such playlist as %v", id)
}

func getTopicVideos(id string) ([]models.Video, error) {

	videos, ok := topicVideoMap[id]

	if ok {
		return videos, nil
	}

	return videos, fmt.Errorf("No videos for id %v", id)
}

func hydrateHero(skel models.HeroSkeleton) (models.Hero, error) {

	var hero models.Hero
	var linkSmpData models.SmpData
	var linkUri string

	if skel.LinkType == "video" {
		video, err := getVideo(skel.LinkId)

		if err != nil {
			return hero, err
		}

		linkSmpData = *video.SmpData
		linkUri = video.LinkUri

	} else {
		playlist, err := getPlaylist(skel.LinkId)

		if err != nil {
			return hero, err
		}

		if len(playlist.Items) == 0 {
			return hero, fmt.Errorf("No videos ids for playlist %v", skel.LinkId)
		}

		video := playlist.Items[0]
		linkSmpData.Title = playlist.Title
		linkSmpData.Summary = playlist.Summary
		linkSmpData.HoldingImageURL = playlist.CoverImageUrl
		linkUri = video.LinkUri
	}

	hero = models.Hero{
		Topic:           skel.Topic,
		BackgroundImage: skel.BackgroundImage,
		SponsorId:       skel.SponsorId,
		SmpData:         &linkSmpData,
		LinkUri:         linkUri,
	}

	if skel.PreviewId != "" {
		video, err := getVideo(skel.PreviewId)
		if err != nil {
			return hero, err
		}
		hero.PreviewSmpData = video.SmpData
	}

	return hero, nil
}

func hydratePlaylist(skel models.PlaylistSkeleton) (models.Playlist, error) {

	var (
		videos   []models.Video
		playlist models.Playlist
	)

	for _, id := range skel.Items {

		video, err := getVideo(id)
		if err != nil {
			return playlist, err
		}
		videos = append(videos, video)
	}

	playlist = models.Playlist{
		Items:         videos,
		Topic:         skel.Topic,
		Title:         skel.Title,
		Summary:       skel.Summary,
		SponsorID:     skel.SponsorID,
		CoverImageUrl: skel.CoverImageUrl,
		Curator:       skel.Curator,
	}

	return playlist, nil
}

func sanitizeParams(offset, limit, size int) (int, int) {

	if offset < 0 {
		offset = 0
	}

	if offset > size {
		offset = size
	}

	if (limit < 1) || (offset+limit > size) {
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

	err := json.Unmarshal(getLocalSkeletonData(), &skel)

	if err != nil {
		fmt.Println(err)
	}

	return skel, err
}


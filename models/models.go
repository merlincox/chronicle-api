// Code generated by schema-generator. DO NOT EDIT.

package models

// CuratedPlaylist: Curated playlist for GNL Chronicle
type CuratedPlaylist struct {
  Curator *Curator `json:"curator"`
  Playlist *Playlist `json:"playlist"`
}

// CuratedPlaylistCollection: Collection of curated playlists within GNL Chronicle
type CuratedPlaylistCollection struct {
  Items []CuratedPlaylist `json:"items"`
  Offset int `json:"offset"`
  Total int `json:"total"`
}

// Curator: A curator within GNL Chronicle
type Curator struct {
  Description string `json:"description,omitempty"`
  Fileid string `json:"fileid"`
  Image string `json:"image"`
  Name string `json:"name"`
}

// Hero: Hero data for GNL Chronicle
type Hero struct {
  LinkUri string `json:"linkUri"`
  Smp *SMP `json:"smp"`
  SponsorId string `json:"sponsorId,omitempty"`
  Topic string `json:"topic,omitempty"`
}

// HeroCollection: Hero collection for GNL Chronicle
type HeroCollection struct {
  Items []Hero `json:"items"`
  Offset int `json:"offset"`
  Total int `json:"total"`
}

// MediaItem: A MediaItem as passed within the array of items to play to SMP
type MediaItem struct {
  Duration int `json:"duration,omitempty"`
  Kind string `json:"kind"`
  Vpid string `json:"vpid"`
}

// Playlist: A Playlist within GNL Chronicle
type Playlist struct {
  Curator *Curator `json:"curator,omitempty"`
  Description string `json:"description,omitempty"`
  Fileid string `json:"fileid"`
  LinkUri string `json:"linkUri"`
  SponsorID string `json:"sponsorID,omitempty"`
  Title string `json:"title"`
  Topic string `json:"topic,omitempty"`
  Videos []Video `json:"videos"`
}

// PlaylistCollection: Collection of Videos within GNL Chronicle
type PlaylistCollection struct {
  Items []Playlist `json:"items"`
  Offset int `json:"offset"`
  Total int `json:"total"`
}

// SMP: SMP data object
type SMP struct {
  Guidance string `json:"guidance,omitempty"`
  HoldingImageURL string `json:"holdingImageURL"`
  Items []MediaItem `json:"items,omitempty"`
  Summary string `json:"summary,omitempty"`
  Title string `json:"title"`
}

// Topic: A topic within GNL Chronicle
type Topic struct {
  Description string `json:"description,omitempty"`
  Fileid string `json:"fileid"`
  Name string `json:"name"`
}

// TopicVideoCollection: Collection of Videos by Topic within GNL Chronicle
type TopicVideoCollection struct {
  Offset int `json:"offset"`
  Topic *Topic `json:"topic"`
  Total int `json:"total"`
  Videos []Video `json:"videos"`
}

// TopicVideoCollectionCollection: Collection of collections of videos by topic within GNL Chronicle
type TopicVideoCollectionCollection struct {
  Items []TopicVideoCollection `json:"items"`
  Offset int `json:"offset"`
  Total int `json:"total"`
}

// Video: A Video within GNL Chronicle
type Video struct {
  Breakpoints []int `json:"breakpoints,omitempty"`
  Id string `json:"id"`
  LinkUri string `json:"linkUri"`
  Smp *SMP `json:"smp"`
  SponsorID string `json:"sponsorID,omitempty"`
  Tags []string `json:"tags,omitempty"`
  Topic string `json:"topic"`
}

// VideoCollection: Collection of Videos within GNL Chronicle
type VideoCollection struct {
  Items []Video `json:"items"`
  Offset int `json:"offset"`
  Total int `json:"total"`
}

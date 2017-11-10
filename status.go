package spotilocal

type Status struct {
	// DO NOT include Version and ClientVersion. Adding so will never cause ErrEmptyResponse.
	// Version         int     `json:"version"`
	// ClientVersion   string  `json:"client_version"`
	Playing         bool    `json:"playing"`
	PlayingPosition float64 `json:"playing_position"`
	Resources       struct {
		Track  Resource `json:"track_resource"`
		Artist Resource `json:"artist_resource"`
		Album  Resource `json:"album_resource"`
	} `json:"track"`
}

type Resource struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

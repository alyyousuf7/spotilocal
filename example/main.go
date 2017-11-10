package main

import (
	"fmt"
	"time"

	"github.com/alyyousuf7/spotilocal"
)

func main() {
	spotify, err := spotilocal.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait for user to open up Spotify Desktop Client
	WaitForClient(spotify)

	for {
		status, err := spotify.StatusOn(spotilocal.OnPlay)
		if err != nil {
			// Ignore empty responses
			if err == spotilocal.ErrEmptyResponse {
				continue
			}

			fmt.Println(err)
			if err != spotilocal.ErrDisconnected {
				return
			}

			// Since it disconnected, lets wait for it again
			WaitForClient(spotify)
		}

		playingStatus := "Playing"
		if !status.Playing {
			playingStatus = "Paused"
		}
		fmt.Printf("Playing Status: %s\nTrack:\t%s\nArtist:\t%s\nAlbum:\t%s\n\n",
			playingStatus,
			status.Resources.Track.Name,
			status.Resources.Artist.Name,
			status.Resources.Album.Name,
		)
	}
}

func Connect(spotify *spotilocal.Spotify) error {
	if err := spotify.Connect(); err != nil {
		return err
	}

	version, clientVersion, err := spotify.Version()
	if err != nil {
		return err
	}

	fmt.Printf("Connected\nVersion: %d\nClient Version: %s\n\n",
		version,
		clientVersion,
	)
	return nil
}

func WaitForClient(spotify *spotilocal.Spotify) {
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-ticker:
			fmt.Println("Trying to connect...")
			if err := Connect(spotify); err == nil {
				return
			}
		}
	}
}

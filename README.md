# Spotilocal
Spotify Desktop Client when running, exposes a HTTPS port which can be used to perform a few actions on the client.
Spotilocal provides a simple API in Golang to perform those actions. It can be fetched using:

```bash
$ go get -d github.com/alyyousuf7/spotilocal
```

## Example
The following will play `Thunder - Imagine Dragon` on an already running Spotify Desktop Client.

```golang
package main

import (
    "fmt"

    "github.com/alyyousuf7/spotilocal"
)

func main() {
    spotify, err := spotilocal.New()
    if err != nil {
        panic(err)
    }

    if err := spotify.Connect(); err != nil {
		panic(err)
	}

    if err := spotify.Play("spotify:track:5VnDkUNyX6u5Sk0yZiP8XB"); err != nil {
		panic(err)
    }
    
    status, err := spotify.StatusOn(spotilocal.OnPlay)
    if err != nil {
		panic(err)
    }

    fmt.Printf("Track:\t%s\nArtist:\t%s\nAlbum:\t%s\n",
        status.Resources.Track.Name,
        status.Resources.Artist.Name,
        status.Resources.Album.Name,
    )
}
```

A detailed example is available in `example/main.go` file.

## APIs
- `Pause(status bool) error` - Pauses a running track on `status = true`; Unpauses on `status = false`
- `Play(uri string) error` - Plays a track
- `PlayWithContext(uri, context string) error` - Plays a track in context with an album or artist
- `Status() (Status, error)` - Gives client's current status in `Status` struct
- `StatusOn(...Event) (Status, error)` - Gives client's status in `Status` struct on an Event; See below for `Event`
- `Version() (string, string, error)` - Gives a pair of version (`version` and `client_version` respectively) from client API

## Events
`Event` is used when requesting for status. Available statuses are:

- `OnAP`
- `OnError`
- `OnLogin`
- `OnLogout`
- `OnPause`
- `OnPlay`

## Troubleshoot
While playing around with Spotify Client Desktop, I faced some issues where `spotilocal` was able to connect but was receiving empty responses from all the calls.
I figured out that sometimes Spotify Desktop Client does not properly closes their ports.
If anything similar happens, try to kill the Spotify Desktop Client using `sudo pkill spotify`.

## License
MIT
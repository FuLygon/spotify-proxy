# Proxy for Spotify API

The actual one I use for [NowPlaying for Spotify](https://github.com/busybox11/NowPlaying-for-Spotify) and my [Homepage](https://github.com/gethomepage/homepage) dashboard.

- For the get current track API, if no track is playing, it will return recently played track data instead, preventing NowPlaying from showing "No track playing."
- For the get current queue API, the response will be cached for my Homepage dashboard, preventing an empty list from being shown.
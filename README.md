## invidious-go

Go bindings for the Invidious API.

Documentation: https://pkg.go.dev/github.com/antoniszymanski/invidious-go

### Installation:

```
go get github.com/antoniszymanski/invidious-go
```

### Endpoints:

| Endpoint                                        | Status | Notes                |
| ----------------------------------------------- | ------ | -------------------- |
| GET /api/v1/stats                               | ❌     |                      |
| GET /api/v1/videos/:id                          | ❌     |                      |
| GET /api/v1/annotations/:id                     | ❌     |                      |
| GET /api/v1/comments/:id                        | ❌     |                      |
| GET /api/v1/captions/:id                        | ❌     |                      |
| GET /api/v1/trending                            | ❌     |                      |
| GET /api/v1/popular                             | ❌     |                      |
| GET /api/v1/search/suggestions                  | ❌     |                      |
| GET /api/v1/search                              | ❌     |                      |
| GET /api/v1/playlists/:plid                     | ❌     |                      |
| GET /api/v1/mixes/:rdid                         | ❌     |                      |
| GET /api/v1/hashtag/:tag                        | ❌     |                      |
| GET /api/v1/resolveurl                          | ❌     |                      |
| GET /api/v1/clips                               | ❌     |                      |
| GET /api/v1/channels/:id                        | ❌     |                      |
| GET /api/v1/channels/:id/channels               | ❌     |                      |
| GET /api/v1/channels/:id/latest                 | ❌     |                      |
| GET /api/v1/channels/:id/playlists              | ❌     |                      |
| GET /api/v1/channels/:id/podcasts               | ❌     |                      |
| GET /api/v1/channels/:id/releases               | ❌     |                      |
| GET /api/v1/channels/:id/shorts                 | ❌     |                      |
| GET /api/v1/channels/:id/streams                | ❌     |                      |
| GET /api/v1/channels/:id/videos                 | ❌     |                      |
| GET /api/v1/channels/:id/community              | ❌     |                      |
| GET /api/v1/channels/:ucid/search               | ❌     |                      |
| GET /api/v1/post/:id                            | ❌     |                      |
| GET /api/v1/post/:id/comments                   | ❌     |                      |
| GET /authorize_token                            | ✅     |                      |
| GET /api/v1/auth/feed                           | ✅     |                      |
| GET /api/v1/auth/notifications                  | ❌     | Won't be implemented |
| POST /api/v1/auth/notifications                 | ❌     | Won't be implemented |
| GET /api/v1/auth/playlists                      | ✅     |                      |
| POST /api/v1/auth/playlists                     | ✅     |                      |
| GET /api/v1/auth/playlists/:id                  | ✅     |                      |
| PATCH /api/v1/auth/playlists/:id                | ✅     |                      |
| DELETE /api/v1/auth/playlists/:id               | ✅     |                      |
| POST /api/v1/auth/playlists/:id/videos          | ✅     |                      |
| DELETE /api/v1/auth/playlists/:id/videos/:index | ✅     |                      |
| GET /api/v1/auth/preferences                    | ✅     |                      |
| POST /api/v1/auth/preferences                   | ❌     |                      |
| GET /api/v1/auth/subscriptions                  | ✅     |                      |
| POST /api/v1/auth/subscriptions/:ucid           | ✅     |                      |
| DELETE /api/v1/auth/subscriptions/:ucid         | ✅     |                      |
| GET /api/v1/auth/tokens                         | ✅     |                      |
| POST /api/v1/auth/tokens/register               | ✅     |                      |
| POST /api/v1/auth/tokens/unregister             | ✅     |                      |
| GET /api/v1/auth/history                        | ✅     |                      |
| POST /api/v1/auth/history/:id                   | ✅     |                      |
| DELETE /api/v1/auth/history/:id                 | ✅     |                      |

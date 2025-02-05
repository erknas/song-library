## SongLibrary

REST API song manager project.

### Endpoints

1. `/songs`

- **[POST]** — Add new song.
- **[GET]** — Get songs with pagination and optional with filters by song, group and realease date.

2. `/song`

- **[GET]** — Get song text by verses with pagination.
- **[DELETE]** — Delete song by ID.
- **[PUT]** — Update song by ID.

3. `/swagger/index.html`
   Or can run in Swagger UI.

Examples:

1. Get songs with pagination and filters

`/songs?page=1&size=10&song=Supermassive%20Black%20Hole&group=Muse&date=16.07.2016`

2. Get song text

`/song?id=1&page=1&size=1`

### RUN

.env file stores all environment variables.

```
docker-compose up
```

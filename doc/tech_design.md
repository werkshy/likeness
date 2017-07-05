# Technical Design

## Data Store: 
Postgres on home server

## Data Models

### Photographer

### Photo
- path
- exif_date
- file_date
- photographer   e.g. alex
- camera_type   e.g. Nexus 5X
- camera_id       e.g. some unique id
- checksum_type   e.g. sha1 / md5sum
- checksum_value

### Album?

## Usage

`photodb index`
-  Scan the main dir (load into db)
- Show possible dupes

`photodb import-dir DIR`
- Import (move into main dir) and insert into the DB any new files
- List possible dupes
- Optionally delete definite dupes (matching checksum)

`photodb import-android`
- Import (move into main dir) and insert into the DB any new files on the phone
- List possible dupes
- Optionally delete definite dupes (matching checksum)

`photodb gen-thumbs`

`photodb serve [--auto-thumbnail]`
- serve a web interface to navigate folders
- react/redux frontend



## Config
```
main_dir
store_host
thumb_dir
listen_port
listen_host
```

## Technical Notes
- Producer/consumer for indexing and thumbnailing
- Add stuff to queue and process in worker threads
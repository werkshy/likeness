TODO
====

## Index Existing Photos

- [x] test tree walking producer test
- [x] test filters
- [x] write hash
- [x] write db schema
- [x] store files in db during index
    - [x] path & hash
    - [x] Detect conflicts / dupes
    - [x] Convert to sqlx
    - [x] Print list of duplicate files
    - [x] File date
    - [x] exif date
    - [x] mp4 date
    - [ ] .MOV date
- [ ] Handle dupes (VERY careful deletion..)
      - which one to delete? There's one that's not in the DB. Should we assume that whatever we are indexing is somewhat reasonably laid out?
      - Should we add '--delete-dupes' and '--move-existing' config options for 'index' command?
      - [ ] Convert to e.g. Cobra
- [ ] TBD should we move existing files to new standard naming convention?

## Import New Directories

- [ ] Walk dir, generate hashes
- [ ] Detect & eliminate duplicates
- [ ] Calculate import path
- [ ] Carefully move (handle existing files carefully, do not overwrite, make that atomic if
      possible)

## Thumbnail generation
- [ ] Generate a thumb on demand
- [ ] Pre-generate thumbnails for any photos that don't have one.

## Web view
- [ ] Web server with hierarchical year/month view
      - Based on date from DB or from filesystem?

## Advanced Dupe Detection

- [ ] File date or Exif date is identical - build a web view with possible dupes and the ability to delete one?

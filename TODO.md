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
    - [ ] File date
    - [ ] exif date & camera model / id
- [ ] Handle dupes (VERY careful deletion..)
      - which one to delete? There's one that's not in the DB. Should we assume that whatever we are indexing is somewhat reasonably laid out?
- [ ] TBD should we move existing files to new standard naming convention?

## Import New Directories

- [ ] Walk dir, generate hashes
- [ ] Detect & eliminate duplicates
- [ ] Calculate import path
- [ ] Carefully move (handle existing files carefully, do not overwrite, make that atomic if
      possible)


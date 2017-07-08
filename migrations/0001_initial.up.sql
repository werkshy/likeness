CREATE TABLE photographers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX ON photographers (name);

CREATE TABLE photos (
    id BIGSERIAL PRIMARY KEY,
    path VARCHAR(1024) NOT NULL,
    exif_date TIMESTAMP WITHOUT TIME ZONE,
    file_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    camera_type VARCHAR(255),
    camera_id VARCHAR(255),
    checksum_type VARCHAR(8) NOT NULL,
    checksum_value BYTEA NOT NULL,
    photographer_id BIGINT REFERENCES photographers(id),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON photos (path);
CREATE UNIQUE INDEX ON photos (checksum_type, checksum_value);

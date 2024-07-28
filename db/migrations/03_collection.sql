-- +goose Up
CREATE TABLE Collection (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
CREATE UNIQUE INDEX idx_collection_name ON Collection (name);

CREATE TABLE CardCollection (
    card_id INTEGER NOT NULL,
    collection_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (card_id, collection_id),
    FOREIGN KEY (card_id) REFERENCES Card (id),
    FOREIGN KEY (collection_id) REFERENCES Collection (id)
);

CREATE TABLE CollectionEnergy (
    collection_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (collection_id, type_id),
    FOREIGN KEY (collection_id) REFERENCES Collection (id),
    FOREIGN KEY (type_id) REFERENCES Type (id)
);

-- +goose Down
DROP INDEX idx_collection_name;
DROP TABLE Collection;
DROP TABLE CardCollection;
DROP TABLE CollectionEnergy;

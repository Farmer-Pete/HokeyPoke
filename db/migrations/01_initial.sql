-- +goose Up
CREATE TABLE Card (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ptcg_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    metadata JSON,
    group_id INTEGER NOT NULL,
    supertype_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES "Group" (id),
    FOREIGN KEY (supertype_id) REFERENCES Supertype (id)
);

CREATE TABLE "Group" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    supertype_id INTEGER NOT NULL,
    FOREIGN KEY (supertype_id) REFERENCES Supertype (id)
);

CREATE TABLE Supertype (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE Type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE CardType (
    card_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    PRIMARY KEY (card_id, type_id),
    FOREIGN KEY (card_id) REFERENCES Card (id),
    FOREIGN KEY (type_id) REFERENCES Type (id)
);

CREATE TABLE GroupType (
    group_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    PRIMARY KEY (group_id, type_id),
    FOREIGN KEY (group_id) REFERENCES "Group" (id),
    FOREIGN KEY (type_id) REFERENCES Type (id)
);

-- +goose Down
DROP TABLE Card;
DROP TABLE "Group";
DROP TABLE Supertype;
DROP TABLE Type;
DROP TABLE CardType;
DROP TABLE GroupType;

-- +goose Up
CREATE TABLE Deck (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    collection_id INTEGER NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES Collection (id)
);

CREATE TABLE CardDeck (
    card_id INTEGER NOT NULL,
    deck_id INTEGER NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    PRIMARY KEY (card_id, deck_id),
    FOREIGN KEY (card_id) REFERENCES Card (id),
    FOREIGN KEY (deck_id) REFERENCES Deck (id)
);

-- +goose Down
DROP TABLE Deck;
DROP TABLE CardDeck;

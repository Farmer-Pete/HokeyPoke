-- +goose Up
CREATE UNIQUE INDEX idx_group_name ON "Group" (name);
CREATE UNIQUE INDEX idx_type_name ON Type (name);
CREATE UNIQUE INDEX idx_supertype_name ON Supertype (name);
CREATE UNIQUE INDEX idx_cardtype_card_id_type_id ON CardType (card_id, type_id);
CREATE UNIQUE INDEX idx_grouptype_group_id_type_id ON GroupType (group_id, type_id);

-- +goose Down
DROP INDEX idx_group_name;
DROP INDEX idx_type_name;
DROP INDEX idx_supertype_name;
DROP INDEX idx_cardtype_card_id_type_id;
DROP INDEX idx_grouptype_group_id_type_id;

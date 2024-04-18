CREATE TABLE echo
(
    id       SERIAL UNIQUE NOT NULL,
    title    VARCHAR DEFAULT '',
    reminder VARCHAR       not null,
    CONSTRAINT unq_id_title_echo unique (id, title)

);
CREATE INDEX btree_index_name ON echo USING btree (id);

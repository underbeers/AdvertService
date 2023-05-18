CREATE TABLE favorites
(
    id              serial primary key,
    advert_id       int references advert_pet (id) ON DELETE CASCADE,
    organization_id int ,
    specialist_id   int,
    event_id        int,
    user_id         UUID
);
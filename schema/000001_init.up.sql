CREATE TABLE advert_pet
(
    id          serial primary key,
    pet_card_id int references pet_card (id) ON DELETE CASCADE,
    user_id     UUID,
    price       int,
    description varchar(2000),
    region      varchar(255),
    locality    varchar(255),
    chat        boolean,
    phone       varchar(255),
    status      varchar(255),
    publication timestamp
);
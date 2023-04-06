CREATE TABLE advert_pet
(
    id          serial primary key,
    pet_card_id int,
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
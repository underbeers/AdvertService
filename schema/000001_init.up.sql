CREATE TABLE city
(
    id   serial primary key,
    city varchar(255)
);


CREATE TABLE district
(
    id       serial primary key,
    city_id  int references city (id) ON DELETE CASCADE,
    district varchar(255)
);

CREATE TABLE advert_pet
(
    id          serial primary key,
    pet_card_id int references pet_card (id) ON DELETE CASCADE,
    user_id     UUID,
    price       int,
    description varchar(2000),
    city_id     int references city (id) ON DELETE CASCADE,
    district_id int references district (id) ON DELETE CASCADE,
    chat        boolean,
    phone       varchar(255),
    status      varchar(255),
    publication timestamp
);
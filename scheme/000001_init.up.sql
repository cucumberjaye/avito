CREATE TABLE users
(
    id      integer not null unique,
    name    varchar(255) not null,
    surname varchar(255) not null
);

CREATE TABLE balance
(
    user_id integer references users (id) not null,
    balance integer
);

CREATE TABLE transactions
(
    user_id             integer references users (id) not null,
    sum                 integer not null,
    comment             varchar(255) not null,
    transaction_date    date not null
)
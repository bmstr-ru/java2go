create table if not exists deal
(
    id              bigserial not null primary key,
    deal_id         bigint    not null,
    client_id       bigint    not null,
    bought_amount   numeric   not null,
    bought_currency char(3)   not null,
    sold_amount     numeric   not null,
    sold_currency   char(3)   not null
);

create table if not exists client_exposure
(
    id                      bigserial not null primary key,
    client_id               bigint    not null,
    total_exposure_amount   numeric   not null,
    total_exposure_currency char(3)   not null
);

create table if not exists client_exposure_detail
(
    id                bigserial not null primary key,
    client_id         bigint    not null,
    exposure_amount   numeric   not null,
    exposure_currency char(3)   not null
);

create table if not exists currency_rate
(
    id              bigserial not null primary key,
    base_currency   char(3)   not null,
    quoted_currency char(3)   not null,
    rate            numeric   not null
);
create type currency_code as enum ('PLN', 'EUR', 'MKD', 'unknown');

alter table accomodations
    add column price numeric(12, 2),
    add column currency currency_code,
    add column description text;

alter table transport
    add column price numeric(12, 2),
    add column currency currency_code,
    add column description text;

alter table activities
    add column price numeric(12, 2),
    add column currency currency_code,
    add column description text;

alter table accomodations
    drop column if exists price,
    drop column if exists currency,
    drop column if exists description;

alter table transport
    drop column if exists price,
    drop column if exists currency,
    drop column if exists description;

alter table activities
    drop column if exists price,
    drop column if exists currency,
    drop column if exists description;

drop type if exists currency_code;

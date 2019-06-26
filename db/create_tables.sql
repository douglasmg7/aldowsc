-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Products.
create table product (
    code                    text primary key,
    brand                   text not null,
    category                text not null,
    description             text not null,
    unit                    text not null,
    multiple                integer not null,
    dealer_price            real not null,
    suggestion_price        real not null,
    technical_description   text not null,
    availability            boolean not null,
    length                  integer not null, -- mm.
    width                   integer not null, -- mm.
    height                  integer not null, -- mm.
    weight                  integer not null, -- gr.
    picture_link            blob,
    warranty_period         integer,  -- Days.
    rma_procedure           text,
    created_at              date not null,
    changed_at              date not null,
    changed                 boolean not null,
    new                     boolean not null,
    removed                 boolean not null
);

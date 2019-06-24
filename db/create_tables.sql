-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Products.
create table product (
  code                text primary key,
  brand               text not null,
  category            text not null,
  description         text not null,
  multiple            integer not null,
  dealerPrice         real not null,
  suggestionPrice     real not null,
  tecnicalDescription text not null,
  availability        integer not null,
  length              integer not null, -- mm.
  width               integer not null, -- mm.
  hight               integer not null, -- mm.
  Weight              integer not null, -- gr.
  pictureLinks        blob,
  warrantyPeriod      integer,  -- Days.
  RMAProcedure        text,
  CreatedAt           date,
  ChangedAt           date,
  changed             boolean not null,
  New                 boolean not null,
  Removed             boolean not null
);

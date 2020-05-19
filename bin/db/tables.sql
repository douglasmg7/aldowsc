-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Parameters.
CREATE TABLE IF NOT EXISTS param (
	name                    TEXT PRIMARY KEY,	-- Name without space.
	value                   TEXT
);

-- Categories.
CREATE TABLE IF NOT EXISTS category (
	name                    TEXT PRIMARY KEY,	-- Name without space.
	text                    TEXT NOT NULL,
	products_qty            INTEGER NOT NULL,
	selected                BOOLEAN NOT NULL
);

-- Products.
CREATE TABLE IF NOT EXISTS product (
	mongodb_id				TEXT DEFAULT "",	-- Store id from mongodb.
	code                    TEXT NOT NULL UNIQUE,	-- Come from dealer.
    store_product_id		INTEGER,
	brand                   TEXT NOT NULL,
	category                TEXT NOT NULL,
	description             TEXT NOT NULL,
	dealer_price            INTEGER NOT NULL,
	suggestion_price        INTEGER NOT NULL,
	technical_description   TEXT NOT NULL,
	availability            BOOLEAN NOT NULL,
	length                  INTEGER NOT NULL, -- mm.
	width                   INTEGER NOT NULL, -- mm.
	height                  INTEGER NOT NULL, -- mm.
	weight                  INTEGER NOT NULL, -- gr.
	picture_link            BLOB,
	warranty_period         INTEGER,  -- Months.
	rma_procedure           TEXT,
	created_at              DATE NOT NULL,
	changed_at              DATE NOT NULL,
	removed_at              DATE DEFAULT "0001-01-01 00:00:00+00:00",
    status_cleaned_at       DATE DEFAULT "0001-01-01 00:00:00+00:00"
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_code ON product(code);

-- Products history.
CREATE TABLE IF NOT EXISTS product_history (
	mongodb_id				TEXT DEFAULT "",	-- Store id from mongodb.
	code					TEXT NOT NULL,	-- Come from dealer.
	store_product_id		INTEGER,
	brand                   TEXT NOT NULL,
	category                TEXT NOT NULL,
	description             TEXT NOT NULL,
	dealer_price            INTEGER NOT NULL,
	suggestion_price        REAL NOT NULL,
	technical_description   TEXT NOT NULL,
	availability            BOOLEAN NOT NULL,
	length                  INTEGER NOT NULL, -- mm.
	width                   INTEGER NOT NULL, -- mm.
	height                  INTEGER NOT NULL, -- mm.
	weight                  INTEGER NOT NULL, -- gr.
	picture_link            BLOB,
	warranty_period         INTEGER,  -- months.
	rma_procedure           TEXT,
	created_at              DATE NOT NULL,
	changed_at              DATE NOT NULL,
	removed_at              DATE DEFAULT "0001-01-01 00:00:00+00:00",
    status_cleaned_at       DATE DEFAULT "0001-01-01 00:00:00+00:00",
	UNIQUE (code, changed_at)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_history_code_changed_at ON product_history(code, changed_at);

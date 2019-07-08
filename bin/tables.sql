-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Products.
CREATE TABLE IF NOT EXISTS product (
	id						INTEGER PRIMARY KEY AUTOINCREMENT,
	code                    TEXT NOT NULL UNIQUE,
	brand                   TEXT NOT NULL,
	category                TEXT NOT NULL,
	description             TEXT NOT NULL,
	unit                    TEXT NOT NULL,
	multiple                INTEGER NOT NULL,
	dealer_price            REAL NOT NULL,
	suggestion_price        REAL NOT NULL,
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
	changed                 BOOLEAN NOT NULL,
	new                     BOOLEAN NOT NULL,
	removed                 BOOLEAN NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_code ON product(code);

-- Products history.
CREATE TABLE IF NOT EXISTS product_history (
	id						INTEGER PRIMARY KEY AUTOINCREMENT,
	code					TEXT NOT NULL,
	brand                   TEXT NOT NULL,
	category                TEXT NOT NULL,
	description             TEXT NOT NULL,
	unit                    TEXT NOT NULL,
	multiple                INTEGER NOT NULL,
	dealer_price            REAL NOT NULL,
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
	changed                 BOOLEAN NOT NULL,
	new                     BOOLEAN NOT NULL,
	removed                 BOOLEAN NOT NULL,
	UNIQUE (code, changed_at)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_product_history_code_changed_at ON product_history(code, changed_at);

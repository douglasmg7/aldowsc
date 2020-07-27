---------------------------------------------------------------------------------------------------------
-- Products history
---------------------------------------------------------------------------------------------------------
BEGIN TRANSACTION;

-- Create backup table.
CREATE TEMPORARY TABLE product_history_backup (
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

-- Copy data to backup table.
INSERT INTO product_history_backup
(
    mongodb_id, 
    code, 
    store_product_id,
    brand, 
    category, 
    description, 
    dealer_price, 
    suggestion_price, 
    technical_description, 
    availability,
    length, 
    width, 
    height, 
    weight, 
    picture_link, 
    warranty_period, 
    rma_procedure, 
    created_at, 
    changed_at
)
SELECT
    mongodb_id, 
    code, 
    store_product_id,
    brand, 
    category, 
    description, 
    dealer_price, 
    suggestion_price, 
    technical_description, 
    availability,
    length, 
    width, 
    height, 
    weight, 
    picture_link, 
    warranty_period, 
    rma_procedure, 
    created_at, 
    changed_at
FROM product_history;

-- Drop old table.
DROP TABLE product_history;

-- Create new table.
CREATE TABLE product_history (
	mongodb_id				TEXT DEFAULT "",	-- Store id from mongodb.
	code					TEXT NOT NULL,	-- Come from dealer.
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
	warranty_period         INTEGER,  -- months.
	rma_procedure           TEXT,
	created_at              DATE NOT NULL,
	changed_at              DATE NOT NULL,
	removed_at              DATE DEFAULT "0001-01-01 00:00:00+00:00",
    status_cleaned_at       DATE DEFAULT "0001-01-01 00:00:00+00:00",
	UNIQUE (code, changed_at)
);
CREATE UNIQUE INDEX idx_product_history_code_changed_at ON product_history(code, changed_at);

-- Copy data to new table.
INSERT INTO product_history
(
    mongodb_id, 
    code, 
    store_product_id,
    brand, 
    category, 
    description, 
    dealer_price, 
    suggestion_price, 
    technical_description, 
    availability,
    length, 
    width, 
    height, 
    weight, 
    picture_link, 
    warranty_period, 
    rma_procedure, 
    created_at, 
    changed_at
)
SELECT
    mongodb_id, 
    code, 
    store_product_id,
    brand, 
    category, 
    description, 
    dealer_price, 
    suggestion_price, 
    technical_description, 
    availability,
    length, 
    width, 
    height, 
    weight, 
    picture_link, 
    warranty_period, 
    rma_procedure, 
    created_at, 
    changed_at
FROM product_history_backup;

-- Drop backup table.
DROP TABLE product_history_backup;

COMMIT;

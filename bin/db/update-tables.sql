---------------------------------------------------------------------------------------------------------
-- Products
---------------------------------------------------------------------------------------------------------
BEGIN TRANSACTION;

-- Create backup table.
CREATE TEMPORARY TABLE product_backup (
	mongodb_id				TEXT DEFAULT "",	-- Store id from mongodb.
	code                    TEXT NOT NULL UNIQUE,	-- Come from dealer.
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
	changed                 BOOLEAN NOT NULL,
	new                     BOOLEAN NOT NULL,
	removed                 BOOLEAN NOT NULL,
	store_product_id		INTEGER
);

-- Copy data to backup table.
INSERT INTO product_backup
(
    mongodb_id, 
    code, 
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
    changed_at,
    changed,
    new,
    removed,
    store_product_id
)
SELECT
    mongodb_id, 
    code, 
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
    changed_at,
    changed,
    new,
    removed,
    store_product_id
FROM product;

-- Drop old table.
DROP TABLE product;

-- Create new table.
CREATE TABLE product (
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
	removed_at              DATE,
    checked_change          BOOLEAN DEFAULT 0
);
CREATE UNIQUE INDEX idx_product_code ON product(code);

-- Copy data to the new table.
INSERT INTO product
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
FROM product_backup;

-- Drop backup table.
DROP TABLE product_backup;

COMMIT;

---------------------------------------------------------------------------------------------------------
-- Products history
---------------------------------------------------------------------------------------------------------
BEGIN TRANSACTION;

-- Create backup table.
CREATE TEMPORARY TABLE product_history_backup (
	mongodb_id				TEXT DEFAULT "",	-- Store id from mongodb.
	code					TEXT NOT NULL,	-- Come from dealer.
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
	changed                 BOOLEAN NOT NULL,
	new                     BOOLEAN NOT NULL,
	removed                 BOOLEAN NOT NULL,
	store_product_id		INTEGER,
	UNIQUE (code, changed_at)
);

-- Copy data to backup table.
INSERT INTO product_history_backup
(
    mongodb_id, 
    code, 
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
    changed_at,
    changed,
    new,
    removed,
    store_product_id
)
SELECT
    mongodb_id, 
    code, 
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
    changed_at,
    changed,
    new,
    removed,
    store_product_id
FROM product_history;

-- Drop old table.
DROP TABLE product_history;

-- Create new table.
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
	removed_at              DATE,
    checked_change          BOOLEAN DEFAULT 0,
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

-- Copy data to new table.
INSERT INTO product_new
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

-- Rename new table to deleted table.
ALTER TABLE product_new RENAME TO product;

-- Copy data to new table.
INSERT INTO product_history_new
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

-- Rename new table to deleted table.
ALTER TABLE product_history_new RENAME TO product_history;

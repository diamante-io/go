-- +migrate Up

ALTER TABLE offers ADD deleted boolean DEFAULT false;

CREATE INDEX best_offer ON offers USING BTREE (selling_asset, buying_asset, deleted, price);
CREATE INDEX live_offers ON offers USING BTREE (deleted, last_modified_ledger);

DROP INDEX offers_by_seller, offers_by_selling_asset, offers_by_buying_asset;

CREATE INDEX offers_by_seller ON offers USING BTREE(seller_id, deleted);
CREATE INDEX offers_by_selling_asset ON offers USING BTREE(selling_asset, deleted);
CREATE INDEX offers_by_buying_asset ON offers USING BTREE(buying_asset, deleted);

-- +migrate Down

DELETE FROM offers where deleted = true;

DROP INDEX offers_by_seller, offers_by_selling_asset, offers_by_buying_asset;

ALTER TABLE offers DROP COLUMN deleted;

CREATE INDEX offers_by_seller ON offers USING BTREE(seller_id);
CREATE INDEX offers_by_selling_asset ON offers USING BTREE(selling_asset);
CREATE INDEX offers_by_buying_asset ON offers USING BTREE(buying_asset);

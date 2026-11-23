-- +goose Up
CREATE TABLE item_photos (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             item_id INT NOT NULL,
                             filename VARCHAR(255) NOT NULL,
                             created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

                             INDEX idx_item_photos_item_id (item_id),

                             CONSTRAINT fk_item_photos_item
                                 FOREIGN KEY (item_id)
                                     REFERENCES items(id)
                                     ON DELETE CASCADE
);

-- +goose Down
DROP TABLE item_photos;

CREATE TABLE `tokens`
(
    `id`            INT unsigned NOT NULL AUTO_INCREMENT,
    `user_id`       INT unsigned NOT NULL,
    `token`         VARCHAR(255) NOT NULL,
    `token_type_id` INT unsigned NOT NULL,
    `meta`          JSON,
    `expires_at`    TIMESTAMP    NULL DEFAULT CURRENT_TIMESTAMP,
    `created_at`    TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP         DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `fk_tokens_user_id_idx` (`user_id` ASC),
    INDEX `fk_tokens_token_type_id_idx` (`token_type_id` ASC),
    CONSTRAINT `fk_tokens_user_id`
        FOREIGN KEY (`user_id`)
            REFERENCES `users` (`id`)
            ON DELETE CASCADE
            ON UPDATE CASCADE,
    CONSTRAINT `fk_tokens_token_type_id`
        FOREIGN KEY (`token_type_id`)
            REFERENCES `token_types` (`id`)
            ON DELETE CASCADE
            ON UPDATE CASCADE
) ENGINE = InnoDB;
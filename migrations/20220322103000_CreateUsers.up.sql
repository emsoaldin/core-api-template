CREATE TABLE `users`
(
    `id`                 INT unsigned NOT NULL AUTO_INCREMENT,
    `first_name`         VARCHAR(255) NOT NULL DEFAULT '',
    `last_name`          VARCHAR(255) NOT NULL DEFAULT '',
    `email`              VARCHAR(120) NOT NULL,
    `created_at`         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    `updated_at`         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`         TIMESTAMP    NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `users_email_idx` (`email` ASC)
) ENGINE = InnoDB;

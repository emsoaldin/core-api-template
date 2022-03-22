CREATE TABLE `auth_providers`
(
    `provider`   VARCHAR(50)  NOT NULL,
    `user_id`    INT unsigned NOT NULL,
    `uid`        VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`provider`, `user_id`),
    CONSTRAINT `fk_auth_providers_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB;
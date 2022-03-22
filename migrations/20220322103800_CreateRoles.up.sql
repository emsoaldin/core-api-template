CREATE TABLE `roles`
(
    `id`          INT unsigned NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(45)  NOT NULL,
    `description` VARCHAR(255) NOT NULL DEFAULT '',
    `created_at`  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;
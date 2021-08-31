# Posts
DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts`
(
    `uuid`       varchar(36) NOT NULL,
    `user_uuid`  varchar(36) NOT NULL,
    `content`    longtext    NOT NULL,
    `created_at` timestamp   NOT NULL DEFAULT current_timestamp(),
    `updated_at` timestamp   NULL     DEFAULT NULL ON UPDATE current_timestamp(),
    `deleted_at` timestamp   NULL     DEFAULT NULL,
    UNIQUE KEY `uuid` (`uuid`)
) ENGINE = InnoDB;


# Accounts
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts`
(
    `uuid`          VARCHAR(36)                           NOT NULL,
    `email`         VARCHAR(255)                          NOT NULL,
    `password_hash` VARCHAR(255)                          NOT NULL,
    `created_at`    TIMESTAMP                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP on update CURRENT_TIMESTAMP NULL     DEFAULT NULL,
    `deleted_at`    TIMESTAMP                             NULL     DEFAULT NULL,
    PRIMARY KEY (`uuid`),
    UNIQUE (`email`)
) ENGINE = InnoDB;

# Comments
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `uuid`        VARCHAR(36)                           NOT NULL,
    `user_uuid`   VARCHAR(36)                           NOT NULL,
    `post_uuid`   VARCHAR(36)                           NOT NULL,
    `parent_uuid` VARCHAR(36)                           NOT NULL,
    `content`     LONGTEXT                              NOT NULL,
    `created_at`  TIMESTAMP                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_atd` TIMESTAMP on update CURRENT_TIMESTAMP NULL     DEFAULT NULL,
    `deleted_at`  TIMESTAMP                             NULL     DEFAULT NULL,
    PRIMARY KEY (`uuid`),
    INDEX (`user_uuid`),
    INDEX (`post_uuid`)
) ENGINE = InnoDB;

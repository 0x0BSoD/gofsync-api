CREATE DATABASE `goFsync` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `goFsync`;

CREATE TABLE `environments` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                `foreman_id` int(11) NOT NULL,
                                `host` varchar(255) NOT NULL,
                                `env` varchar(255) NOT NULL,
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `environments_id_uindex` (`id`),
                                KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `hg` (
                      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                      `foreman_id` int(11) NOT NULL,
                      `name` varchar(255) NOT NULL,
                      `host` varchar(255) NOT NULL,
                      `dump` text NOT NULL,
                      `pcList` text NOT NULL,
                      `created_at` datetime NOT NULL,
                      `updated_at` datetime NOT NULL,
                      PRIMARY KEY (`id`),
                      UNIQUE KEY `hg_id_uindex` (`id`),
                      KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `hg_parameters` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `hg_id` int(11) DEFAULT NULL,
                                 `name` varchar(255) DEFAULT NULL,
                                 `value` varchar(255) DEFAULT NULL,
                                 `priority` int(11) DEFAULT NULL,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `hg_parameters_id_uindex` (`id`),
                                 KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `locations` (
                             `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                             `foreman_id` int(11) NOT NULL,
                             `host` varchar(255) NOT NULL,
                             `loc` varchar(255) NOT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `locations_id_uindex` (`id`),
                             KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `override_values` (
                                   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                   `sc_id` int(11) DEFAULT NULL,
                                   `match` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                   `value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
                                   `use_puppet_default` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                   `foreman_id` int(11) DEFAULT NULL,
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `override_values_id_uindex` (`id`),
                                   KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25091 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `puppet_classes` (
                                  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                  `foreman_id` int(11) NOT NULL,
                                  `host` varchar(255) NOT NULL,
                                  `class` varchar(255) NOT NULL,
                                  `subclass` varchar(255) NOT NULL,
                                  `sc_ids` text NOT NULL,
                                  `env_ids` text NOT NULL,
                                  PRIMARY KEY (`id`),
                                  UNIQUE KEY `puppet_classes_id_uindex` (`id`),
                                  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `smart_classes` (
                                 `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                 `foreman_id` INT(11) NOT NULL,
                                 `host` VARCHAR(255) NOT NULL,
                                 `puppetclass` VARCHAR(255) DEFAULT NULL,
                                 `parameter` VARCHAR(255) DEFAULT NULL,
                                 `parameter_type` VARCHAR(255) DEFAULT NULL,
                                 `override_values_count` INT(11) DEFAULT NULL,
                                 `dump` LONGTEXT,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `smart_classes_id_uindex` (`id`),
                                 KEY `id` (`id`)
)  ENGINE=INNODB DEFAULT CHARSET=UTF8MB4 COLLATE = UTF8MB4_UNICODE_CI;

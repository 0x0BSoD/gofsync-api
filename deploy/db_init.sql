DROP database goFsync;
CREATE DATABASE `goFsync` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE goFsync;

CREATE TABLE `hosts` (
                         `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `env` set('stage','prod','error') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'stage',
                         `trend` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `success` int(10) unsigned DEFAULT '0',
                         `failed` int(10) unsigned DEFAULT '0',
                         `rFailed` int(10) unsigned DEFAULT '0',
                         `total` int(10) unsigned DEFAULT '0',
                         `last` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `hosts_id_uindex` (`id`),
                         KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `environments` (
                                `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                `foreman_id` int(11) NOT NULL,
                                `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                                `meta` json DEFAULT NULL,
                                `state` set('ok','outdated','absent','error') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'absent',
                                `repo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'svn://svn.dins.ru/Vportal/trunk/setup/automation/puppet/environments/',
                                `host_id` int(10) unsigned NOT NULL,
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `environments_id_uindex` (`id`),
                                KEY `id` (`id`),
                                KEY `host_id_idx` (`host_id`),
                                CONSTRAINT `environments_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=380 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `hg` (
                      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                      `foreman_id` int(11) NOT NULL,
                      `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                      `dump` text COLLATE utf8mb4_unicode_ci NOT NULL,
                      `pcList` text COLLATE utf8mb4_unicode_ci NOT NULL,
                      `status` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                      `created_at` datetime NOT NULL,
                      `updated_at` datetime NOT NULL,
                      `host_id` int(10) unsigned NOT NULL,
                      PRIMARY KEY (`id`),
                      UNIQUE KEY `hg_id_uindex` (`id`),
                      KEY `id` (`id`),
                      KEY `host_id_idx` (`host_id`),
                      CONSTRAINT `hg_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3561 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `hg_parameters` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `hg_id` int(11) unsigned DEFAULT NULL,
                                 `foreman_id` int(11) NOT NULL,
                                 `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `value` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `priority` int(11) DEFAULT NULL,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `hg_parameters_id_uindex` (`id`),
                                 KEY `id` (`id`),
                                 KEY `hg_id_idx` (`hg_id`),
                                 CONSTRAINT `hg_id` FOREIGN KEY (`hg_id`) REFERENCES `hg` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2539 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `locations` (
                             `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                             `foreman_id` int(11) NOT NULL,
                             `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                             `host_id` int(11) unsigned NOT NULL,
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `locations_id_uindex` (`id`),
                             KEY `id` (`id`),
                             KEY `host_id_idx` (`host_id`),
                             CONSTRAINT `locations_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `smart_classes` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `foreman_id` int(11) NOT NULL,
                                 `puppetclass` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `parameter` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `parameter_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `override` tinyint(1) DEFAULT '0',
                                 `override_values_count` int(11) DEFAULT NULL,
                                 `dump` longtext COLLATE utf8mb4_unicode_ci,
                                 `host_id` int(11) unsigned NOT NULL,
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `smart_classes_id_uindex` (`id`),
                                 KEY `id` (`id`),
                                 KEY `host_id_idx` (`host_id`),
                                 CONSTRAINT `smart_classes_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6269 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `override_values` (
                                   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                   `sc_id` int(11) unsigned DEFAULT NULL,
                                   `match` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                   `value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
                                   `use_puppet_default` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                   `foreman_id` int(11) DEFAULT NULL,
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `override_values_id_uindex` (`id`),
                                   KEY `id` (`id`),
                                   KEY `sc_id_idx` (`sc_id`),
                                   CONSTRAINT `sc_id` FOREIGN KEY (`sc_id`) REFERENCES `smart_classes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8550 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `puppet_classes` (
                                  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                  `foreman_id` int(11) NOT NULL,
                                  `class` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                                  `subclass` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                                  `sc_ids` text COLLATE utf8mb4_unicode_ci NOT NULL,
                                  `env_ids` text COLLATE utf8mb4_unicode_ci NOT NULL,
                                  `host_id` int(11) unsigned NOT NULL,
                                  PRIMARY KEY (`id`),
                                  UNIQUE KEY `puppet_classes_id_uindex` (`id`),
                                  KEY `id` (`id`),
                                  KEY `host_id_idx` (`host_id`),
                                  CONSTRAINT `puppet_classes_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5426 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `goFsync`.`hosts` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`environments` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`hg` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`hg_parameters` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`locations` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`override_values` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`puppet_classes` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`smart_classes` AUTO_INCREMENT=0;


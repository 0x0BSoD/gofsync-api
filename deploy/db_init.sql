DROP database goFsync;
CREATE DATABASE `goFsync` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE goFsync;

CREATE TABLE IF NOT EXISTS `goFsync`.`hosts` (
                                                 `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                 `name` VARCHAR(255) NOT NULL,
                                                 `env` SET('stage', 'prod', 'error') NOT NULL DEFAULT 'stage',
                                                 `trend` VARCHAR(255) NULL DEFAULT NULL,
                                                 `success` INT(10) UNSIGNED NULL DEFAULT '0',
                                                 `failed` INT(10) UNSIGNED NULL DEFAULT '0',
                                                 `rFailed` INT(10) UNSIGNED NULL DEFAULT '0',
                                                 `total` INT(10) UNSIGNED NULL DEFAULT '0',
                                                 `last` VARCHAR(255) NULL DEFAULT '',
                                                 PRIMARY KEY (`id`),
                                                 UNIQUE INDEX `hosts_id_uindex` (`id` ASC) VISIBLE,
                                                 INDEX `id` (`id` ASC) VISIBLE)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`locations` (
                                                     `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                     `foreman_id` INT(11) NOT NULL,
                                                     `name` VARCHAR(255) NOT NULL,
                                                     `host_id` INT(11) UNSIGNED NOT NULL,
                                                     PRIMARY KEY (`id`),
                                                     UNIQUE INDEX `locations_id_uindex` (`id` ASC) VISIBLE,
                                                     INDEX `id` (`id` ASC) VISIBLE,
                                                     INDEX `host_id_idx` (`host_id` ASC) VISIBLE,
                                                     CONSTRAINT `locations_host_id`
                                                         FOREIGN KEY (`host_id`)
                                                             REFERENCES `goFsync`.`hosts` (`id`)
                                                             ON DELETE NO ACTION
                                                             ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`environments` (
                                                        `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                        `foreman_id` INT(11) NOT NULL,
                                                        `name` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL,
                                                        `meta` JSON NULL DEFAULT NULL,
                                                        `state` SET('ok', 'outdated', 'absent', 'error') NOT NULL DEFAULT 'absent',
                                                        `repo` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NOT NULL DEFAULT 'svn://svn.dins.ru/Vportal/trunk/setup/automation/puppet/environments/',
                                                        `host_id` INT(10) UNSIGNED NOT NULL,
                                                        PRIMARY KEY (`id`),
                                                        UNIQUE INDEX `environments_id_uindex` (`id` ASC) VISIBLE,
                                                        INDEX `id` (`id` ASC) VISIBLE,
                                                        INDEX `host_id_idx` (`host_id` ASC) VISIBLE,
                                                        CONSTRAINT `environments_host_id`
                                                            FOREIGN KEY (`host_id`)
                                                                REFERENCES `goFsync`.`hosts` (`id`)
                                                                ON DELETE NO ACTION
                                                                ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`hg` (
                                              `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                              `foreman_id` INT(11) NOT NULL,
                                              `name` VARCHAR(255) NOT NULL,
                                              `dump` TEXT NOT NULL,
                                              `pcList` TEXT NOT NULL,
                                              `status` VARCHAR(255) NULL DEFAULT NULL,
                                              `created_at` DATETIME NOT NULL,
                                              `updated_at` DATETIME NOT NULL,
                                              `host_id` INT(10) UNSIGNED NOT NULL,
                                              PRIMARY KEY (`id`),
                                              UNIQUE INDEX `hg_id_uindex` (`id` ASC) VISIBLE,
                                              INDEX `id` (`id` ASC) VISIBLE,
                                              INDEX `host_id_idx` (`host_id` ASC) VISIBLE,
                                              CONSTRAINT `hg_host_id`
                                                  FOREIGN KEY (`host_id`)
                                                      REFERENCES `goFsync`.`hosts` (`id`)
                                                      ON DELETE NO ACTION
                                                      ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`hg_parameters` (
                                                         `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                         `hg_id` INT(11) UNSIGNED NULL DEFAULT NULL,
                                                         `foreman_id` INT(11) NOT NULL,
                                                         `name` VARCHAR(255) NULL DEFAULT NULL,
                                                         `value` VARCHAR(255) NULL DEFAULT NULL,
                                                         `priority` INT(11) NULL DEFAULT NULL,
                                                         PRIMARY KEY (`id`),
                                                         UNIQUE INDEX `hg_parameters_id_uindex` (`id` ASC) VISIBLE,
                                                         INDEX `id` (`id` ASC) VISIBLE,
                                                         INDEX `hg_id_idx` (`hg_id` ASC) VISIBLE,
                                                         CONSTRAINT `hg_id`
                                                             FOREIGN KEY (`hg_id`)
                                                                 REFERENCES `goFsync`.`hg` (`id`)
                                                                 ON DELETE NO ACTION
                                                                 ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`puppet_classes` (
                                                          `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                          `foreman_id` INT(11) NOT NULL,
                                                          `class` VARCHAR(255) NOT NULL,
                                                          `subclass` VARCHAR(255) NOT NULL,
                                                          `sc_ids` TEXT NOT NULL,
                                                          `env_ids` TEXT NOT NULL,
                                                          `host_id` INT(11) UNSIGNED NOT NULL,
                                                          PRIMARY KEY (`id`),
                                                          UNIQUE INDEX `puppet_classes_id_uindex` (`id` ASC) VISIBLE,
                                                          INDEX `id` (`id` ASC) VISIBLE,
                                                          INDEX `host_id_idx` (`host_id` ASC) VISIBLE,
                                                          CONSTRAINT `puppet_classes_host_id`
                                                              FOREIGN KEY (`host_id`)
                                                                  REFERENCES `goFsync`.`hosts` (`id`)
                                                                  ON DELETE NO ACTION
                                                                  ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`smart_classes` (
                                                         `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                         `foreman_id` INT(11) NOT NULL,
                                                         `puppetclass` VARCHAR(255) NULL DEFAULT NULL,
                                                         `parameter` VARCHAR(255) NULL DEFAULT NULL,
                                                         `parameter_type` VARCHAR(255) NULL DEFAULT NULL,
                                                         `override` TINYINT(1) NULL DEFAULT '0',
                                                         `override_values_count` INT(11) NULL DEFAULT NULL,
                                                         `dump` LONGTEXT NULL DEFAULT NULL,
                                                         `host_id` INT(11) UNSIGNED NOT NULL,
                                                         PRIMARY KEY (`id`),
                                                         UNIQUE INDEX `smart_classes_id_uindex` (`id` ASC) VISIBLE,
                                                         INDEX `id` (`id` ASC) VISIBLE,
                                                         INDEX `host_id_idx` (`host_id` ASC) VISIBLE,
                                                         CONSTRAINT `smart_classes_host_id`
                                                             FOREIGN KEY (`host_id`)
                                                                 REFERENCES `goFsync`.`hosts` (`id`)
                                                                 ON DELETE NO ACTION
                                                                 ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `goFsync`.`override_values` (
                                                           `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                                                           `sc_id` INT(11) UNSIGNED NULL DEFAULT NULL,
                                                           `match` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT NULL,
                                                           `value` LONGTEXT CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT NULL,
                                                           `use_puppet_default` VARCHAR(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT NULL,
                                                           `foreman_id` INT(11) NULL DEFAULT NULL,
                                                           PRIMARY KEY (`id`),
                                                           UNIQUE INDEX `override_values_id_uindex` (`id` ASC) VISIBLE,
                                                           INDEX `id` (`id` ASC) VISIBLE,
                                                           INDEX `sc_id_idx` (`sc_id` ASC) VISIBLE,
                                                           CONSTRAINT `sc_id`
                                                               FOREIGN KEY (`sc_id`)
                                                                   REFERENCES `goFsync`.`smart_classes` (`id`)
                                                                   ON DELETE NO ACTION
                                                                   ON UPDATE NO ACTION)
    ENGINE = InnoDB
    AUTO_INCREMENT = 0
    DEFAULT CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;

ALTER TABLE `goFsync`.`environments` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`hg` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`hg_parameters` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`locations` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`override_values` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`puppet_classes` AUTO_INCREMENT=0;
ALTER TABLE `goFsync`.`smart_classes` AUTO_INCREMENT=0;


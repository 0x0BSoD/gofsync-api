package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func dbActions() {

	if _, err := os.Stat(globConf.DBFile); os.IsNotExist(err) {

		db := *globConf.DB
		_, err = db.Exec(`
# HOST GROUPS ==================================================================
create table hg
(
  id          integer unsigned auto_increment,
  foreman_id integer not null,
  name       varchar(255) NOT NULL,
  host       varchar(255) NOT NULL,
  dump       text         NOT NULL,
  pcList     text not null,
  locList    text not null,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL,
  key(id)
);
create unique index hg_id_uindex on hg (id);
alter table hg add constraint hg_pk primary key (id);

# LOCATIONS ==================================================================
create table locations
(
  id          integer unsigned auto_increment,
  foreman_id integer not null,
  host       varchar(255) NOT NULL,
  loc       varchar(255) NOT NULL,
  key(id)
);
create unique index locations_id_uindex
  on locations (id);
alter table locations
  add constraint locations_pk
    primary key (id);

# ENVIRONMENTS ==================================================================
create table environments
(
  id          integer unsigned auto_increment,
  foreman_id  integer not null,
  host        varchar(255) NOT NULL,
  env         varchar(255) NOT NULL,
  key(id)
);
create unique index environments_id_uindex
  on environments (id);
alter table environments
  add constraint environments_pk
    primary key (id);

# HOST GROUPS Parameters ==================================================================
create table hg_parameters
(
  id       integer unsigned auto_increment,
  hg_id    integer,
  name     varchar(255),
  value    varchar(255),
  priority integer,
  key(id)
);
create unique index hg_parameters_id_uindex
  on hg_parameters (id);
alter table hg_parameters
  add constraint hg_parameters_pk
    primary key (id);

# Puppet classes Parameters ==================================================================
create table puppet_classes
(
  id          integer unsigned auto_increment,
  foreman_id  integer not null,
  host        varchar(255) NOT NULL,
  class       varchar(255) NOT NULL,
  subclass    varchar(255) NOT NULL,
  sc_ids      varchar(255) NOT NULL,
  env_ids     varchar(255) NOT NULL,
  hg_ids      varchar(255) NOT NULL,
  key(id)
);
create unique index puppet_classes_id_uindex
  on puppet_classes (id);
alter table puppet_classes
  add constraint puppet_classes_pk
    primary key (id);

# Smart classes Parameters ==================================================================
create table smart_classes
(
  id                    integer unsigned auto_increment,
  foreman_id            integer not null,
  host                  varchar(255) NOT NULL,
  parameter             varchar(255),
  parameter_type        varchar(255),
  override_values_count integer,
  dump                  varchar(255),
  key(id)
);
create unique index smart_classes_id_uindex
  on smart_classes (id);
alter table smart_classes
  add constraint smart_classes_pk
    primary key (id);

# Override Values ==================================================================
create table override_values(
  id                 integer unsigned auto_increment,
  sc_id              integer,
  'match'            varchar(255),
  value              varchar(255),
  use_puppet_default varchar(255),
  key(id)
);
create unique index override_values_id_uindex
  on override_values (id);
alter table override_values
  add constraint override_values_pk
    primary key (id);
`)
		if err != nil {
			log.Printf("%q\n", err)
			return
		}
	} else {
		//fmt.Println("Base file exist")
	}
}

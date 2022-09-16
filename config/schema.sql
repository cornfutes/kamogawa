create table users
(
    id            bigserial,
    created_at    timestamp with time zone,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone,
    email         text not null,
    password      text,
    gmail         text,
    scope         text,
    access_token  text,
    refresh_token text,
    constraint users_pkey
        primary key (id, email)
);

create index idx_users_deleted_at
    on users (deleted_at);

create table project_auths
(
    gmail      text not null,
    project_id text not null,
    constraint project_auths_pkey
        primary key (gmail, project_id)
);

create table project_dbs
(
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone,
    project_number   text,
    project_id       text not null,
    life_cycle_state text,
    name             text,
    create_time      text,
    has_gce_enabled  boolean,
    constraint project_dbs_pkey
        primary key (project_id)
);

create index idx_project_dbs_deleted_at
    on project_dbs (deleted_at);

create index project_dbs_idx
    on project_dbs using gin (((name || ' '::text) || project_id) gin_trgm_ops);

create table gce_instance_auths
(
    gmail text not null,
    id    text not null,
    constraint gce_instance_auths_pkey
        primary key (gmail, id)
);

create table gce_instance_dbs
(
    id         text not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name       text,
    status     text,
    project_id text,
    zone       text,
    constraint gce_instance_dbs_pkey
        primary key (id)
);

create index idx_gce_instance_dbs_deleted_at
    on gce_instance_dbs (deleted_at);

create index gce_instance_dbs_idx
    on gce_instance_dbs using gin (((((((((id || ' '::text) || name) || ' '::text) || status) || ' '::text) ||
                                      project_id) || ' '::text) || zone) gin_trgm_ops);

create table gae_service_auths
(
    gmail text not null,
    id    text not null,
    constraint gae_service_auths_pkey
        primary key (gmail, id)
);

create table gae_service_dbs
(
    name       text,
    id         text not null,
    project_id text not null,
    constraint gae_service_dbs_pkey
        primary key (id, project_id)
);

create index gae_service_dbs_idx
    on gae_service_dbs using gin (((((name || ' '::text) || id) || ' '::text) || project_id) gin_trgm_ops);

create table gae_version_auths
(
    gmail text not null,
    id    text not null,
    constraint gae_version_auths_pkey
        primary key (gmail, id)
);

create table gae_version_dbs
(
    name           text,
    id             text not null,
    project_id     text,
    service_name   text,
    service_id     text,
    serving_status text,
    constraint gae_version_dbs_pkey
        primary key (id)
);

create index gae_version_dbs_idx
    on gae_version_dbs using gin (((((((((((name || ' '::text) || id) || ' '::text) || project_id) || ' '::text) ||
                                       service_name) || ' '::text) || service_id) || ' '::text) || serving_status)
                                  gin_trgm_ops);

create table gae_instance_auths
(
    gmail text not null,
    id    text not null,
    constraint gae_instance_auths_pkey
        primary key (gmail, id)
);

create table gae_instance_dbs
(
    name         text,
    id           text not null,
    project_id   text,
    service_name text,
    version_name text,
    version_id   text,
    vm_name      text,
    constraint gae_instance_dbs_pkey
        primary key (id)
);


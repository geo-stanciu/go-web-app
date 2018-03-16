create table if not exists currency (
    currency_id int auto_increment primary key,
    currency    varchar(8) not null,
    constraint currency_uk unique (currency)
);

create table if not exists exchange_rate (
    currency_id           int            not null,
    exchange_date         date           not null,       
    rate                  numeric(18, 6) not null,
    reference_currency_id int            not null,
    constraint exchange_rate_pk primary key (currency_id, exchange_date),
    constraint exchange_rate_currency_fk foreign key (currency_id)
        references currency (currency_id),
    constraint exchange_rate_ref_currency_fk foreign key (reference_currency_id)
        references currency (currency_id)
);

create index if not exists idx_exchange_rate_curr_id on exchange_rate (currency_id);
create index if not exists idx_exchange_rate_refcurr_id on exchange_rate (reference_currency_id);
create index if not exists idx_exchange_rate_date on exchange_rate (exchange_date);

create table if not exists audit_log (
    audit_log_id   bigint auto_increment PRIMARY KEY,
    source         varchar(64) not null,
    source_version varchar(16) not null,
    log_time       datetime(3) not null,
    log_msg        JSON not null
);

create index idx_time_audit_log on audit_log (log_time);
CREATE INDEX idx_log_source_audit_log ON audit_log (source);


CREATE TABLE request (
  request_id        int AUTO_INCREMENT PRIMARY KEY,
  request_template  varchar(64)  not null DEFAULT '-',
  request_url       varchar(128) not null DEFAULT '-',
  controller        varchar(64)  not null DEFAULT '-',
  action            varchar(64)  not null DEFAULT '-',
  redirect_url      varchar(256) not null DEFAULT '-',
  redirect_on_error varchar(256) not null DEFAULT '-',
  request_type      varchar(8)   not null DEFAULT 'GET',
  index_level       int,
  order_number      int,
  fire_event        int          not null DEFAULT 1,
  parent_id       int,
  constraint request_url_uk unique (request_url, request_type),
  constraint request_type_chk check (request_type in ('GET', 'POST')),
  constraint request_idx_uk unique (index_level, order_number),
  constraint request_event_chk check (fire_event in (0, 1)),
    constraint request_parent foreign key (parent_id)
        references request (request_id)
);

create index if not exists idx_request_parent on request (parent_id);
);

CREATE TABLE role (
  role_id   int AUTO_INCREMENT PRIMARY KEY,
  role      varchar(64) not null,
  loweredrole varchar(64) not null
);

CREATE UNIQUE INDEX IF NOT EXISTS role_uk ON role (loweredrole);

CREATE TABLE IF NOT EXISTS request_name (
    request_id int NOT NULL,
    language varchar(8) NOT NULL,
    name varchar(64) NOT NULL,
    constraint request_name_pk PRIMARY KEY (request_id, language),
    constraint request_name_fk FOREIGN KEY (request_id)
      REFERENCES request (request_id)
);

create index if not exists idx_request_name_id on request_name (request_id);

CREATE TABLE IF NOT EXISTS request_role (
    role_id int NOT NULL,
    request_id int NOT NULL,
    constraint request_role_pk PRIMARY KEY (role_id, request_id),
    constraint request_role_fk FOREIGN KEY (role_id)
      REFERENCES role (role_id),
    constraint request_role_req_fk FOREIGN KEY (request_id)
      REFERENCES request (request_id)
);

create index if not exists idx_request_role_id on request_role (role_id);
create index if not exists idx_request_role_re1_id on request_role (request_id);

CREATE TABLE user (
  user_id                bigint AUTO_INCREMENT PRIMARY KEY,
  username               varchar(64) not null,
  loweredusername        varchar(64) not null,
  name                   varchar(64) not null,
  surname                varchar(64) not null,
  email                  varchar(64) not null,
  loweredemail           varchar(64) not null,
  creation_time          datetime(3) not null,
  last_update            datetime(3) not null,
  activated              int         not null DEFAULT 0,
  activation_time        datetime(3),
  last_password_change   datetime(3),
  failed_password_atmpts int         not null DEFAULT 0,
  first_failed_password  datetime(3),
  last_failed_password   datetime(3),
  last_connect_time      datetime(3),
  last_connect_ip        varchar(128),
  valid                  int         not null DEFAULT 1,
  locked_out             int         not null DEFAULT 0,
  CONSTRAINT user_uk unique(loweredusername)
);

CREATE TABLE user_password (
  password_id   bigint       AUTO_INCREMENT PRIMARY KEY,
  user_id       bigint       NOT NULL,
  password      VARCHAR(256) NOT NULL,
  password_salt VARCHAR(256) NOT NULL,
  valid_from    datetime(3) NOT NULL,
  valid_until   datetime(3),
  temporary     INT          NOT NULL DEFAULT 0,
  constraint user_password_fk foreign key (user_id)
    references user(user_id)
);

create index if not exists idx_user_password_usr_id on user_password (user_id);

CREATE TABLE user_role (
  user_role_id bigint            AUTO_INCREMENT PRIMARY KEY,
  user_id      bigint not null,
  role_id      int not null,
  valid_from   datetime(3) not null,
  valid_until  datetime(3),
  constraint user_role_fk foreign key (role_id)
    references role(role_id),
  constraint user_role_usr_fk foreign key (user_id)
    references user(user_id)
);

create index if not exists idx_user_role_role_id on user_role (role_id);
create index if not exists idx_user_role_usr_id on user_role (user_id);

CREATE TABLE user_role_history (
  user_role_id bigint PRIMARY KEY,
  user_id      bigint not null,
  role_id      int not null,
  valid_from   datetime(3) not null,
  valid_until  datetime(3),
  constraint user_role_h_fk foreign key (role_id)
    references role(role_id),
  constraint user_role_h_usr_fk foreign key (user_id)
    references user(user_id)
);

create index if not exists idx_user_role_h_role_id on user_role_history (role_id);
create index if not exists idx_user_role_h_usr_id on user_role_history (user_id);

CREATE TABLE user_ip (
  user_ip_id bigint       AUTO_INCREMENT PRIMARY KEY,
  user_id    bigint       NOT NULL,
  ip         varchar(256) NOT NULL,
  constraint user_ip_fk foreign key (user_id)
    references user(user_id)
);

create index if not exists idx_user_ip_usr_id on user_ip (user_id);

CREATE TABLE cookie_encode_key (
  cookie_encode_key_id int auto_increment PRIMARY KEY,
  encode_key           varchar(256) not null,
  valid_from           datetime(3) not null,
  valid_until          datetime(3) not null
);


create table if not exists user (
  id integer primary key autoincrement,
  team_id integer not null,
  account_id integer not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists user_index on user(team_id, id, account_id);

create table if not exists account (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);

create index if not exists user_index on user(id);

create table if not exists access_token(
  id integer primary key autoincrement,
  account_id integer not null unique,
  session_token text not null unique,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists token_index on access_token(session_token);

create table if not exists team (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists team_index on team(id);

create table if not exists borrowing(
  id integer primary key autoincrement,
  user_id integer not null,
  name text not null,
  memo text not null,
  has_return text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists borrowing_index on borrowing(user_id, id);

create table if not exists history(
  id integer primary key autoincrement,
  team_id integer not null,
  notion text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists history_index on history(team_id);

-- insert into team (name, password_hash) values ('pixiv-ios', 'testtesttest');
-- select id as team_id from team where team.name = 'pixiv-ios';
-- 
-- insert into account(password_hash, name) values ('1', 'kameike');
-- insert into account(password_hash, name) values ('2', 'kwzr');
-- insert into account(password_hash, name) values ('3', 'fromatom');
-- insert into account(password_hash, name) values ('4', 'nono');
-- 
-- insert into user (team_id, account_id)
-- select account.id, team.id from team, account where team.name = 'pixiv-ios' and account.password_hash = '1';
-- 
-- insert into user (team_id, account_id)
-- select account.id, team.id from team, account where team.name = 'pixiv-ios' and account.password_hash = '2';
-- 
-- insert into user (team_id, account_id)
-- select account.id, team.id from team, account where team.name = 'pixiv-ios' and account.password_hash = '3';
-- 
-- insert into user (team_id, account_id)
-- select account.id, team.id from team, account where team.name = 'pixiv-ios' and account.password_hash = '4';
-- 
-- 

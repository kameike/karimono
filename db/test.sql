
create table if not exists account_team (
  id integer primary key autoincrement,
  team_id integer not null,
  account_id integer not null,
  created_at text default (datetime('now', 'localtime')),
	unique(team_id, account_id)
);
create index if not exists user_index on account_team(team_id, id, account_id);

create table if not exists account (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);

create table if not exists access_token(
  id integer primary key autoincrement,
  account_id integer not null unique,
  token text not null unique,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists token_index on access_token(token);

create table if not exists team (
  id integer primary key autoincrement,
  name text not null unique,
  password_hash text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists team_index on team(id);

create table if not exists borrowing(
  id integer primary key autoincrement,
	account_id integer not null,
	team_id integer not null,
	hashed_id text not null unique,
  name text not null,
  memo text not null,
  has_return text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists borrowing_index on borrowing(account_id, team_id, hashed_id, id);

create table if not exists history(
  id integer primary key autoincrement,
  team_id integer,
	account_id integer,
	text text not null,
  created_at text default (datetime('now', 'localtime'))
);
create index if not exists history_index on history(team_id);


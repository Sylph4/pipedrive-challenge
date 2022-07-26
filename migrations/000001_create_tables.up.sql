CREATE TABLE IF NOT EXISTS Gists
(
ID         varchar(32) PRIMARY KEY,
user_name    varchar(100) NOT NULL,
created_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS  Users
(
user_name       varchar(100) NOT NULL,
github_api_Key    varchar(40)  NOT NULL,
pipedrive_api_Key varchar(40)  NOT NULL,
pipedrive_user_id  varchar(8) NOT NULL
);

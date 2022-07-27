-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE Gists
(
    ID         varchar(32) PRIMARY KEY,
    user_name    varchar(100) NOT NULL,
    created_at timestamp NOT NULL,
	is_checked bool DEFAULT false
);

CREATE TABLE Users
(
    user_name       varchar(100) NOT NULL,
    github_api_Key    varchar(40)  NOT NULL,
    pipedrive_api_Key varchar(40)  NOT NULL,
    pipedrive_user_id  INTEGER
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS Gists;
DROP TABLE IF EXISTS Users;
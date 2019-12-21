create table users (
    id SMALLINT UNSIGNED ZEROFILL NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username CHAR(10) NOT NULL,
    email CHAR(30) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

create user 'takeru'@'localhost' identified by 'Takeru0219!';
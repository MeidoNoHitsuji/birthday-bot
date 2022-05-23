CREATE TABLE users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NULL,
    birthday DATE
);
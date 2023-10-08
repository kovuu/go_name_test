    CREATE TABLE IF NOT EXISTS Person
    (
        id         serial PRIMARY KEY,
        name       varchar(400) NOT NULL,
        surname    varchar(400) NOT NULL,
        patronymic varchar(400),
        age        int,
        sex varchar(10),
        nationality varchar(100)
    )


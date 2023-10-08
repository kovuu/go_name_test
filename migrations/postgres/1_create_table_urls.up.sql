    CREATE TABLE IF NOT EXISTS Person
    (
        id         serial NOT NULL,
        name       varchar(400) NOT NULL,
        surname    varchar(400) NOT NULL,
        patronymic varchar(400),
        age        int,
        gender varchar(10),
        nationality varchar(100),
        PRIMARY KEY (id, name, surname, age, gender, nationality)
    )


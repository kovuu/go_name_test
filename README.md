
Тестовая программа для работы с информацией о людях
с использованием graphql, rest, kafka и redis



Для запуска окружения для программы запустите следующие команды из корневого каталога
- 
-  docker compose up -d          
-  Запустить загруженные в докер образы(kafka запускать через некоторое время, как прогрузится зукипер) 
- Запустить приложение cmd/person_processing,при запуске которого отработает миграция структуры БД(если понадобится)

Так же присутствует дополнительное приложение add_person-producer для ручного наполнения базы посредством сообщений кафки 


Программа имеет возможность получения и изменяния данных о людях из бд:
--

REST:
- /persons - Get метод(получение всех людей , доступны параметры(limit=10, offset=10, filter=name=Vasya))
-   /persons/<id> - Get метод(получение человека по айди)
-   /persons Post метод сохрание человека в БД(требуется передать в теле весь обьект Person(кроме айди и необязательных полей age, nationality, gender))
-   /persons/<id> - Delete метод удаления человека из базы
-  /persons patch метод для обновления данных о человеке(требуется передать полный обьект в каком виде нужно сохранить в базу)


GraphQL:
- доступны альтернативные варианты использования приложения через GraphQL
-  Persons(limit: Int! = 10, offset: Int! = 0, filter: String! = ""): [Person!]!
-  PersonById(id: Int! = 0): Person!
-  createPerson(person: NewPerson): PersonMutationPayload!
-   deletePerson(id: Int!): PersonMutationPayload!
-   updatePerson(person: PersonInput!): PersonMutationPayload!
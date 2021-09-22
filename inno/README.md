![Основная страница](https://github.com/IVZaytsev/myte-solutions/blob/main/inno/treeview.png?raw=true)

## Описание
Приложение (веб-страница) с иерархией сотрудников предприятия.
База данных - SQLite3 с возможностью переключиться на любую другую
Сделано без использования JS-фреймворков

## Зависимости (requirements.txt)
- django
- virtualenv

## Описание содержимого
**/virtualenv/** - содержит виртуальное окружение для Python. Создаеется и настраивается с помощью последовательности команд:
```
$ python -m venv virtualenv
$ myvenv\Scripts\activate
$ python -m pip install --upgrade pip
$ pip install -r requirements.txt
```
**/www/** - каталог с Django-проектом. Создаем его с помощью команды:
```
$ django-admin.exe startproject www .
```

**/workers/** - каталог приложениея **Workers**. Создаем его с помощью команды:
```
$ python manage.py startapp workers
```
**/workers/fixtures/** - тестовый набор данных для проверки работы приложения
**/workers/static/** - скрипты, стили и иконки
**/workers/templates/** - шаблон страницы со списком сотрудников

## Запуск приложения
Заходим в витруальное окружение и запускаем командами:
```
$ myvenv\Scripts\activate
$ python manage.py runserver
```

Приложение доступно по адресу: **[http://127.0.0.1:8000/](http://127.0.0.1:8000/)**

Логин\пароль Администратора: **root : toor** 

Также доступен адрес для отображения отдельного подразделения:
**[http://127.0.0.1:8000/department/[int]/](http://127.0.0.1:8000/department/1/)**
# micro-blogger - powered by Golang + HTMX!

Note: this is a study project and is only meant to give you ideas or serve
for educational purposes. For this reason, certain things might
be implemented non-optimally or incorrectly.

## how-to
- run with: ```make run-dev or make run```
- build with: ```make build```
- test with: ```make test```
- css,html,js static building with: ```npm run dev```

## features
- minimalistic blogging web-app
- admin panel for managing posts and users
- HTMX to reduce full-page rerenderings
- postgres database
- MVC code structure
- UI components with DaisyUI and tailwind.
- light and dark mode supported.

## tech
- Go language
- Golang's html/template
- HTMX
- Tailwind + DaisyUI + Vercel for static building
- Postgresql

## learning resources
- [Boredstack repository from Anthdm](https://github.com/anthdm/boredstack/)
- [Web Development w/ Google’s Go (golang) Programming Language from Todd McLeod](https://www.udemy.com/course/go-programming-language/)

## images

![Snapshot of homepage on 28-08-2023](2023-08-28_16-09-snapshot.jpg)

## Migrate lib
export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'
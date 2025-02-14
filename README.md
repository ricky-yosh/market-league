![MarketLeague Image Wide](./readme-images/logo_wide.jpg)
# MarketLeague
St. Edward's University Senior Project repository for MarketLeague, a fantasy-football style approach to learning about the stock market. Users will be able to create a portfolio of stocks that act like trading cards, and users will be able to join leagues where they can trade stocks with other players.

## Running MarketLeague Locally
Requirements:
[Docker v27.2.0](https://www.docker.com/products/docker-desktop/)
```sh
./run_docker_dev.sh
```
This will run the `docker-compose.dev.yml` and will run the dev environments for the Angular Frontend and Gin Backend. Both will be updated as changes are made live.

### .env file
In order for MarketLeague to run correctly locally you will need to create a `.env` file with these properties:
```
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=database
JWT_KEY=secretkey
```
Change `user`, `password`, and `database` to appropriate values.

## Running MarketLeague for Production:
Requirements:
[Docker v27.2.0](https://www.docker.com/products/docker-desktop/)
```sh
docker compose -f docker-compose.prod.yml up --build -d
```
This will run the `docker-compose.prod.yml` and will run the prod environments for the Angular Frontend and Gin Backend. The `-d` flag will hide docker logs. If you want to see the logs on the production environment run:
```sh
docker compose -f docker-compose.prod.yml logs -f
```

### Ensuring that websockets are secure (wss) not insecure (ws)
If you try and run the secure websockets (wss) it will not work unless the server is running on a secure connection (https). If you do not have a domain name you can use a self-signed certificate:
```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
-keyout ./ssl/nginx-selfsigned.key \
-out ./ssl/nginx-selfsigned.crt
```
Make sure that the `/ssl` folder that is created is on the same level as the `docker-compose.prod.yml` files:
```
market-league
├── README.md
├── docker-compose.prod.yml
├── market-league-back-end/
├── market-league-front-end/
├── ssl/
└── ...
```

### .env
We also need to make sure that an `.env` file exists on the server that is hosting the production version of the webapp.
```
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=database
JWT_KEY=secretkey
```
Change `user`, `password`, and `database` to appropriate values.

## MarketLeague Roadmap
Projected plan for features and presentations.
![MarketLeague Roadmap](./readme-images/marketleague-roadmap.png)

## Credits & Thanks
Thanks to Liam Molina for creating the logo for MarketLeague!
- [Instagram](https://www.instagram.com/designedbyliamm/)
- [LinkedIn](https://www.linkedin.com/in/liam-molina-ab3211290/)

## Creators:
#### Timothée Pommier
- [LinkedIn](https://www.linkedin.com/in/timoth%C3%A9e-pommier-81749a251/)
- [GitHub](https://github.com/TimotheePommier)
#### Ricky Yoshioka
- [LinkedIn](https://www.linkedin.com/in/r1chard-yoshioka/)
- [GitHub](https://github.com/ricky-yosh)

## Resources
- [unDraw Illustrations](https://undraw.co/illustrations)
- [Docker](https://www.docker.com/)
- [Typescript](https://www.typescriptlang.org/)
- [Angular](https://angular.dev/)
- [Go](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [D2Lang](https://d2lang.com/)
- [Mermaid.js](https://mermaid.js.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Digital Ocean Droplet](https://www.digitalocean.com/)
- [Postman](https://www.postman.com/)
- [Sourcetree](https://www.sourcetreeapp.com/)
- [VSCode](https://code.visualstudio.com/)
- [Slack](https://slack.com/)

## MarketLeague Version
v4.3

![Repeating logo checkered pattern](./readme-images/logo_repeat.jpg)
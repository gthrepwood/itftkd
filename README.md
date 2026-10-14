# ITF Taekwondo Kvíz

Interaktív kvízalkalmazás ITF Taekwondo koreai szókincse tanulásához. A webes felületen választhatsz témakört, és a rendszer a `backend/data/data.json` alapján generál magyar fordítással többszörös választású kérdéseket.

## Tech Stack

- **Backend**: Go + Gin
- **Frontend**: HTML5 + Tailwind CSS + Vanilla JavaScript

## Projektstruktúra

```
questionary1/
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── data/
│   │   ├── questions.go
│   │   └── data.json
│   └── handlers/
│       └── quiz.go
├── frontend/
│   ├── index.html
│   ├── app.js
│   └── styles.css
├── Dockerfile
└── docker-compose.yml
```

## Előfeltételek

- Go 1.23+

## Futtatás

```bash
cd backend
go mod tidy
go run main.go
```

A szerver elérhető: `http://localhost:9000`

## Használat

1. Nyisd meg a böngészőben `http://localhost:9000`
2. Válassz témakört a legördülő menüből
3. Add meg a kívánt kérdésszámot
4. Kattints az **Indítás** gombra
5. Válaszd a helyes fordítást a lehetőségek közül
6. A végeredményt a Kvíz véget ért! képernyőn láthatod

## API végpontok

- `GET /api/topics` - Elérhető témakörök listája
- `GET /api/questions?topic=&count=` - Kérdések lekérése témakör alapján
- `POST /api/quiz/start` - Új kvíz indítása
- `GET /api/quiz/:id/question` - Aktuális kérdés lekérése
- `POST /api/quiz/:id/answer` - Válasz beküldése

## Adatforrás

A kérdések a `backend/data/data.json` fájlból generálódnak. A JSON tartalmazza a koreai szavak, kifejezések, számok, állások, technikák és formagyakorlatok magyar jelentéseit.

## Docker

### Előfeltételek

- Docker
- Docker Compose

### Futtatás Dockerrel

```bash
# Build
docker build -t taekwondo-quiz .

# Run
docker run -d -p 9000:9000 --name taekwondo-quiz taekwondo-quiz
```

### Docker Compose

A projekt tartalmaz egy készenlét `docker-compose.yml` fájlt:

```bash
docker compose up -d
```

A szerver elérhető: `http://localhost:9000`

### Docker-frissítés

```bash
docker compose pull
docker compose up -d
```

## Unraid telepítés

### Előfeltételek

- Unraid 6.12+ 
- A Community Applications (CA) plugin telepítve
- Egy mappingelve könyvtár, pl. `/mnt/user/appdata/taekwondo-quiz`

### Lépések

1. Nyisd meg a **Community Applications**-t az Unraid felületén
2. Keresd meg a **Docker Compose Manager**-t, vagy használd a `UserTemplates` funkciót
3. Helyezd a következő `docker-compose.yml` tartalmat a Compose sablonba:

```yaml
version: '3.8'

services:
  taekwondo-quiz:
    image: ghcr.io/<felhasznalo>/taekwondo-quiz:latest
    container_name: taekwondo-quiz
    restart: unless-stopped
    ports:
      - "9000:9000"
    volumes:
      - /mnt/user/appdata/taekwondo-quiz:/data
```

4. Módosítsd a `<felhasznalo>` részt a saját GitHub felhasználónevedre
5. Kattints az **Apply** gombra

### Image építése és pusholása GitHub Container Registry-ra

```bash
# Bejelentkezés GitHub Container Registry-be
echo $CR_PAT | docker login ghcr.io -u <felhasznalo> --password-stdin

# Build és tag
docker build -t ghcr.io/<felhasznalo>/taekwondo-quiz:latest .

# Push
docker push ghcr.io/<felhasznalo>/taekwondo-quiz:latest
```

### Port megjelenítése

- Ügyelj arra, hogy az **Unraid** dashboardon bejegyrizd a 9000-es portot,
  vagy a `docker-compose.yml` ports szakaszában használd a host oldali fix portot.
- Ha már foglalt a 9000, módosítsd a host oldali portot, pl. `"9001:9000"`

### Frissítés

Az Unraid **Docker** oldalon kattints a **Check for Updates** gombra a `taekwondo-quiz` konténerenél.

## Licensz

MIT

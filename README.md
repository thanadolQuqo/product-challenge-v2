# Product challenge

---
## Start Cockroach Local
1. install
```
brew install cockroachdb/tap/cockroach
```
2. start local db 
    - single cluster startup 
    - for accessing db : port 26257
      - admin dashboard : 8080
```
make db
```

---
## Start Program
```
make run
```

---
program will start up and load config from `.env` file which contain db connection config.


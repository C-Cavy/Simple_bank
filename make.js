const path = require("path");

const dir = path.resolve("./");

// docker
const docker = {
  sqlc: `docker run --rm -v ${dir}:/src -w /src kjconroy/sqlc`,
  "postgres start": "docker start postgres14",
  "postgres stop": "docker stop postgres14",
  postgres: "docker exec -it postgres14",
};

const db = {
  createdb: "jv postgres createdb --username=root --owner=root simple_bank",
  dropdb: "jv postgres dropdb simple_bank",
  migrate:
    "migrate -path db/migration -database postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable",
};

// export
module.exports = {
  ...docker,
  ...db,
};

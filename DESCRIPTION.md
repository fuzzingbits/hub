# Project Description

This project is an experiment is building a full stack web application with a focus on maintainability. Currently it's a platform that I'll be building some tools on in the future.

- Complies into a single [Go](https://golang.org/) binary
- [Nuxt.js](https://nuxtjs.org/) files are embedded into the binary for the front end
- [MariaDB](https://mariadb.org/) for relational data storage
- [MongoDB](https://www.mongodb.com/) for document storage
- Redis for temporary cache
- Enforced Contracts between Go and TypeScript:
    - A file containing all of the TypeScript interfaces needed for the application is generated using the Go matching structs
    - API Routes are defined in Go structs
        - Definitions verified via unit tests
        - Definitions are used to generate a TypeScript Class to be used as the API client
    - These 2 features mean that changes to the Go API will automatically generate the corresponding changes in TypeScript and cause the TypeScript to fail to build if there are breaking changes.

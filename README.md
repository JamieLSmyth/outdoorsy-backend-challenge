# Coding Challenge

An implementation of the Outdoorsy Backend Coding Challenge.

## Running the application
This application is entirely containerized and designed using the [Dev Container](https://containers.dev/) standard and works best in VSCode with the [Dev Container Extension](https://code.visualstudio.com/docs/devcontainers/containers) installed.  You can run the application locally with Go installed on your local machine but it is not the recommended way and so it is not documented here.

The easiest way to run the applicaiton, if you want easy access to debug and run the tests, is to open the root directory in VSCode with the dev container extension installed and clicking "Open in Container" when the prompt pops up. From there you can open a new terminal and run the application as follows.

```sh
go run app.go
```

This will start the applicaiton and expose the port 8080 to the host machine.  You can then use curl or any tool you would like (such as postman) to exercise the API.

You can also run the application in a container without opening it as a dev container simply by running the following on a machine with docker and docker compose installed.

```sh
docker-compose up
```

This will start up the two docker containers, one for postgres and one for the go app and run the application.  As above this will expose the API on port 8080 to the local machine.

## Running the test
If you connect to the application through vscode as a dev container you can run the tests in a vscode termina as follows (assuming you are in the `/workspace` directory)
```sh
go test ./...
```

If you decided to go the standalone container route you can connect to a shell on the app container as follows
```sh
docker exec -it interview-challenge-backend-app-1 /bin/bash
```
:warning: The container name might be different on you machine so if you can't connect use `docker ps` to find the container name.

This should open a bash shell in the `/workspace` directory and then you can run the tests as above.

## Developer Notes
This is not a 100% complete and production ready implementation of the specification.  To keep the implementation to a reasonable time, considering the fact that I haven't worked with Go in a couple years, I took a few short cuts and didn't cover every test case, but should show that I have the capability to work my way around Go and the available libraries.

The structure of the project and some of the styles used in the code are shaped by my years of working in Java and recent working in Python, so it doesn't always follow cannonical Go idioms and I would correct that if I was working on this project for a production release. 

I chose to stick as close to vanilla go as was reasonable.  I did chose to use the [GORM](https://gorm.io/) ORM for postgres access, the [Gin framework](https://gin-gonic.com/) for serving web endpoints and [sqlmock](https://pkg.go.dev/github.com/data-dog/go-sqlmock)  for create some unit tests for the repository.

If I was going to do any clean up I would probably separate the Route Handlers from the routing so the Route handlers could stand on thier own.  I would restruct the testing to be follow more of Goes standards and use a more detailed testing library like [testify](https://pkg.go.dev/github.com/stretchr/testify) for assertions and [Godog](https://github.com/cucumber/godog) for functional tests and probably figure out if there is a BDD framework available for Go that is comparable to Groovy's [Spock framework](https://spockframework.org/) (the best framework ever created).

I did have to make some guess, particularly as it was unclear what options should be available regarding sorting and how the parameter values should be present. For example the sample uses the sort parameter of "price"  but there really is no such thing as a "price" as a sortable thing in the API. I also wasn't sure if every field should be sort able and how much safety was wanted around that sorting.  This area of the code could certianly be more robust.

# Original Specification
This is the original specification as provided by Outdoorsy.

## Functionality
The task is to develop a rentals JSON API that returns a list of rentals that can be filtered, sorted, and paginated. We have included files to create a database of rentals.

Your application should support the following endpoints.

- `/rentals/<RENTAL_ID>` Read one rental endpoint
- `/rentals` Read many (list) rentals endpoint
    - Supported query parameters
        - price_min (number)
        - price_max (number)
        - limit (number)
        - offset (number)
        - ids (comma separated list of rental ids)
        - near (comma separated pair [lat,lng])
        - sort (string)
    - Examples:
        - `rentals?price_min=9000&price_max=75000`
        - `rentals?limit=3&offset=6`
        - `rentals?ids=3,4,5`
        - `rentals?near=33.64,-117.93` // within 100 miles
        - `rentals?sort=price`
        - `rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price`

The rental object JSON in the response should have the following structure:
```json
{
  "id": "int",
  "name": "string",
  "description": "string",
  "type": "string",
  "make": "string",
  "model": "string",
  "year": "int",
  "length": "decimal",
  "sleeps": "int",
  "primary_image_url": "string",
  "price": {
    "day": "int"
  },
  "location": {
    "city": "string",
    "state": "string",
    "zip": "string",
    "country": "string",
    "lat": "decimal",
    "lng": "decimal"
  },
  "user": {
    "id": "int",
    "first_name": "string",
    "last_name": "string"
  }
}
```

## Notes
- Running `docker-compose up` will automatically generate a postgres database and some data to work with. Connect and use this database.
- Write production ready code.
- Please make frequent, and descriptive git commits.
- Use third-party libraries or not; your choice.
- Please use Golang to complete this task.
- Feel free to add functionality as you have time, but the feature described above is the priority.
- Please add tests

## What we're looking for
- The functionality of the project matches the description above
- An ability to think through all potential states
- In the README of the project, describe exactly how to run the application and execute the tests

When complete, please push your code to Github to your own account and send the link to the project or zip the project (including the `.git` directory) and send it back.

Thank you and please ask if you have any questions!

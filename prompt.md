you are an expert golang developer tasked with developing a back-end API for a new application that your company is building.

Your manager has provided you with the following description of the application. It is up to you to create the most appropriate API from the description in the <api_spec> tags.

<api_spec>

# ergracer-api

## overview

ergracer is an application that allows users to racer each other using the concept2 rowing machine. This application is split into a traditional client/server architecture. The api (server) handles requests from all clients, including android, ios, and web users.

## high level elements

- postgresql database. the api should expect a connection string variable to be exposed in the environment that allows it to connect to the database
- The api should be written in golang using the latest and most well documented routing framework
- The api should be built with an extensible middleware layer that allows for detailed logging, authentication, authorization, etc.
- authentication should use a JWT bearer token

## application level elements

- a user can create an account for this application by providing an email, password, and username. The email should be verified, and the username should be unique across all users
- when a user signs in, they receive a JWT to use for further API access
- each user can have a list of friends. They can invite another user to be their friend, and if the other user accepts, both users end up on each other's friends list.
- a user can create a race. A race has a specific distance in meters. A race has a uuid in addition to an auto-increment primary key. The uuid is used to share the race with others who want to join.
- once each user in a race has been marked with "ready" status, a 10 second countdown will begin. When the countdown is complete, the race will enter the "active" state, at which point it will accept incoming data on the distance rowed for each user participating in the race.
- when a user's reported distance meets or exceeds the race's distance, that user should be marked as having completed the race. once all users have completed the race, the race will enter the "finished" state
- when the race enters a finished state, it should automatically compute the pace of each participating user, which is reported as mm:ss per 500 meters
- a user can invite another user to become friends only if they have participated in at least 1 race together
- a user can fetch a list of their past races, which should include details about the race like the distance, the other users who participated, the order of completion, and the time to completion.

</api_spec>

This project should include a dockerfile that can automatically build and package the golang binary into a docker image.

Make sure to include documentation, and think step-by-step before writing any code.

# Purpose
- The `commit` package handles the CRUD lifecycle of sessions within the database.
- The package contains routines such as 
  - `CreateSession` : Creates a session & saves to the psql db.
  - `BeginSession` : Loads a session in from the psql db, attempts to update the state and notifies the 
  **Lock Manager** Microservice & **Notification** Microservice through a Kafka Event & saves the session to a redis db. 
  - `UpdateSession` : Updates a session entry & saves the updated entry to the psql db.
  - `DeleteSession` : Removes session from the database.
  - `EndSesssion` : Removes the session from the redis DB and Notifies the **Lock Manager** & **Notification** Microservices.

# Imperative Designs of Routines 
- Each routine within the package has the following 
  - A query builder. 
  - An Error handler which takes in the original error and returns a more descriptive one using context from the state of the program ( If this is not needed we can just return the normal error as recieved from the server) .

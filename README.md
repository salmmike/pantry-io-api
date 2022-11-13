Go API for pantry-io hosted in heroku.

Methods:

    GET /db:
        - Returns a JSON field containing a unit_id paired with items list
        - The items list contains ones and zeros, one meaning item is in inventory, zero meaning not.
        - Caller needs to have a priori knownledge of what item is in which index.
        - Not having item names defined in the database is a deliberate choise so minimum user data is required.
        - Uses authentication token from GET /create

    POST /db:
        - Add infromation of inventory to database.
        - JSON data field should contain unit_id and items list.

    GET /create:
        - Create a new entry, returns unique authentication token.

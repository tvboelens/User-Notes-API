# Design
- `main`
    - Entrypoint for the app where config is read, DB initialized and the server is started.
- `config`
    - Reading config for DB (and other things?)
- `controllers`
    - Business logic for HTTP requests
        - Handling incoming requests
        - Talking to models
        - Returning responses
- `middlewares`
    - Authentication
    - Rate limiting?
    - Logging
    - Error handling
- `models`
    - Define the structures for DB, i.e. ORM mapping
- `repositories`
    - Repository layer on top of models to prevent direct DB access.
- `routes`
    - Register routes
    - Connect URL to controllers and applies middleware
- `services`
    - handle login and registration
    - handle notes
- `utils`
    - Helper functions
    - JWT generation/verification, password hashing, input validation
- `migrations`
    - Creating tables and updating schema.
    - Might be enough to call Automigrate at first
- `tests`
    - These should mostly be integration and end-to-end tests, the unit tests will stay in the packages itself for now.
- perhaps `.env`, `docker-compose.yml` and Dockerfile.

## Todo
- [] Testing for Config
    - [] Missing values?
- [x] Testing for User model
    - [x] Creation
    - [x] Unique username (cannot create second user with same name)
    - [x] Reading out user with specific name
    - [x] Update
    - [x] Delete
    - [x] Preload Notes when fetching user
- [x] Testing for Notes model
    - [x] CRUD
    - [x] List notes by user, maybe get User from Note
    - [x] Try to create note for non existing user (may depend on the DB?)
- [] Repository + testing
    - [] CRUD User
        - [x] Find by ID and name -> return user object
        - [x] Create User -> should return the User object or modify in place
            - Which fields necessary?
        - [] Update User
            - [] Make sure the User has all fields defined
            - [] Error if the user not found
            - Should only modify DB, not the User object -> Not quite true, gorm.Model has an UpdatedAt field
        - [x] Cascading delete for notes
            - [x] Implemented, but needs to be tested
        - [x] When deleting user modify the name so that it can be reused
    - [] CRUD Notes
        - [x] Find by Id
        - [x] Create same as user
        - [] Update same as user
        - [x] Delete by id
    - [] Operations for multiple objects, i.e. allow arrays/slices of Users and Notes?
- [] Error messaging
- [] utils
    - [] encode hash string
    - [] parse hash string
- [x] auth
    - [x] Custom error when user not found
    - [x] testing
        - [x] Login 
            - [x] Succesful login with the right credentials
            - [x] Credentials wrong then login fails
                - [x] Wrong pwd
                - [x] username missing
        - [x] Registration
            - [x] Login works after registration
            - [x] Before registration cannot login -> user not found error
            



## Next steps
1. Repository layer → encapsulate DB CRUD operations.
    1. Should probably include the DB by dependency injection
    2. maybe also need an init to either create the schema or build connection?
    2. Create update, delete users
    2. Create update, delete notes
    3. List notes for single user
    4. Find user by name
2. Utils → hashing & JWT.
    1. Pwd hashing
        1. Argon2 and scrypt seem to be the best choices, bcrypt not so
        2. Use salt?
        3. Consider practicality between security and performance (iterations for hashing)
        4. Probably only need two functions: generating the hash and verifying the password
            1.  Use iterations as parameter (and other params?)
        5. It seems best to have an interface with methods to generate a salt, a password and compare
    3. JWT
        1. Need a secret in order to create JWT
        2. create header and payload
        3. sign the jwt
        4. Verifying jwt
        5. Parsing?
        4. Transmitting comes later
        6. Fields
            1. 
3. Auth service → register/login business logic.
4. Note service → CRUD & ownership rules.
5. Middleware → JWT auth for requests.
6. Controllers & Routes → HTTP layer.

### Services
1. Credentials struct (already defined in utils or auth)
2. Login service
    1. Should contain user repository and password comparer via DI (use constructor?)
        1. Have two interfaces for different parts of the repository, could inject via pointers
        2. Same for password comparer
    2. Call auth.LoginUser
    5. Generate JWT and return it
3. Registration service
    1. Inject pwd hasher and user repo (use constructor?).
    2. Check if user already exists
    2. Call utils.RegisterUser
    3. Generate JWT and return it
4. Note handling service
    1. Fetch note ids, maybe with titles (usecase would be showing a list of notes by title + some metadata)
    2. Fetch note contents by id
    3. Save a new note
    4. Edit (i.e. update) a note
    5. Delete a note
    6. Block a user from accessing notes of another user

### Middleware
1. auth middleware
    1. Check if auth header and bearer present
    2. Check the jwt validity and call next handler
    3. Use gin context to pass claims to next handlers
    4. If anything fails return status unauthorized



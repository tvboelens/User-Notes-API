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

## Stories
- user registers
    - if everything ok user receives the jwt token and can start creating notes and do everything else
    - if something is wrong with the password user receives an error that password must have certain form
    - if username already exists user receives error with message that username already exists
- user logs in
    - if credentials correct user receives jwt token and can proceed
    - if credentials incorrect user receives status unauthorized
    - if username wrong/unknown user receives corresponding error
- user wants to change password
    - here needs to send old and new password
    - login first using old password, then update password
    - can either decide to let user login with new password or already logged in
- changing username (if we want to allow it) using the same logic
- Everything else always needs succesful login (non-expired jwt token)
- User receives a list of all notes (maybe only ids) and as a next step can query a specific note (by id) to receive title and content
    - these can only be their own notes
    - if user tries to query a note that does not belong to them they get status unauthorized
- User creates a note and after sending the request receives the id for later use
- User queries a note, edits it and saves it via id

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
        - [x] implement and test FindNotesByUserId
        - [x] Create same as user
        - [] Update same as user
        - [x] Delete by id
    - [] Operations for multiple objects, i.e. allow arrays/slices of Users and Notes?
- [] Error messaging
    - [] custom error if username already exists
    - [] custom error if password is empty (or handle it directly in the controller)
- [] controllers
    - [] password validation when registering?
    - 
- [] services
    - [x] auth service
    - [] note service
        - [] test read functions
            - [x] GetNote
            - [] GetNotes
        - [] modifying functions
            - [x] CreateNote
            - [] Update
            - [] Delete
- [x] utils
    - [x] encode hash string
    - [x] parse hash string
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
    - [x] jwt
        - [x] encode user id in jwt?
- [] controllers
    - [x] auth controllers
    - [] note controllers
        - [x] create note
        - [] GET all notes
        - [] GET specific note
        - [] POST create note
        - [] DELETE note
        - [] PUT update note
- testing
    - [] e2e
        - [] Register/login/create flow
        - [] two users, try to get/edit note of other user
    - [] integration
        - [] services and controllers
        - [] services and repos
        - [] auth service and registration/login manager?
        - [] repos and models
            



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
        6. Create the hash string. This will be stored in the db.
    3. JWT -> I don't think this needs to be in the utils, creation can be done in the login controller/handler and verification can be done in middleware.
        1. Need a secret in order to create JWT
        2. create header and payload
            1. RFC 7519 gives the spec
            2. Claims: 
                1. Issuer (set name via config, can use something like `auth.user-notes-api.local`)
                2. sub
                3. exp
                3. iat
                4. jti? Seems like a good idea, possibly there are scenarios where revocation is necessary.
                5. role? permissions? I would say not necessary, since every user will only be able to edit their own notes and there will be no admin.
                7. session id? Not now, but maybe add this later.
        3. sign the jwt
        4. Verifying jwt
        5. Parsing?
        4. Transmitting comes later
        6. Fields
            1. 
3. Auth service → register/login business logic.
4. Note service → CRUD & ownership rules.
5. Middleware
    1. Authentication
    2. Logging
    3. Rate limiting?
    4. Panic/error recovery (handle service/handler panics/errors)
    6. Metrics?
6. Controllers & Routes → HTTP layer.




# Ama'Shan

World of Warcraft utilities or something. I basically just get bored and fill my time by working on random projects. 
This is one of them.

---

## Configuration

### Flags

- `addr` sets the string address to bind to
  - defaults to `:9090`
- `denv` is a boolean, when true attempts to read a `.env` file in running directory
  - defaults to `false`

### Environment

There are a few environment variables that we work with. These can be set in your environment or through a `.env` file 
and set auto-magically with the `-denv` runtime flag.

We'll include an example `.env` file below, followed by information regarding each of the variables with links for 
obtaining any relevant access credentials.

```shell
# BattleNet
BNET_CLIENT_ID="<id>"
BNET_CLIENT_SECRET="<secret>"
BNET_REDIRECT_URL="<callback_url>"

# Session
SESSION_KEY="<your_session_key>"
```

---

#### BattleNet

Otherwise, referred to as BNET is Blizzard's API for all things Blizzard (e.g. World of Warcraft). We require an API
Client which can be [managed here](https://develop.battle.net/access/clients).

We redirect to the following OAuth2 URLs:

```text
https://amashan.com/api/auth/battlenet/callback
http://localhost:9090/api/auth/battlenet/callback
```

#### Session

This is the value that will be used for the CookieStore.

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux)
- [hashicorp/go-hclog](https://github.com/hashicorp/go-hclog)
- [stretchr/testify](https://github.com/stretchr/testify)
- [joho/godotenv](https://github.com/joho/godotenv)
- [x/oauth2](https://github.com/golang/oauth2)
- [x/time/rate](https://cs.opensource.google/go/x/time)

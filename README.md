# Lemonoid

# Preconditions
The following environment variables are expected at startup:
| Environment variable | Description |
|-|-|
| LEMONOID_TOKEN_TELEGRAM | Telegram bot token |
| LEMONOID_TOKEN_OPENAI | OpenAi token |
| LEMONOID_TOKEN_ELEVENLABS | ElevenLabs token |
| LEMONOID_TOKEN_DISCORD | Discord token |

# Running
Run the application:

```bash
go run ./...
```

# Testing
Run all tests:

```bash
go test ./...
```

# TODO
- [Persistance with Redis](https://redis.io/docs/latest/develop/connect/clients/go/)
- Use better Telegram bot client (that supports topic and more modern features)
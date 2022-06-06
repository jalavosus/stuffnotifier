# stuffnotifier

A thing which can notify you of:
- cryptocurrency price movements, limits, and more with Gemini's API
- Flight tracking and status via FlightAware's AeroAPI

via Discord, SMS, and potentially other forms of communication.

## API Keys

Keys can be set in your shell's environment (ex. via a `.env` file), or passed as command line flags.

Mapping of keys to environment variables:

|         Key         | Description                         | Environment Variable  |
|:-------------------:|-------------------------------------|:---------------------:|
|   Gemini API Key    | API key for Gemini                  |   `GEMINI_API_KEY`    |
|  Gemini API Secret  | API Secret for Gemini               |  `GEMINI_API_SECRET`  |
| FlightAware API Key | API key for FlightAware Aero API    | `FLIGHTAWARE_API_KEY` |
| Twilio Account SID  | Twilio Account SID (or API key SID) | `TWILIO_ACCOUNT_SID`  |
|   Twilio API Key    | Twilio API key (for SMS)            |   `TWILIO_API_KEY`    |
|  Twilio API Secret  | Twilio API secret                   |  `TWILIO_API_SECRET`  |
|  Twilio Auth Token  | Twilio API Auth token               |  `TWILIO_API_TOKEN`   |
|       Discord       | Discord Bot Token                   |    `DISCORD_TOKEN`    |

## Supported notification methods

- [x] CLI
- [x] SMS
- [ ] Discord
- [ ] Email
- [ ] [Avian Carrier](https://datatracker.ietf.org/doc/html/rfc1149)

## TODO

- [ ] CLI (sorta done)
- [ ] Gemini
  - [x] Rest API integration
  - [ ] Websocket API integration
- [x] FlightAware integration (Flights, Airports)
- [ ] Discord integration
- [x] Twilio integration
- [ ] Email integration
- [ ] REST API service
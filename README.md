# Claptrap [WIP]
[![Build Status](https://travis-ci.org/DSchalla/Claptrap.svg?branch=master)](https://travis-ci.org/DSchalla/Claptrap)
[![codecov](https://codecov.io/gh/DSchalla/Claptrap/branch/master/graph/badge.svg)](https://codecov.io/gh/DSchalla/Claptrap)

Claptrap is a rule-based bot engine for the Mattermost Chat platform. It allows you to define rules for various events
and react upon them. Also, it is possible to implement custom responses to extend it with custom API integrations.

The following conditions and responses are supported:

**Conditions:**

* Text Starts With
* Text Equals
* Text Matches (RegEx)
* Random
* User Equals (User/Actor)
* User Is Role (User/Actor)
* Channel Equals
* Channel Is Type

Multiple conditions are supported, currently all conditions have to be met to trigger the responses (WIP)

**Responses:**

* Message (Public/DM)
* Invite
* Kick
* Delete Message
* Custom Callback Function

Currently Claptrap requires an system administrator account to execute delete/kick/invite actions. It is possible that
Claptrap joins all public channels automatically on system startup, for private groups an invite is required. DM support
is right now not given (unless direct chat with Claptrap itself).

## Installation

The simplest way to get Claptrap up and running will be to download the binary releases from Github. 
After downloading the release for your platform, place it in a directory and create a config file for
it (See [Configuration](#Configuration)). Afterwards you can start claptrap:

```
./claptrap -config_file=config.yaml
```

The process does not fork and logs to stdout in the current version.

## Configuration

The application is configured using a single YAML file with the following parameters:
```
general:
  auto_join_all_channel: true [true/false]
  case_dir: cases/ [directory path]

mattermost:
  api_url: localhost:5000 [domain/ip:port]
  username: claptrap@example.com [email]
  password: hunter2 [password]
  team: nsf [mattermost team name]
```

There are various message types the bot can react on, currently implemented:

* message
* user_add
* user_removed

Rules for the message types are defined in JSON files in the `case_dir`, e.g. `cases/message.json`:

```
[{
  "name": "Regexp Message",
  "conditions": [
    {"type": "text_matches", "condition": "^a[0-9]b$"}
  ],
  "responses": [
    {"action": "message_channel", "message": "Yes, Regexp works!"}
  ]
},
{
    "name": "Debug",
    "conditions": [
      {"type": "text_equals", "condition": "!debug"}
    ],
    "responses": [
      {"action": "kick_user"},
      {"action": "delete_message"}
    ]
  }
  ]
```

**TODO:** Add Table for condition / response per message type supported and their parameters.

## Development

If you are interested in extending or supporting development of Claptrap, create a fork and download it to your local
Gopath. If you are interested in extending Claptrap with own callbacks for cases, you can modify `cmd/claptrap/claptrap.go`
and extend it with your own functions, e.g.:

```
func main() {
    // ...
	lmgtfyCase := rules.Case{
		Name: "LMGTFY",
		Conditions: []rules.Condition{
			rules.TextStartsWithCondition{Condition:"!lmgtfy"},
		},
		ResponseFunc: lmgtfyCaseCallback,
	}
	botServer.AddCase("message", lmgtfyCase)
	botServer.Start()
}

func lmgtfyCaseCallback(event provider.Event, p provider.Provider) bool {
	mattermostHandler := p.(*provider.Mattermost)
	message := strings.Replace(event.Text, "!lmgtfy ", "", 1)
	message = "http://lmgtfy.com/?q=" + strings.Replace(message, " ", "+", -1)
	mattermostHandler.MessagePublic(event.ChannelID, message)
	return true
}

```

**TODO:** Add dependency management

### Testing your changes

The test suite is still work in progress, right now the only major parts that are tested in-depth are the test conditions.
Tests can be invoked with:

```
go test ./...
```

## Deployment

When deploying Claptrap in a production environment, it is recommended to configure it as a system service. Claptrap
can be running either on an independent machine or on the Mattermost host itself. An example configuration for systemd
can be found here: [claptrap.service](contrib/systemd/claptrap.service)

The service expects that the Claptrap binary is placed in `/opt/claptrap/claptrap` and a dedicated `claptrap` user will
be created.


## Authors

* Daniel Schalla

See also the list of [contributors](https://github.com/dschalla/claptrap/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

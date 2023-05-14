# Internal

This document explains the internal workings of SpectroMate. This document will be technical and intended for application authors and contributors. 

The following topics will be covered in detail.

- [Overview](#overview)
- [API Server](#api-server)
- [Routes](#routes)
    - [Health](#health)
    - [Commands](#commands)
    - [Actions](#actions)
- [Cache](#cache)


# Overview

SpectroMate was designed to provide consumers flexibility, minimum maintainability, and scale. These principles are why SpectroMate was designed using Go and structured as an API server.
 
The Go language lends itself well to creating applications built with concurrency while leveraging the positive performance attributes of the language. Additionally, using Go, a strongly statically typed language, errors are detected earlier in the development cycle, and runtime is optimized thanks to the compiler.

The application is structured as an API server using the standard library HTTP package. The decision for structuring the application as an API server is to support the flexibility principle and enable consumers to add new capabilities to SpectroMate.  Although SpectroMate is a great fit for Slack bot purposes, consumers of SpectroMate could add other capabilities to SpectroMate by adding new routes and creating the logic for those routes. 

For example, a consumer could create a new route that is used to generate an on-demand report specific to an internal business process, such as creating an inventory list or activating an internal process, such as adding a user to a platform or tool.

SpectroMate is compiled and distributed as a multi-platform binary. The binary can be installed in a system and start-up without requiring the installation of software dependencies. SpectroMate is also distributed as a container image. The container image is the preferred consumption method as it lends itself nicely to modern infrastructure platforms supporting the deployment of containerized workloads. 

# API Server

SpectroMate's entry point is found in the **main.go** file. The API server is initialized using the `init()` function. The init function sets up the cache network connection, and it's also used to gather all environment variables applicable to the application, such as the cache connection URL or the log output level.

## Environment Variables

The following environment variables are available:

| Variable | Description | Required | Default |
|---|---|---|---|
| `TRACE`| Set the debug level output. Available values are `INFO`, `DEBUG`, `TRACE`. | No| `INFO`|
| `SLACK_SIGNING_SECRET` | The Slack application has a unique signing secret. This value is used to validate the request is originating from the Slack application. | Yes | `""`|
| `MENDABLE_API_KEY` | The client API used to authenticate with the Mendable API. | Yes| `""`|
| `REDIS_TLS`| Enable to require TLS when communicating with the Redis server| No | `false`|
| `PORT` | Specify the network port for the SpectroMate server to listen on.| No| `3000`|
| `HOST`| Specify the network interface the SpectroMate server should listen on. | No | `0.0.0.0`|
| `REDIS_URL` | The URL of the Redis server.| No| `localhost`|
| `REDIS_PASSWORD`| The password of the Redis user.| No | `""`|
| `REDIS_USER`| The username of the Redis user to use for all Redis interactions.| No| `""`|

In the `main()` function, the HTTP server is started by using the `http.ListenAndServe()` function. Before starting the HTTP server, all routes and their respective handler are declared and added to the API server. 

In the following code snippet, three routes are declared. The endpoints are `/health` , `/slack`, `/slack/actions`. 

```go
    healthRoute := endpoints.NewHealthHandlerContext(ctx)
    slackRoute := endpoints.NewSlackHandlerContext(ctx, globalSigningSecret, globalMendableAPIKey, rdb)
    slackActionsRoute := endpoints.NewActionsHandlerContext(ctx, globalSigningSecret, globalMendableAPIKey)

    http.HandleFunc(internal.ApiPrefixV1+"health", healthRoute.HealthHTTPHandler)
    http.HandleFunc(internal.ApiPrefixV1+"slack", slackRoute.SlackHTTPHandler)
    http.HandleFunc(internal.ApiPrefixV1+"slack/actions", slackActionsRoute.ActionsHTTPHandler)
```

The `http.HandlerFunc` for an endpoint accepts a unique type representing the route. This type is in the [endpoint package](../endpoint/).


Ideally, the current route-type approach should be converted to an interface approach. However, to release an application with SpectroMate's capabilities more quickly, a simpler approach with a type approach for each route was used.

Each route type is defined in the file **endpoints/types.go**. The `/slack` endpoint route type is displayed in the snippet below. The route type contains all the required dependencies for the route.

```go
type SlackRoute struct {
    ctx            context.Context
    signingSecret  string
    mendableApiKey string
    SlackEvent     *internal.SlackEvent
    cache          internal.Cache
}
``` 

If you are creating a new route, create a new type struct. Specify all the required dependencies your route will need.

Example:

```go
type MyNewRouteExample struct {
    ctx            context.Context
    payload        *myPayloadStruct
    anotherAPIKey  string
}
```

# Routes

This section provides an overview of each of the available routes in SpectroMate. 

## Health 

Endpoint: `/health`

You can use the health route to monitor the availability of the application. The route accepts GET requests without any parameters. The return payload is a 200 HTTP status reply.

This route does not perform Slack secret validation as the route is intended for non-Slack purposes. HTTP requests not of the type GET will return a 405 HTTP status code.

## Commands

Endpoint: `/slack`

Slack commands are accepted through the Slack route. The Slack route takes HTTP POST requests from all domains. If the request is not of the type POST an HTTP error response with the 405 HTTP status code is returned. 

The Slack route requires validation of the [Slack signature secret](https://api.slack.com/authentication/verifying-requests-from-slack). Validation failures return an error and a 403 HTTP status code.

The Slack payload for the command is converted to a Go struct, and the type is `SlackEvent`, which is defined in **internal/types.go**

Once the HTTP method and the Slack signature are verified, and the  Slack payload is marshaled, the route `getHandler()` function is invoked. 

The route handler function is a common pattern across all routes, and it's where you can find the core logic of all routes. The Slack handler function acts as the entry point for all Slack commands, as all commands are defined in this route.

The command route has some string logic to extract the Slack sub-command. For example, users can define multiple Slack commands such as `/docs ask` or `/docs pask` by using the second argument as the command identifies. The core command is the application's entry point, and the sub-commands are how additional functionality is exposed. 

Once the sub-command is extracted from the Slack payload, the command string value is validated to ensure a valid command was provided. The return value is a type value of the Slack sub-command.

```go
    cmd, err := determineCommand(userCmd)
    if err != nil {
        log.Debug().Msg("Error converting string to SlackCommands type.")
    }
```

You must define each command as a type. If you are adding new Slack command, define the new command in the **endpoints/type.go** file. Start by adding a new const value for your route.

For example, assume you are adding a new command titled "coffee". The first step is to add a const titled `Coffee`

```go
type SlackCommands int

const (
    Help SlackCommands = iota
    Ask
    PAsk
    Coffee
)
```
Next, update the SlackCommands's string function to convert the type to a string value and vice versa.

```go 
// String converts a SlackCommands type to a string.
func (s SlackCommands) String() string {
    switch s {
    case Help:
        return "help"
    case Ask:
        return "ask"
    case PAsk:
        return "pask"
   case Coffee:
       return "coffee"
    default:
        return "unknown"
    }
}

// SlackCommandsFromString converts a string to a SlackCommands type.
func SlackCommandsFromString(s string) (SlackCommands, error) {
    switch s {
    case "help":
        return Help, nil
    case "ask":
        return Ask, nil
    case "pask":
        return PAsk, nil
  case "coffee":
      return Coffee, nil
    default:
        return -1, fmt.Errorf("unknown command: %s", s)
    }
}
```
In the file **endpoints/slack.go**. The Slack endpoint's switch statement routes each incoming command to the correct case logic. Take the Ask command as an example. The Slack request for AskCommand is prepared, followed by a reply to the Slack server to address the three-second timeout requirement. Lastly, the logic for the Ask command is invoked through a concurrent function call. 

```go
    switch cmd {
    .... //Abbreviated code
    case Ask:
        slackRequestInfo := slackCmds.NewSlackAskRequest(
            slack.ctx,
            Slack.SlackEvent,
            slack.mendableApiKey,
            Slack.cache,
        )
        reply200Payload, err := internal.ReplyStatus200(slack.SlackEvent.ResponseURL, writer, false)
        if err != nil {
            log.Info().Err(err).Msg("failed to reply to slack with status 200.")
            return nil, err
        }
        // Reply back to Slack with a 200 status code to avoid the 3 second timeout.
        returnPayload = reply200Payload
        // Start Go routine to call the command function.
        go slackCmds.AskCmd(slackRequestInfo, false)
    case ... //Abbreviated code
}
```

The `slackCmds.AskCmd()` function is invoked to start a Go routine so that the logic required for the command can continue without being limited to the current request-reply, which is used to reply with an HTTP status code of 200 to address the Slack timeout requirement. This design also allows multiple requests to be handled by the available CPU cores in the system to improve performance. 

If you add a new slack command, add a new case to the switch statement and handle the logic accordingly.

```go
    switch cmd {
    case Coffee:
        slackRequestInfo := slackCmds.NewSlackCoffeeRequest(
            slack.ctx,
            slack.SlackEvent,
        )

        reply200Payload, err := internal.ReplyStatus200(slack.SlackEvent.ResponseURL, writer, false)
        if err != nil {
            log.Info().Err(err).Msg("failed to reply to slack with status 200.")
            return nil, err
        }
        returnPayload = reply200Payload
        // Start Go routine to call the command function.
        go slackCmds.CoffeeCmd(slackRequestInfo)
    case ... //Abbreviated code
    }
```

Notice how the `CoffeeCmd()` function and the other commands are sourced from the `slackCmds` package. The `slackCmds` package is sourced from the [**slackCmds**](../slackCmds/) folder, containing the core logic for each command. All new commands must have their own logic file in the **slackCmds** folder.

# Actions

Endpoint: `/slack/actions/`

The actions endpoint supports Slack [application interactions](https://api.slack.com/interactivity#responses). The actions endpoint accepts HTTP POST requests and requires Slack signature secret verification. 

The actions route handler is located in the **endpoints/slack-actions.go**. The internal route handler uses the action identifier to route the request to the appropriate action logic function. 

```go
// getHandler invokes the modelFeedbackHandler function from the slackActions package
func (actions *ActionsRoute) getHandler(routeRequest *ActionsRoute, reqeust *http.Request, action *internal.SlackActionEvent) ([]byte, error) {
    var returnPayload []byte

    slackRequestInfo := slackActions.NewSlackActionFeedback(routeRequest.ctx, action, routeRequest.mendableApiKey)

    switch action.Actions[0].ActionID {

    case internal.ActionsAskModelPositiveFeedbackID:
        log.Debug().Msg("Positive feedback action triggered.")
        go slackActions.ModelFeedbackHandler(slackRequestInfo, internal.PositiveFeedbackScore)
    case internal.ActionsAskModelNegativeFeedbackID:
        log.Debug().Msg("Negative feedback action triggered.")
        go slackActions.ModelFeedbackHandler(slackRequestInfo, internal.NegativeFeedbackScore)
    default:
        log.Debug().Msg("Unknown action.")
    }

    return returnPayload, nil

}
```
The action identifier is defined in the **internal/constants.go** file.
```go
    ActionsAskModelPositiveFeedbackID string = "ask_model_positive_feedback"
    // ActionsAskModelNegativeFeedbackID is the ID for the negative feedback action.
    ActionsAskModelNegativeFeedbackID string = "ask_model_negative_feedback"
```
The action identifier is an application-defined value that can be applied to a Slack message. For example, the Mendable ask, and pask command includes the action identifier in the return message and embeds the ID in the message's buttons. When a Slack user clicks on the feedback button, the respective action identifier is included in the Slack action event payload. 

To create a new action handler, create a new action logic file in the **slackActions** folder. In the new route, ensure to create an action route type.

For example, if creating an action called "coffeeRating", create a new action type.

```go
type SlackActionCoffeeRating struct {
	ctx            context.Context
	action         *internal.SlackActionEvent
}
```

An action type also requires an action handler. The action handler is where the core logic of the action resides.

```go
func CoffeeRatingHandler(action *SlackActionFeedback, ratingScore internal.MendableRatingScore) {
    ... # Your logic here
}
```

# Cache

The `Cache` interface provides an abstraction layer over the underlying cache technology. The Cache interface defines a contract for a cache system, made up of the following methods:


- `StoreHashMap`: This method is intended to store a hash map in the cache system. It accepts a context, a primary key, and a hash map as inputs.

- `GetHashMap`: This method is designed to retrieve a hash map from the cache system. It returns a boolean indicating the existence of the key, the hash map associated with the key, and any error that might occur during the operation.

- `ExpireKey`: This method sets an expiration time on a specific key in the cache system. If there's an error setting the expiration, it will be returned.

- `Ping`: This method checks the connectivity to the cache system and returns an error if there's any issue.

The `RedisCache` type is the default cache provider supported out-of-the-box but you can swap out the cache provider by creating your cache type that complies with the requirements of the `Cache` interface.


```go
type RedisCache struct {
	redis *redis.Client
}
```
# Internal

This document explains the internal workings of SpectroMate. This document will be technical and intended for application authors and contributors. 

The following topics will be covered in detail.

- [Overview](#overview)
- [API Server](#api-server)
- Routes
    -  Health
    - Commands
    - Actions
    - Cache
- Go Types


# Overview

SpectroMate was designed to provide consumers flexibility, minimum maintainability, and scale. These principles are why SpectroMate was designed using Go and structured as an API server.
 
The Go language lends itself well to creating applications built with concurrency while leveraging the positive performance attributes of the language. Additionally, using Go,  a strongly statically typed language, errors are detected earlier in the development cycle, and runtime is optimized thanks to the compiler.

The application is structured as an API server using the standard library HTTP package. The decision for structuring the application as an API server is to support the flexibility principle and enable consumers to add new capabilities to SpectroMate.  Although SpectroMate is a great fit for Slack bot purposes, consumers of SpectroMate could add other capabilities to SpectroMate by adding new routes and creating the logic for those routes. 

For example, a consumer could create a new route that is used to generate an on-demand report specific to an internal business process, such as creating an inventory list or activating an internal process, such as adding a user to a platform or tool.

SpectroMate is compiled and distributed as a multi-platform binary.

SpectroMate is compiled and distributed as a multi-platform binary. The binary can be installed in a system and start-up without requiring the installation of software dependencies. SpectroMate is also distributed as a container image. The container image is the preferred consumption method as it lends itself nicely to modern infrastructure platforms supporting the deployment of containerized workloads. 

# API Server

SpectroMate's entry point is found in the **main.go** file. The API server is initialized using the `init()` function. The init function sets up the cache network connection, and it's also used to gather all environment variables applicable to the application, such as the cache connection URL or the log output level.

In the `main()` the HTTP server is started by using the `http.ListenAndServe()` function. Before starting the HTTP server, all routes and their respective handler are declared and added to the API server. 

In the following code snippet, three routes are declared. The endpoints are `/health` , `/slack`, a `/slack/actions`. 

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

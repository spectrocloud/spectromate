# SpectroMate


<p align="center">The home of SpectroMate :robot: </p>

<p align="center">
  <img src="/static/images/mascot.png" alt="drawing" width="250"/>
</p>


## Overview

SpectroMate is an API server with extended functionality designed for Slack integration in the form of a bot. You can use SpectroMate to handle [slash commands](https://api.slack.com/interactivity/slash-commands), and [message actions](https://api.slack.com/reference/interaction-payloads). You can also use SpectroMate to handle non-slack-related events by creating API endpoints for other purposes. 

SpectroMate is designed for deployment in [Palette](https://console.spectrocloud.com) using the Palette Dev Engine (PDE). Using simplifies the management and deployment of SpectroMate.

The following is an architecture overview of SpectroMate. 

![An architecture diagram with all the components that support SpectroMate](./static/images/infrastructure-architecture.png)


## API Endpoints

The following endpoints are available.

| Description                                               | Endpoint           |
| ----------------------------------------------------------|-------------------|
| Used for health checks by external resources.             | `/health`          |
| A slack endpoint that can be used to handle slash commands.| `/slack`           |
| A slack endpoint for handling slack message actions.      | `/slack/actions`   |


## Slack Commands

The following endpoints are available.

| Description                                               | Endpoint           |
| ----------------------------------------------------------|-------------------|
| Displays information to the user for how to use SpectroMate. Invalid commands return the help response.             | `/help`          |
| Used to query the Mendable and ask documentation questions to a trained model.| `/ask`           |
| Same as the `/ask` but responses are only visible to the user versus the entire channel.      | `/pask`   |
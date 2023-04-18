# SpectroMate


<p align="center">The home of SpectroMate :robot: </p>

<p align="center">
  <img src="/static/images/mascot.png" alt="drawing" width="250"/>
</p>


## Overview 👩‍🚀 🧑‍🚀 🧑🏿‍🚀

SpectroMate is an API server with extended functionality designed for Slack integration in the form of a bot. You can use SpectroMate to handle [slash commands](https://api.slack.com/interactivity/slash-commands), and [message actions](https://api.slack.com/reference/interaction-payloads). You can also use SpectroMate to handle non-slack-related events by creating API endpoints for other purposes. 

SpectroMate is designed for deployment in [Palette](https://console.spectrocloud.com) using the Palette Dev Engine (PDE). Using simplifies the management and deployment of SpectroMate.

---

## Getting Started 🚀

To get started with Spectromate follow the steps outlined in the [Getting Started](./docs/getting-started.md) guide.

---

## API Endpoints 🕹️

The following endpoints are available.

| Description                                               | Endpoint           |
| ----------------------------------------------------------|-------------------|
| Used for health checks by external resources.             | `/health`          |
| A slack endpoint that can be used to handle slash commands.| `/slack`           |
| A slack endpoint for handling slack message actions.      | `/slack/actions`   |


## Slack Commands 🛠️

The following Slack commands are available.

| Description                                               | Command           |
| ----------------------------------------------------------|-------------------|
| Displays information to the user for how to use SpectroMate. Invalid commands return the help response.             | `/help`          |
| Used to query the Mendable and ask documentation questions to a trained model.| `/ask`           |
| Same as the `/ask` but responses are only visible to the user versus the entire channel.      | `/pask`   |


## Slack Actions 🪡

Spectromate supports the following actions.

| Description                                               | Action           |
| ----------------------------------------------------------|-------------------|
| Handles the possitive feedback button and submits the feedback to Mendable.  | `ask_model_positive_feedback` |
| Handles the negavtive feedback button and submits the feedback to Mendable.| `ask_model_negative_feedback` |


## Architecture

The following is an architectural overview of SpectroMate. 

![An architecture diagram with all the components that support SpectroMate](./static/images/infrastructure-architecture.png)


## Supported Features and Limitations

|Action| Supported | Notes |
|---|---|---|
| Slash command| ✅ | Supported through the `/slack` endpoint.|
| Message buttons | ✅| Supported through the `/slack/actions` endpoint.|
| Mentions | ❌ | Currently unavailable. |
| Threads | ❌ | Currently unavailable. |
| Health checks | ✅ | Supported through the `/health` endpoint.|
| Verify Slack signature| ✅ | Verification of Slack signature is applied to all Slack endpoints.|
| Metrics | ❌ | Currently unavailable. |


!> There is a limitation with `pask` messages and submitting feedback. The original message is replaced with a feedback acknowledgment. This behavior stems from the Slack API not including the original message when handling action events from an ephemeral message.

# Contribution

We welcome all types of contributions. Please take a moment and review our [contribution guidelines](./docs/contributions.md).

# Open Source Licenses

Review the [Open Source Acknowledgment]((./docs/open-source.md)) reference resource for a complete list of open-source licenses used in this project.
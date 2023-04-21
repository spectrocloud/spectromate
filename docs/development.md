# Development Environment

To get started, use the following steps.

## Software Requirements


If you want to work on Spectromate, you will need to install the following software and actions completed.

:warning: This project uses Go Modules making it safe to work with it outside of your existing GOPATH


* [Docker](https://docs.docker.com/get-docker/) v20.10.0 or greater.

* [Go](https://golang.org/doc/install) 1.20 or greater.

* [Ngrok](https://ngrok.com/download) v3.2.2 or greater.

* [A Slack App](https://api.slack.com/apps) available that has access to your target workspace.

## Setup
. The instructions that follow assume a directory in your home directory outside of the standard GOPATH.

1. Clone the repository into by using the following command.

    ```shell
    $ mkdir -p $HOME/projects/spectromate/; cd $HOME/projects/spectromate/
    $ git clone git@github.com:spectrocloud/spectromate.git
    ```

2. Issue the command `make init`. This will download all the required project dependencies.
To set up the local development support dependencies, such as Redis and Ngork, use the following command.

    ```shell
    make start
    ```

3. Copy the URL Ngrok exposed in the output. 

    ```shell
    $ make start
    Starting the application and ngrok...
    [+] Running 3/3
    ✔ Network docs-slack-bot_default        Created                                                                                                                                                       0.0s
    ✔ Container docs-slack-bot-redis-1      Started                                                                                                                                                       0.3s
    ✔ Container docs-slack-bot-redis-cli-1  Started                                                                                                                                                       0.5s
    Ngrok URL: https://76ac-184-179-188-29.ngrok-free.app
    ```

4. Create a `env` file at the root of the project.

    ```shell
    touch .env
    ```

5. Populate the **.env** file with your respective credentials.

    ```shell
    export SLACK_SIGNING_SECRET=....
    export MENDABLE_API_KEY=....
    ```

6.  Start the local server.

```shell
make server
```

7. Navigate to your [Slack app ](https://api.slack.com/apps) and open up your app's details page. 


8. You will access both elements highlighted in the image. Start with the first element, by clicking on the left **Main Menu** and selecting **Slash Command**.

    ![A view of the slack app details view](../static/images/slack_app_main_view.png)

9. Create a new slash command titled `docs` and provide your Ngrok URL to the command. Append `/api/v1/slack` to the end of URL.

    ![View of the Slack slash command page](../static/images/slack_slash_command.png)

10. Save your changes and navigate back to the main details page.


12. From the left **Main Menu**, select **Interactivity & Shortcuts**.

13. Toggle the **Interactvity** button to **On**. 

14. Use your Ngrok URL and past the URL into the text input box. Append `/api/v1/slack/actions` to the end of URL and save your changes.

    ![View of the interactivity options view](../static/images/slack_interactivity_view.png)

15. Navigate to a Slack channel and issue the Slack command `/docs help`.


You should receive an output similar to the following.

![View of the Slack command in use](../static/images/slack_command_issued.png)

You now have a functional test application connected to your local workstation. 

16. Stop the local server by pressing the `Ctrl` + `C`.


Make your changes and start the server with the `make server` command to test the new behavior. 

```shell
make server
Starting the server...
go run main.go
2023/04/20 19:44:15 maxprocs: Leaving GOMAXPROCS=8: CPU quota undefined
{"level":"info","time":1682045055,"message":"Server is configured for port 3000 and listing on 0.0.0.0:3000"}
{"level":"info","time":1682045055,"message":"Redis is configured for localhost:6379"}
{"level":"info","time":1682045055,"message":"Trace level set to: INFO"}
{"level":"info","time":1682045055,"message":"Starting server..."}


```


## Testing

Add test cases to new functions and new commands. Invoke the Go tests from the root namespace. The pipeline will invoke the Go tests as well.
```shell
go test -v ./...
```

### Terraform

A terraform folder is included in this project. The Terraform code deploys  55 Lambdas. This is used to test the LambdaListFunctions API logic. To use this Terraform code you must provide your own AWS environment and credentials.
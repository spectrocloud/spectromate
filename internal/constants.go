package internal

const (
	ApiVersion                    string = "v1/"
	ApiPath                       string = "/api/"
	ApiPrefixV1                   string = ApiPath + ApiVersion
	SlackPostMessageURL           string = "https://slack.com/api/chat.postMessage"
	SlackDefaultUserErrorMessage  string = "An error occured with the help command. Please reach out to `#docs` for assistance."
	MendandableNewConversationURL string = "https://api.mendable.ai/v0/newConversation"
)

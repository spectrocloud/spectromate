package internal

import "math/rand"

const (
	ApiVersion                    string = "v1/"
	ApiPath                       string = "/api/"
	ApiPrefixV1                   string = ApiPath + ApiVersion
	SlackPostMessageURL           string = "https://slack.com/api/chat.postMessage"
	SlackDefaultUserErrorMessage  string = "An error occured with the help command. Please reach out to `#docs` for assistance."
	MendandableNewConversationURL string = "https://api.mendable.ai/v0/newConversation"
	MendableChatQueryURL          string = "https://api.mendable.ai/v0/mendableChat"
	DefaultUserErrorMessage       string = `:warning: I'm sorry, I'm having technical issues. Notify the the docs team @ #docs and please try again later.`
)

// GetRandomWaitMessage returns a random wait message.
func GetRandomWaitMessage() string {
	waitMessages := []string{
		":hourglass_flowing_sand: Hang tight while I review the docs...",
		":hourglass_flowing_sand: Just a moment while I explore the documentation rabbit hole...",
		":hourglass_flowing_sand: Hold tight while I decode this technical jargon for you!",
		":hourglass_flowing_sand: Please be patient as I consult the ancient scrolls of documentation...",
		":hourglass_flowing_sand: Just a little bit longer, we're digging through the treasure trove of documentation to find what you need!",
		":hourglass_flowing_sand: Hold on tight, we're flipping through pages to bring you the answer!",
		":hourglass_flowing_sand: Just a sec, we're diving deep into the ocean of documentation to find the hidden gems.",
		":hourglass_flowing_sand: Please hold while we decipher the clues in the documentation mystery!",
		":hourglass_flowing_sand: Hold tight while we navigate through the labyrinth of documentation and find the solution.",
		":hourglass_flowing_sand: Just a few more seconds... we're almost done unraveling the mysteries of the documentation.",
		":hourglass_flowing_sand: We're almost there, just a little bit longer until we decode the secrets of the documentation universe!",
	}

	return waitMessages[rand.Intn(len(waitMessages))]
}

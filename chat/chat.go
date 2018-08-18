package chat

import "net/http"

// Chat inteface should be implemented for all messengers(facebook, slack, telegram, whatever)
type Chat interface {
	SendUpdate(string)
	PullRequestCreated(http.ResponseWriter, *http.Request)
	PullRequestCommented(http.ResponseWriter, *http.Request)
	PullRequestApproved(http.ResponseWriter, *http.Request)
	PullRequestMerged(http.ResponseWriter, *http.Request)
}

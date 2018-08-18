package chat

// Chat inteface should be implemented for all messengers(facebook, slack, telegram, whatever)
type Chat interface {
	SendUpdate(string)
}

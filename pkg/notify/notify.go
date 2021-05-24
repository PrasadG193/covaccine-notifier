package notify

type Notifier interface {
	SendMessage(string) error
}

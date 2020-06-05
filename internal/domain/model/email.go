package model

type Email struct {
	Subject string
	Size    int64
	ID      string
	Date    string
	Snippet string
	Content string
}

type MailParsingStatus string

const (
	MailParsingStatusDone    = "DONE"
	MailParsingStatusPending = "PENDING"
	MailParsingStatusError   = "ERROR"
)

type EmailParsingJob struct {
	ID       string
	Status   MailParsingStatus
	Warnings []string
	Detail   string
	Subject  string
	Trip     Trip
}

package domain

import "amadeus-trip-parser/internal/domain/model"

type EmailProvider interface {
	GetEmails(filter string) []*model.Email
}

type EmailParser interface {
	CreateJob(mail *model.Email) (*model.EmailParsingJob, error)
	GetJobStatus(job model.EmailParsingJob) (*model.EmailParsingJob, error)
	GetJobResult(job model.EmailParsingJob) (*model.EmailParsingJob, error)
}

package usecase

import (
	"amadeus-trip-parser/internal/domain"
	"amadeus-trip-parser/internal/domain/model"
	"github.com/rs/zerolog/log"
	"time"
)

type emailProcessor struct {
	provider       domain.EmailProvider
	parser         domain.EmailParser
	repo           domain.TripRepository
	parserInterval time.Duration
	mailInterval   time.Duration
	emails         chan *model.Email
	toRefresh      chan *model.EmailParsingJob
	resultReady    chan *model.EmailParsingJob
	done           chan bool
}

func NewEmailProcessor(provider domain.EmailProvider, parser domain.EmailParser,
	repo domain.TripRepository) domain.EmailProcessor {
	return &emailProcessor{
		provider,
		parser,
		repo,
		15 * time.Second,
		10 * time.Minute,
		make(chan *model.Email),
		make(chan *model.EmailParsingJob),
		make(chan *model.EmailParsingJob),
		make(chan bool),
	}
}

func (e *emailProcessor) Process() {
	go e.fetchEmail()
	go e.createJob()
	go e.checkJobStatus()
	go e.getResult()
}

func (e *emailProcessor) Stop() {
	e.done <- true
}

func (e *emailProcessor) fetchEmail() {
	for {
		//TODO allow mail filter configuration
		emails := e.provider.GetEmails("is:unread")
		select {
		case <-e.done:
			close(e.done)
			return
		default:
			for _, em := range emails {
				e.emails <- em
			}
		}
		time.Sleep(e.mailInterval)
	}
}

func (e *emailProcessor) createJob() {
	for {
		select {
		case email := <-e.emails:
			job, err := e.parser.CreateJob(email)
			if err != nil {
				log.Debug().Msgf("error when creating job for %v: %v", email, err)
				continue
			}
			log.Debug().Msgf("job created %v", job)
			e.toRefresh <- job
		case <-e.done:
			close(e.toRefresh)
			return
		}
		time.Sleep(e.parserInterval)
	}
}

func (e *emailProcessor) checkJobStatus() {
	for {
		select {
		case job := <-e.toRefresh:
			refreshedJob, err := e.parser.GetJobStatus(*job)
			if err != nil {
				log.Debug().Msgf("error when refreshing job %s: %v", job, err)
				continue
			}
			switch refreshedJob.Status {
			case model.MailParsingStatusPending:
				log.Debug().Msgf("job %s is still pending", refreshedJob.ID)
				go func() {
					time.Sleep(e.parserInterval)
					select {
					case <-e.done:
						return
					default:
						e.toRefresh <- refreshedJob
					}
				}()
			case model.MailParsingStatusError:
				log.Debug().Msgf("job %s is in error: %s", refreshedJob.ID, refreshedJob.Detail)
			case model.MailParsingStatusDone:
				log.Debug().Msgf("job %s is done", refreshedJob.ID)
				e.resultReady <- refreshedJob
			default:
				log.Debug().Msgf("job %s has unknown parsing status %s", refreshedJob.ID, refreshedJob.Status)
			}
		case <-e.done:
			close(e.resultReady)
			return
		}
		time.Sleep(e.parserInterval)
	}
}

func (e *emailProcessor) getResult() {
	for {
		select {
		case job := <-e.resultReady:
			jobWithResult, err := e.parser.GetJobResult(*job)
			if err != nil {
				log.Debug().Msgf("failed to retrieve result for job %s : %v", job.ID, err)
				continue
			}
			e.storeTrip(jobWithResult.Trip)
		case <-e.done:
			return
		}
		time.Sleep(e.parserInterval)
	}
}

func (e *emailProcessor) storeTrip(trip model.Trip) {
	if err := e.repo.Create(&trip); err != nil {
		log.Debug().Msgf("failed to store trip %v: %v", trip, err)
	} else {
		log.Debug().Msgf("trip %s (ref: %s) written in repository", trip.ID, trip.Reference)
	}
}

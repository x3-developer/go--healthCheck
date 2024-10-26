package notifier

import "HealthCheck/config"

type NotificationType string

const (
	Alert NotificationType = "alert"
	Calm                   = "calm"
)

type Notifier interface {
	Notify(nt NotificationType, site config.Site) error
}

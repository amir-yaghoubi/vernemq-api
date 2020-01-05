package auth

import (
	"errors"

	mqttpattern "github.com/amir-yaghoubi/mqtt-pattern"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasttemplate"
)

// InvalidSubscription qos value for invalid subscription
const InvalidSubscription = 128

// New returns new Service instance
func New(repo Repository, logger *logrus.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Service auth service
type Service struct {
	repo   Repository
	logger *logrus.Logger
}

// Subscription type
type Subscription struct {
	Topic string `json:"topic"`
	Qos   uint8  `json:"qos"`
}

// Modifiers type
type Modifiers struct {
	ClientID   string `json:"client_id"`
	Mountpoint string `json:"mountpoint"`
}

func (s *Service) authorizePublish(clientID string, username string, topic string, qos uint8, retain bool, publishACL []PublishACL) bool {
	if publishACL == nil {
		return true
	}

	for i := range publishACL {
		// TODO we could improve performance by caching templates
		tmp, err := fasttemplate.NewTemplate(publishACL[i].Pattern, "{{", "}}")
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"clientID": clientID,
				"username": username,
				"topic":    topic,
				"qos":      qos,
				"retain":   retain,
				"pattern":  publishACL[i].Pattern,
				"err":      err.Error(),
			}).Warn("cannot parse pattern")
			return false
		}

		pattern := tmp.ExecuteString(map[string]interface{}{
			"client_id": clientID,
			"username":  username,
		})

		isMatch := mqttpattern.Matches(pattern, topic)
		if isMatch {
			if publishACL[i].AllowedRetain != nil && *publishACL[i].AllowedRetain != retain {
				continue
			}

			if publishACL[i].MaxQos != nil && *publishACL[i].MaxQos < qos {
				continue
			}
			return true
		}
	}

	return false
}

func (s *Service) authorizeSubscribe(clientID string, username string, subscriptions []Subscription, subACL []SubACL) []Subscription {
	if subACL == nil {
		return subscriptions
	}

	var isProcessed bool
	results := make([]Subscription, 0, len(subscriptions))
	for ss := range subscriptions {
		isProcessed = false

		for i := range subACL {
			// TODO we could improve performance by caching templates
			tmp, err := fasttemplate.NewTemplate(subACL[i].Pattern, "{{", "}}")
			if err != nil {
				s.logger.WithFields(logrus.Fields{
					"clientID": clientID,
					"username": username,
					"pattern":  subACL[i].Pattern,
					"err":      err.Error(),
				}).Warn("cannot parse pattern")
				continue
			}

			pattern := tmp.ExecuteString(map[string]interface{}{
				"client_id": clientID,
				"username":  username,
			})

			isMatched := mqttpattern.Matches(pattern, subscriptions[ss].Topic)
			if isMatched {
				isProcessed = true
				if subACL[i].MaxQos != nil && *subACL[i].MaxQos < subscriptions[ss].Qos {
					results = append(results, Subscription{Topic: subscriptions[ss].Topic, Qos: InvalidSubscription})
				} else {
					results = append(results, subscriptions[ss])
				}
			}
		}

		if !isProcessed {
			results = append(results, Subscription{Topic: subscriptions[ss].Topic, Qos: InvalidSubscription})
		}
	}

	return results
}

// Authenticate with clientID, username and password
func (s *Service) Authenticate(clientID string, username string, password string) (*Modifiers, error) {
	user, err := s.repo.Get(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUnAuthorizeAccess
		}
		return nil, err
	}

	// TODO hash password
	if user.Password != password {
		return nil, ErrUnAuthorizeAccess
	}

	return &Modifiers{ClientID: clientID, Mountpoint: user.Mountpoint}, nil
}

// AuthorizePublish with user ACL
func (s *Service) AuthorizePublish(clientID string, username string, topic string, qos uint8, retain bool) (bool, error) {
	user, err := s.repo.Get(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	return s.authorizePublish(clientID, username, topic, qos, retain, user.PublishACL), nil
}

// AuthorizeSubscribe with user ACL
func (s *Service) AuthorizeSubscribe(clientID string, username string, subscriptions []Subscription) ([]Subscription, error) {
	user, err := s.repo.Get(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			invalidTopics := make([]Subscription, 0, len(subscriptions))
			for i := range subscriptions {
				invalidTopics = append(invalidTopics, Subscription{Topic: subscriptions[i].Topic, Qos: InvalidSubscription})
			}
			return invalidTopics, nil
		}
		return nil, err
	}

	return s.authorizeSubscribe(clientID, username, subscriptions, user.SubACL), nil
}

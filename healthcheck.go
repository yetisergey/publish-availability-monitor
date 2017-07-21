package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
	"github.com/Financial-Times/service-status-go/gtg"
	log "github.com/Sirupsen/logrus"
)

const requestTimeout = 4500

// Healthcheck offers methods to measure application health.
type Healthcheck struct {
	client          *http.Client
	config          *AppConfig
	consumer        consumer.MessageConsumer
	metricContainer *publishHistory
}

func newHealthcheck(config *AppConfig, metricContainer *publishHistory) *Healthcheck {
	httpClient := &http.Client{Timeout: requestTimeout * time.Millisecond}
	c := consumer.NewConsumer(config.QueueConf, func(m consumer.Message) {}, httpClient)
	return &Healthcheck{
		client:          httpClient,
		config:          config,
		consumer:        c,
		metricContainer: metricContainer,
	}
}

type readEnvironmentHealthcheck struct {
	env    Environment
	client *http.Client
}

const pam_run_book_url = "https://dewey.ft.com/publish-availability-monitor.html"

var readCheckEndpoints = map[string]func(string) (string, error){
	"S3": buildAwsHealthcheckUrl,
	// only exceptions need to be listed here - everything else will default to standard FT healthcheck URLs
}

var noReadEnvironments = fthealth.Check{
	BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
	Name:             "ReadEnvironments",
	PanicGuide:       pam_run_book_url,
	Severity:         1,
	TechnicalSummary: "There are no read environments to monitor. This could be because none have been configured, or that etcd is not reachable/healthy",
	Checker: func() (string, error) {
		return "", errors.New("There are no read environments to monitor.")
	},
}

func (h *Healthcheck) checkHealth(writer http.ResponseWriter, req *http.Request) {
	checks := make([]fthealth.Check, 4)
	checks[0] = h.messageQueueProxyReachable()
	checks[1] = h.reflectPublishFailures()
	checks[2] = h.validationServicesReachable()
	checks[3] = isConsumingFromPushFeeds()

	readEnvironmentChecks := h.readEnvironmentsReachable()
	if len(readEnvironmentChecks) == 0 {
		checks = append(checks, noReadEnvironments)
	} else {
		for _, hc := range readEnvironmentChecks {
			checks = append(checks, hc)
		}
	}

	fthealth.HandlerParallel(
		"Dependent services healthcheck", "Checks if all the dependent services are reachable and healthy.",
		checks...,
	)(writer, req)
}

func (h *Healthcheck) GTG() gtg.Status {
	consumerCheck := func() gtg.Status {
		return gtgCheck(h.consumer.ConnectivityCheck)
	}

	validationServiceCheck := func() gtg.Status {
		return gtgCheck(h.checkValidationServicesReachable)
	}

	return gtg.FailFastParallelCheck([]gtg.StatusChecker{
		consumerCheck,
		validationServiceCheck,
	})()
}

func gtgCheck(handler func() (string, error)) gtg.Status {
	if _, err := handler(); err != nil {
		return gtg.Status{GoodToGo: false, Message: err.Error()}
	}
	return gtg.Status{GoodToGo: true}
}

func isConsumingFromPushFeeds() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
		Name:             "IsConsumingFromNotificationsPushFeeds",
		PanicGuide:       pam_run_book_url,
		Severity:         1,
		TechnicalSummary: "The connections to the configured notifications-push feeds are operating correctly.",
		Checker: func() (string, error) {
			var failing []string
			result := true
			for _, val := range subscribedFeeds {
				for _, feed := range val {
					push, ok := feed.(*feeds.NotificationsPushFeed)
					if ok && !push.IsConnected() {
						log.Warnf("Feed \"%s\" with URL \"%s\" is not connected!", feed.FeedName(), feed.FeedURL())
						failing = append(failing, feed.FeedURL())
						result = false
					}
				}
			}

			if !result {
				return "Disconnection detected.", errors.New("At least one of our Notifcations Push feeds in the delivery cluster is disconnected! Please review the logs, and check delivery healthchecks. We will attempt reconnection indefinitely, but there could be an issue with the delivery cluster's notifications-push services. Failing connections: " + strings.Join(failing, ","))
			}
			return "", nil
		},
	}
}

func (h *Healthcheck) messageQueueProxyReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
		Name:             "MessageQueueProxyReachable",
		PanicGuide:       pam_run_book_url,
		Severity:         1,
		TechnicalSummary: "Message queue proxy is not reachable/healthy",
		Checker:          h.consumer.ConnectivityCheck,
	}
}

func (h *Healthcheck) reflectPublishFailures() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "At least two of the last 10 publishes failed. This will reflect in the SLA measurement.",
		Name:             "ReflectPublishFailures",
		PanicGuide:       pam_run_book_url,
		Severity:         1,
		TechnicalSummary: "Publishes did not meet the SLA measurments",
		Checker:          h.checkForPublishFailures,
	}

}

func (h *Healthcheck) checkForPublishFailures() (string, error) {
	h.metricContainer.RLock()
	failures := make(map[string]struct{})
	var emptyStruct struct{}
	for i := 0; i < len(h.metricContainer.publishMetrics); i++ {

		if !h.metricContainer.publishMetrics[i].publishOK {
			failures[h.metricContainer.publishMetrics[i].UUID] = emptyStruct
		}
	}
	h.metricContainer.RUnlock()

	failureThreshold := 2 //default
	if h.config.HealthConf.FailureThreshold != 0 {
		failureThreshold = h.config.HealthConf.FailureThreshold
	}

	if len(failures) >= failureThreshold {
		return "", fmt.Errorf("%d publish failures happened during the last 10 publishes", failures)
	}
	return "", nil
}

func (h *Healthcheck) validationServicesReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics might not be correct. False positive failures might be recorded. This will impact the SLA measurement.",
		Name:             "validationServicesReachable",
		PanicGuide:       pam_run_book_url,
		Severity:         1,
		TechnicalSummary: "Validation services are not reachable/healthy",
		Checker:          h.checkValidationServicesReachable,
	}
}

func (h *Healthcheck) checkValidationServicesReachable() (string, error) {
	endpoints := h.config.ValidationEndpoints
	var wg sync.WaitGroup
	hcErrs := make(chan error, len(endpoints))
	for _, url := range endpoints {
		wg.Add(1)
		healthcheckURL, err := inferHealthCheckUrl(url)
		if err != nil {
			log.Errorf("Validation Service URL: [%s]. Err: [%v]", url, err.Error())
			continue
		}

		go checkServiceReachable(healthcheckURL, h.client, hcErrs, &wg)
	}

	wg.Wait()
	close(hcErrs)
	for err := range hcErrs {
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func checkServiceReachable(healthcheckURL string, client *http.Client, hcRes chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Checking: %s", healthcheckURL)

	resp, err := client.Get(healthcheckURL)
	if err != nil {
		hcRes <- fmt.Errorf("Healthcheck URL: [%s]. Error: [%v]", healthcheckURL, err)
		return
	}
	defer cleanupResp(resp)
	if resp.StatusCode != 200 {
		hcRes <- fmt.Errorf("Unhealthy statusCode received: [%d] for URL [%s]", resp.StatusCode, healthcheckURL)
		return
	}
	hcRes <- nil
}

func (h *Healthcheck) readEnvironmentsReachable() []fthealth.Check {
	hc := make([]fthealth.Check, len(environments))

	i := 0
	for _, env := range environments {
		hc[i] = fthealth.Check{
			BusinessImpact:   "Publish metrics might not be correct. False positive failures might be recorded. This will impact the SLA measurement.",
			Name:             env.Name + " readEndpointsReachable",
			PanicGuide:       pam_run_book_url,
			Severity:         1,
			TechnicalSummary: "Read services are not reachable/healthy",
			Checker:          (&readEnvironmentHealthcheck{env, h.client}).checkReadEnvironmentReachable,
		}
		i++
	}
	return hc
}

func (h *readEnvironmentHealthcheck) checkReadEnvironmentReachable() (string, error) {
	var wg sync.WaitGroup
	hcErrs := make(chan error, len(appConfig.MetricConf))

	for _, metric := range appConfig.MetricConf {
		var endpointURL *url.URL
		var err error

		if absoluteUrlRegex.MatchString(metric.Endpoint) {
			endpointURL, err = url.Parse(metric.Endpoint)
		} else {
			if metric.Alias == "S3" {
				endpointURL, err = url.Parse(h.env.S3Url + metric.Endpoint)
			} else {
				endpointURL, err = url.Parse(h.env.ReadUrl + metric.Endpoint)
			}
		}

		if err != nil {
			log.Errorf("Cannot parse url [%v], Err: [%v]", metric.Endpoint, err.Error())
			continue
		}

		var healthcheckURL string
		if fn, ok := readCheckEndpoints[metric.Alias]; ok {
			healthcheckURL, err = fn(endpointURL.String())
		} else {
			healthcheckURL, err = buildFtHealthcheckUrl(*endpointURL, metric.Health)
		}

		if err != nil {
			log.Errorf("Service URL: [%s]. Err: [%v]", endpointURL.String(), err.Error())
			continue
		}

		wg.Add(1)
		go checkServiceReachable(healthcheckURL, h.client, hcErrs, &wg)
	}

	wg.Wait()
	close(hcErrs)
	for err := range hcErrs {
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func inferHealthCheckUrl(serviceUrl string) (string, error) {
	parsedURL, err := url.Parse(serviceUrl)
	if err != nil {
		return "", err
	}

	var newPath string
	if strings.HasPrefix(parsedURL.Path, "/__") {
		newPath = strings.SplitN(parsedURL.Path[1:], "/", 2)[0] + "/__health"
	} else {
		newPath = "/__health"
	}

	parsedURL.Path = newPath
	return parsedURL.String(), nil
}

func buildFtHealthcheckUrl(endpoint url.URL, health string) (string, error) {
	endpoint.Path = health
	endpoint.RawQuery = "" // strip query params
	return endpoint.String(), nil
}

func buildAwsHealthcheckUrl(serviceUrl string) (string, error) {
	return serviceUrl + "healthCheckDummyFile", nil
}

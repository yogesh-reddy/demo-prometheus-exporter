package collector

import "github.com/prometheus/client_golang/prometheus"

// ApplicationMetrics implements of  prometheus.collector interface
type ApplicationMetrics struct {
	scrapeCount       prometheus.Counter
	failureScrapes    prometheus.Counter
	SuccessfulScrapes prometheus.Counter
	dataMetrics       []*MetricDefinition
}

// MetricDefinition defines metric description type and value function
type MetricDefinition struct {
	Type  prometheus.ValueType
	Desc  *prometheus.Desc
	Value func(response *ExampleData) float64
}

//ExampleData response from app that exposes metric data that needs to be instrumented
type ExampleData struct {
	AvailableReplicas int
	DesiredReplicas   int
	Requests          int
	Errors            int
	MemoryUsed        float64
	MemoryAvailable   float64
}

func NewApplicationMetrics() *ApplicationMetrics {
	subsystem := "app"

	return &ApplicationMetrics{
		scrapeCount: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "scrapes_total"),
			Help: "Current total scrapes.",
		}),
		failureScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "scrapes_failure_count"),
			Help: "Number of errors while scraping.",
		}),
		SuccessfulScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subsystem, "scrapes_success_count"),
			Help: "Total success scrapes.",
		}),
		dataMetrics: []*MetricDefinition{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "available_replicas"),
					"The number of available replicas.",
					nil, nil,
				),
				Value: func(resp *ExampleData) float64 {
					return float64(resp.AvailableReplicas)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "desired_replicas"),
					"The number of desired replicas.",
					nil, nil,
				),
				Value: func(resp *ExampleData) float64 {
					return float64(resp.DesiredReplicas)
				},
			},
			{
				Type: prometheus.CounterValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "errors"),
					"The number of errors in app.",
					nil, nil,
				),
				Value: func(resp *ExampleData) float64 {
					return float64(resp.Errors)
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(namespace, subsystem, "memory_used"),
					"Memory used by the App",
					nil, nil,
				),
				Value: func(resp *ExampleData) float64 {
					return resp.MemoryUsed
				},
			},
		},
	}

}


// Describe set Prometheus metrics descriptions.
func (c *ApplicationMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.dataMetrics {
		ch <- metric.Desc
	}
	ch <- c.failureScrapes.Desc()
	ch <- c.scrapeCount.Desc()
	ch <- c.SuccessfulScrapes.Desc()
}

// Collect collects Application metrics.
func (c *ApplicationMetrics) Collect(ch chan<- prometheus.Metric) {
	var err error
	c.scrapeCount.Inc()

	appResponse, err := fetchAppData()
	if err != nil {
		c.failureScrapes.Inc()
		return
	}
	c.SuccessfulScrapes.Inc()

	for _, metric := range c.dataMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(appResponse),
		)
	}
	defer func() {
		ch <- c.scrapeCount
		ch <- c.SuccessfulScrapes
		ch <- c.failureScrapes
	}()
}

//would be actual api call
func fetchAppData() (*ExampleData ,error) {
	//here goes the code for contacting the actuall app and fetching the data
	return &ExampleData{
		AvailableReplicas: 1,
		DesiredReplicas: 2,
		Requests: 30,
		Errors: 5,
		MemoryUsed: 512,
		MemoryAvailable: 2048,
	} ,nil
}
package collector

import "github.com/prometheus/client_golang/prometheus"

type simpleCollector struct {
	metric1 *prometheus.Desc
	metric2 *prometheus.Desc
}

const (
	namespace = "demo"
	subsystem1= "collector1"
)

func NewSimpleCollector() *simpleCollector {
	return &simpleCollector{
		metric1: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem1, "metric_one"),
			"help message for metric 1", []string{"label1","label2"}, nil),
		metric2: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem1, "metric_two"),
			"help message for metric 2", []string{}, nil),
	}
}

func (collector *simpleCollector) Describe(ch chan<- *prometheus.Desc) {
	//write as metric descriptors into the channel
	ch <- collector.metric1
	ch <- collector.metric2
}

func (collector *simpleCollector) Collect(ch chan<- prometheus.Metric) {
	//Write value for each metric in the prometheus metric channel.
	ch <- prometheus.MustNewConstMetric(collector.metric1, prometheus.GaugeValue, 1,"lvalue1","lvalue2" )
	ch <- prometheus.MustNewConstMetric(collector.metric2, prometheus.CounterValue, 2, )
}
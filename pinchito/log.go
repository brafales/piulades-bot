package pinchito

type Log struct{}

func (*Log) Body() string {
	return "Testing"
}

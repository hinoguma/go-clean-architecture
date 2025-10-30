package crosscuttinglayer

import "app/threelayeredarchitecture/crosscuttinglayer/adapter"

var globalTempIDGenerator adapter.StrIDGenerator

func SetGlobalTempIDGenerator(gen adapter.StrIDGenerator) {
	globalTempIDGenerator = gen
}

func IssueRandomStrID() string {
	return globalTempIDGenerator.Issue()
}

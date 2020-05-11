package main

// jobDetailMap, key is the repo name, value is the list of job types, like continuous, nightly, etc., as well as custome names
type JobDetailMap map[string][]string

// testGridMetaData saves the meta data needed to generate the final config file.
// key is the main project version, value is another map containing job details
type TestGridMetaData struct {
	md map[string]JobDetailMap
	// projNames save the proj names in a list when parsing the simple config file, for the purpose of maintaining the output sequence
	projNames []string
	// repoNames save the repo names in a list when parsing the simple config file, for the purpose of maintaining the output sequence
	repoNames []string
}

func (j JobDetailMap) Add(repo, jt string) {
	j.EnsureExists(repo)
	j[repo] = append(j[repo], jt)
}

func NewJobDetailMap() JobDetailMap {
	return make(JobDetailMap)
}

// EnsureExists returns true if already existed or false if newly-created
func (j JobDetailMap) EnsureExists(repo string) bool {
	if _, exists := j[repo]; exists == false {
		j[repo] = make([]string, 0)
		return false
	}
	return true
}

func NewTestGridMetaData() TestGridMetaData {
	return TestGridMetaData{
		md:        make(map[string]JobDetailMap),
		projNames: make([]string, 0),
		repoNames: make([]string, 0),
	}
}

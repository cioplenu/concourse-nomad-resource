package resource

type Source struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type JobVersion struct {
	Version    int
	Stable     bool
	SubmitTime int
}

type History []JobVersion

func (h History) Len() int           { return len(h) }
func (h History) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h History) Less(i, j int) bool { return h[i].Version < h[j].Version }

type Version struct {
	Version int `json:"Version,string"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

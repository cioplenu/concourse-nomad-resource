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

type Version struct {
	Version int `json:"Version,string"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

package github

type Schema struct {
	Data Data `json:"data"`
}
type ReleaseNode struct {
	TagName      string `json:"tagName"`
	IsPrerelease bool   `json:"isPrerelease"`
}
type Releases struct {
	Nodes []ReleaseNode `json:"nodes"`
}
type RepoNode struct {
	Name     string   `json:"name"`
	Releases Releases `json:"releases"`
}
type Repositories struct {
	Nodes []RepoNode `json:"nodes"`
}
type Team struct {
	Description  string       `json:"description"`
	Repositories Repositories `json:"repositories"`
}
type Organization struct {
	Team Team `json:"team"`
}
type Data struct {
	Organization Organization `json:"organization"`
}

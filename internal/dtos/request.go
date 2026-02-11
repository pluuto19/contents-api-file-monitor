package dtos

type ReadmeResponseDTO struct {
	Type        string  `json:"type"`
	Encoding    string  `json:"encoding"`
	Size        int     `json:"size"`
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Content     string  `json:"content"`
	Sha         string  `json:"sha"`
	Url         string  `json:"url"`
	GitUrl      *string `json:"git_url,omitempty"`
	HtmlUrl     *string `json:"html_url,omitempty"`
	DownloadUrl *string `json:"download_url,omitempty"`
	Links       struct {
		Git  *string `json:"git"`
		Html *string `json:"html"`
		Self string  `json:"self"`
	} `json:"_links"`
	Target          *string `json:"target,omitempty"`
	SubmoduleGitUrl *string `json:"submodule_git_url,omitempty"`
}
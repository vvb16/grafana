package search

type panelInfo struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	Type            string   `json:"type"`                      // PluginID
	Datasource      []string `json:"datasource,omitempty"`      // UIDs
	DatasourceType  []string `json:"datasourceType,omitempty"`  // PluginIDs
	Transformations []string `json:"transformations,omitempty"` // ids of the transformation steps
}

type dashboardInfo struct {
	UID            string      `json:"uid"`
	Path           string      `json:"path"`
	Title          string      `json:"title"`
	Description    string      `json:"description,omitempty"`
	Tags           []string    `json:"tags"`                     // UIDs
	Datasource     []string    `json:"datasource,omitempty"`     // UIDs
	DatasourceType []string    `json:"datasourceType,omitempty"` // PluginIDs
	PanelTypes     []string    `json:"panelTypes"`               // PluginIDs
	TemplateVars   []string    `json:"templateVars,omitempty"`   // the keys used
	Panels         []panelInfo `json:"panels"`                   // nesed documents
	SchemaVersion  int64       `json:"schemaVersion"`
}

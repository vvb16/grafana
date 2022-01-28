package search

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

func readDashboardJSON(iter *jsoniter.Iterator) *dashboardInfo {
	dash := &dashboardInfo{}

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		switch l1Field {
		case "uid":
			dash.UID = iter.ReadString()

		case "title":
			dash.Title = iter.ReadString()

		case "description":
			dash.Description = iter.ReadString()

		case "schemaVersion":
			dash.SchemaVersion = iter.ReadInt64()

		case "tags":
			for iter.ReadArray() {
				dash.Tags = append(dash.Tags, iter.ReadString())
			}

		case "panels":
			for iter.ReadArray() {
				dash.Panels = append(dash.Panels, readPanelInfo(iter))
			}

		case "templating":
			for sub := iter.ReadObject(); sub != ""; sub = iter.ReadObject() {
				if "list" == sub {
					for iter.ReadArray() {
						for k := iter.ReadObject(); k != ""; k = iter.ReadObject() {
							if k == "name" {
								dash.TemplateVars = append(dash.TemplateVars, iter.ReadString())
							} else {
								iter.Skip()
							}
						}
					}
				} else {
					iter.Skip()
				}
			}

		default:
			v := iter.Read()
			fmt.Printf("[DASHBOARD] support key: %s / %v\n", l1Field, v)
		}
	}

	return dash
}

// will always return strings for now
func readPanelInfo(iter *jsoniter.Iterator) panelInfo {
	panel := panelInfo{}

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		switch l1Field {
		case "id":
			panel.ID = iter.ReadInt64()

		case "type":
			panel.Type = iter.ReadString()

		case "title":
			panel.Title = iter.ReadString()

		case "description":
			panel.Description = iter.ReadString()

		case "datasource":
			v := iter.Read()
			fmt.Printf(">>Panel.datasource = %v\n", v) // string or object!!!

		case "targets":
			for iter.ReadArray() {
				v := iter.Read()
				fmt.Printf("[Panel.TARGET] %v\n", v)
			}

		case "transformations":
			for iter.ReadArray() {
				for sub := iter.ReadObject(); sub != ""; sub = iter.ReadObject() {
					if sub == "id" {
						panel.Transformations = append(panel.Transformations, iter.ReadString())
					} else {
						iter.Skip()
					}
				}
			}

		case "options":
			fallthrough

		case "gridPos":
			fallthrough

		case "fieldConfig":
			iter.Skip()

		// case "error":
		// case "errorType":
		// case "warnings":
		default:
			v := iter.Read()
			fmt.Printf("[PANEL] support key: %s / %v\n", l1Field, v)
		}
	}

	return panel
}

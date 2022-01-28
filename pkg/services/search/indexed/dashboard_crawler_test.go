package search

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestReadDashboard(t *testing.T) {
	devdash := filepath.Clean("../../../../devenv/dev-dashboards")
	name := "panel-graph/graph-ng-soft-limits.json"

	walker := func(fpath string, info os.FileInfo, e error) error {
		f, err := os.Open(fpath)
		require.NoError(t, err)

		iter := jsoniter.Parse(jsoniter.ConfigDefault, f, 1024)
		dash := readDashboardJSON(iter)

		out, err := json.MarshalIndent(dash, "", "  ")
		require.NoError(t, err)

		fmt.Printf("%s", string(out))
		return nil
	}

	// Single file
	err := walker(path.Join(devdash, name), nil, nil)
	require.NoError(t, err)

	// Walk the whole folder
	err = filepath.Walk(devdash, walker)
	require.NoError(t, err)

	// This makes the printf do something
	t.Fail()
}

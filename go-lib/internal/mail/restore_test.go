package mail

import (
	"testing"

	"github.com/ProtonMail/go-proton-api"
	"github.com/stretchr/testify/require"
)

func TestSortLabels(t *testing.T) {
	l := []proton.Label{
		{ID: "1", ParentID: ""},
		{ID: "2", ParentID: ""},
		{ID: "3", ParentID: ""},
		{ID: "4", ParentID: ""},
	}

	result, err := sortLabels(l)
	require.NoError(t, err)
	require.Equal(t, result, l)

	l = []proton.Label{
		{ID: "1", ParentID: ""},
		{ID: "2", ParentID: "1"},
		{ID: "4", ParentID: "2"},
		{ID: "3", ParentID: "1"},
	}
	result, err = sortLabels(l)
	require.NoError(t, err)
	require.Equal(t, result, l)

	l = []proton.Label{
		{ID: "1", ParentID: "3"},
		{ID: "2", ParentID: ""},
		{ID: "3", ParentID: "2"},
		{ID: "4", ParentID: "3"},
	}
	result, err = sortLabels(l)
	require.NoError(t, err)
	require.Equal(t, result, []proton.Label{
		{ID: "2", ParentID: ""},
		{ID: "3", ParentID: "2"},
		{ID: "4", ParentID: "3"},
		{ID: "1", ParentID: "3"},
	})

	l = []proton.Label{
		{ID: "1", ParentID: "2"},
		{ID: "2", ParentID: "3"},
		{ID: "3", ParentID: "4"},
		{ID: "4", ParentID: ""},
	}
	result, err = sortLabels(l)
	require.NoError(t, err)
	require.Equal(t, result, []proton.Label{
		{ID: "4", ParentID: ""},
		{ID: "3", ParentID: "4"},
		{ID: "2", ParentID: "3"},
		{ID: "1", ParentID: "2"},
	})

	l = []proton.Label{
		{ID: "1", ParentID: "3"},
		{ID: "2", ParentID: ""},
		{ID: "3", ParentID: "1"}, // circular reference 3 <-> 1.
		{ID: "4", ParentID: "3"},
	}

	_, err = sortLabels(l)
	require.Error(t, errCircularLabelReference, err)

	l = []proton.Label{
		{ID: "1", ParentID: "2"},
		{ID: "2", ParentID: "3"},
		{ID: "3", ParentID: "4"},
		{ID: "4", ParentID: "2"}, // circular reference 4 <-> 3 <-> 2.
	}

	_, err = sortLabels(l)
	require.Error(t, errCircularLabelReference, err)
}

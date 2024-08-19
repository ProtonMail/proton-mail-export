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

func TestFindFirstAvailableLabelIncrementalName(t *testing.T) {
	remoteLabels := []proton.Label{
		{Name: "Folder (1)"},
		{Name: "FOLDER (2)"},
		{Name: "FOLDER(3)"},
		{Name: "Folder (4)"},
	}
	require.Equal(t, findFirstAvailableLabelIncrementalName("Folder", remoteLabels), "Folder (3)")
	require.Equal(t, findFirstAvailableLabelIncrementalName("Folder", nil), "Folder (1)")
	require.Equal(t, findFirstAvailableLabelIncrementalName("Folders", remoteLabels), "Folders (1)")
	require.Equal(t, findFirstAvailableLabelIncrementalName("Folder ", remoteLabels), "Folder  (1)")
}

func TestMatchLocalLabelWithRemote(t *testing.T) {
	remoteLabels := []proton.Label{
		{ID: "remoteID_F1", Name: "F1", Type: proton.LabelTypeFolder},
		{ID: "remoteID_F1", Name: "F1 (1)", Type: proton.LabelTypeFolder},
		{ID: "remoteID_L1", Name: "L1", Type: proton.LabelTypeLabel},
	}

	// folder with name does not exist
	folder := proton.Label{ID: "localID_F", Name: "F", Type: proton.LabelTypeFolder}
	labelID, newName := matchLocalLabelWithRemote(folder, remoteLabels)
	require.Len(t, labelID, 0)
	require.Equal(t, newName, "F")

	// folder with name exists and is of the right type
	folder = proton.Label{ID: "localID_F1", Name: "F1", Type: proton.LabelTypeFolder}
	labelID, newName = matchLocalLabelWithRemote(folder, remoteLabels)
	require.Equal(t, labelID, "remoteID_F1")
	require.Len(t, newName, 0)

	// folder with name exists but is not of the right type
	folder = proton.Label{ID: "localID_F1", Name: "F1", Type: proton.LabelTypeLabel}
	labelID, newName = matchLocalLabelWithRemote(folder, remoteLabels)
	require.Len(t, labelID, 0)
	require.Equal(t, newName, "F1 (2)")

	// label with name does not exist
	label := proton.Label{ID: "localID_L", Name: "l", Type: proton.LabelTypeLabel}
	labelID, newName = matchLocalLabelWithRemote(label, remoteLabels)
	require.Len(t, labelID, 0)
	require.Equal(t, newName, "l")

	// label with name exists and is of the right type
	label = proton.Label{ID: "localID_L1", Name: "l1", Type: proton.LabelTypeLabel}
	labelID, newName = matchLocalLabelWithRemote(label, remoteLabels)
	require.Equal(t, labelID, "remoteID_L1")
	require.Len(t, newName, 0)

	// folder with name exists but is not of the right type
	label = proton.Label{ID: "localID_L1", Name: "l1", Type: proton.LabelTypeFolder}
	labelID, newName = matchLocalLabelWithRemote(label, remoteLabels)
	require.Len(t, labelID, 0)
	require.Equal(t, newName, "l1 (1)")
}

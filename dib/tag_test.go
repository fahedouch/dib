package dib_test

import (
	"testing"

	"github.com/radiofrance/dib/dag"
	"github.com/radiofrance/dib/dib"
	"github.com/radiofrance/dib/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Retag_DoesNotRetagIfNoRetagNeeded(t *testing.T) {
	t.Parallel()

	DAG := &dag.DAG{}
	DAG.AddNode(dag.NewNode(&dag.Image{
		Name:      "registry.example.org/image",
		ShortName: "image",
		RetagDone: false,
	}))

	tagger := &mock.Tagger{}
	err := dib.Retag(DAG, tagger, false)

	assert.NoError(t, err)
	assert.Empty(t, tagger.RecordedCallsArgs)
}

func Test_Retag_RetagWhenRebuild(t *testing.T) {
	t.Parallel()

	DAG := &dag.DAG{}
	DAG.AddNode(dag.NewNode(&dag.Image{
		Name:         "registry.example.org/image",
		ShortName:    "image",
		Hash:         "myhash",
		RetagDone:    false,
		NeedsRebuild: true,
	}))

	tagger := &mock.Tagger{}
	err := dib.Retag(DAG, tagger, false)

	require.NoError(t, err)
	require.Len(t, tagger.RecordedCallsArgs, 1)
	args := tagger.RecordedCallsArgs[0]
	assert.Equal(t, "registry.example.org/image:dev-myhash", args.Src)
	assert.Equal(t, "registry.example.org/image:myhash", args.Dest)
}

func Test_Retag_ReleaseWithExtraTags(t *testing.T) {
	t.Parallel()

	img := &dag.Image{
		Name:      "registry.example.org/image",
		ShortName: "image",
		Hash:      "myhash",
		ExtraTags: []string{"latest1", "latest2"},
		RetagDone: false,
	}
	DAG := &dag.DAG{}
	DAG.AddNode(dag.NewNode(img))

	tagger := &mock.Tagger{}
	err := dib.Retag(DAG, tagger, true)

	assert.NoError(t, err)

	require.Len(t, tagger.RecordedCallsArgs, 2)

	assert.Equal(t, "registry.example.org/image:myhash", tagger.RecordedCallsArgs[0].Src)
	assert.Equal(t, "registry.example.org/image:latest1", tagger.RecordedCallsArgs[0].Dest)
	assert.Equal(t, "registry.example.org/image:myhash", tagger.RecordedCallsArgs[1].Src)
	assert.Equal(t, "registry.example.org/image:latest2", tagger.RecordedCallsArgs[1].Dest)

	assert.True(t, img.RetagDone)
}

package lsp

import (
	"testing"

	"github.com/fhs/acme-lsp/internal/golang_org_x_tools/lsp/protocol"
	"github.com/google/go-cmp/cmp"
)

func TestCompatibleCodeActions(t *testing.T) {
	for _, tc := range []struct {
		name        string
		cap         protocol.ServerCapabilities
		kinds, want []protocol.CodeActionKind
	}{
		{
			"True",
			protocol.ServerCapabilities{CodeActionProvider: true},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
		},
		{
			"False",
			protocol.ServerCapabilities{CodeActionProvider: false},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
			nil,
		},
		{
			"AllFound",
			protocol.ServerCapabilities{
				CodeActionProvider: protocol.CodeActionOptions{
					CodeActionKinds: []protocol.CodeActionKind{
						protocol.QuickFix,
						protocol.SourceOrganizeImports,
					},
				},
			},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
		},
		{
			"NoneFound",
			protocol.ServerCapabilities{
				CodeActionProvider: protocol.CodeActionOptions{
					CodeActionKinds: []protocol.CodeActionKind{
						protocol.QuickFix,
					},
				},
			},
			[]protocol.CodeActionKind{protocol.SourceOrganizeImports},
			nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := CompatibleCodeActions(&tc.cap, tc.kinds)
			want := tc.want
			if !cmp.Equal(got, want) {
				t.Errorf("got %v; want %v", got, want)
			}
		})
	}
}

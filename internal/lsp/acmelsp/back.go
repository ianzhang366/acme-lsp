package acmelsp

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/fhs/acme-lsp/internal/lsp/protocol"
)

func getBookmarkFilepath() (string, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return h, err
	}

	bSuffix := ".bookmark.json"

	return fmt.Sprintf("%s/%s", h, bSuffix), nil
}

func CreateBookmarkFile() {
	bFile, err := getBookmarkFilepath()
	if err != nil {
		log.Printf("failed to generate bookmark filepath, err: %v", err)
		return
	}

	f, err := os.OpenFile(bFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Printf("failed to create bookmark file, err: %v", err)
	}
	defer f.Close()
	return
}

func (rc *RemoteCmd) AddBookmark() {
	if rc.bookmarkFile == "" {
		return
	}

	pos, _, err := rc.getPosition()
	if err != nil {
		log.Print("failed to get position for adding bookmark")
	}

	p, err := ioutil.ReadFile(rc.bookmarkFile)
	if err != nil {
		log.Print("failed read bookmarks")
	}

	prev := []*protocol.TextDocumentPositionParams{}

	if len(p) != 0 {
		if err := json.Unmarshal(p, &prev); err != nil {
			log.Printf("failed to unmarshal previous bookmarks, current: %v, prev: %v, err: %v", p, prev, err)
			return
		}
	}

	prev = append(prev, pos)

	d, err := json.Marshal(prev)
	if err != nil {
		log.Print("failed marshal bookmarks")
		return
	}

	if err = ioutil.WriteFile(rc.bookmarkFile, d, fs.FileMode(os.O_TRUNC)); err != nil {
		log.Print("failed write bookmarks to disk")
	}
}

func (rc *RemoteCmd) PopBookmark() (*protocol.TextDocumentPositionParams, error) {
	var out *protocol.TextDocumentPositionParams

	prev := []*protocol.TextDocumentPositionParams{}

	p, err := ioutil.ReadFile(rc.bookmarkFile)
	if err != nil {
		log.Printf("failed read bookmarks, err: %v", err)
	}

	if len(p) == 0 {
		return out, nil
	}

	if err := json.Unmarshal(p, &prev); err != nil {
		return out, fmt.Errorf("failed to unmarshal previous bookmarks, err: %v", err)
	}

	out = prev[len(prev)-1]

	prev = prev[:len(prev)-1]

	d := []byte{}

	if len(prev) != 0 {
		d, err = json.Marshal(prev)
		if err != nil {
			return out, fmt.Errorf("failed marshal bookmarks, err: %v", err)
		}
	}

	if err = ioutil.WriteFile(rc.bookmarkFile, d, fs.FileMode(os.O_TRUNC)); err != nil {
		log.Print("PopBookmark failed write bookmarks to disk %v", err)
	}

	return out, nil
}

func (rc *RemoteCmd) PlumbBookmark(pos *protocol.TextDocumentPositionParams) error {
	loc := protocol.Location{
		URI:   pos.TextDocument.URI,
		Range: protocol.Range{Start: pos.Position},
	}

	return PlumbLocations([]protocol.Location{loc})
}

func (rc *RemoteCmd) Back() error {
	if rc.bookmarkFile == "" {
		return fmt.Errorf("failed to pop bookmark, err: bookmark file path doesn't exist")
	}

	pos, err := rc.PopBookmark()
	if err != nil {
		return fmt.Errorf("failed to pop bookmark %w", err)
	}

	if pos == nil {
		return nil
	}

	if err := rc.PlumbBookmark(pos); err != nil {
		return fmt.Errorf("failed to plumb bookmark %w", err)
	}

	return nil
}

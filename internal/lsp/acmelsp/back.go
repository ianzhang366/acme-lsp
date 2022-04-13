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

func createBookmarkFile(bFile string) {
	f, err := os.OpenFile(bFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("failed to create bookmark file, err: %v", err)
	}
	f.Close()

	return
}

func (rc *RemoteCmd) AddBookmark() {
	pos, _, err := rc.getPosition()
	if err != nil {
		log.Print("failed to get position for adding bookmark")
	}

	p, err := ioutil.ReadFile(rc.bookmarkFile)
	if err != nil {
		log.Print("failed read bookmarks")
	}

	log.Printf("AddBookmark before %v", string(p))
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

	log.Printf("AddBookmark prev %v", prev)
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

	log.Printf("PopBookmark prev {%v} from %v", string(p), rc.bookmarkFile)

	if err := json.Unmarshal(p, &prev); err != nil {
		return out, fmt.Errorf("failed to unmarshal previous bookmarks, err: %v", err)
	}

	out = prev[len(prev)-1]

	prev = prev[:len(prev)-1]

	d := []byte{}

	log.Printf("PopBookmark prev %v", prev)
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

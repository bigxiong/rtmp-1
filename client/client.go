package client

import (
	"io"

	"github.com/WatchBeam/rtmp/chunk"
	"github.com/WatchBeam/rtmp/control"
	"github.com/WatchBeam/rtmp/handshake"
)

// Client represents a client connected to a RTMP server (see
// github.com/WatchBeam/rtmp/server for more). Clients are able to be written to
// and read from, and may have additional metadata attached to them in the
// future.
type Client struct {
	chunks      *chunk.Parser
	chunkWriter chunk.Writer

	controlStream *control.Stream

	// Conn represents the readable and writeable connection that links to
	// the client. This may be a net.Conn, or even just a bytes.Buffer.
	Conn io.ReadWriter
}

// New instantiates and returns a pointer to a new instance of type Client. The
// client is initialized with the given connection.
func New(conn io.ReadWriter) (*Client, error) {
	chunkWriter := chunk.NewWriter(conn, chunk.DefaultReadSize)
	chunks := chunk.NewParser(
		chunk.NewReader(conn, chunk.DefaultReadSize),
		chunk.NewNormalizer(),
	)

	controlChunks, err := chunks.Stream(2)
	if err != nil {
		return nil, err
	}

	return &Client{
		chunks:      chunks,
		chunkWriter: chunkWriter,

		controlStream: control.NewStream(
			controlChunks,
			chunkWriter,
			control.NewParser(),
			control.NewChunker(),
		),

		Conn: conn,
	}, nil
}

// Handshake preforms the handshake operation against the connecting client. If
// an error is encountered during any point of the handshake process, it will be
// returned immediately.
//
// See github.com/WatchBeam/RTMP/handshake for details.
func (c *Client) Handshake() error {
	return handshake.With(&handshake.Param{
		Conn: c.Conn,
	}).Handshake()
}

// Controls returns the stream of control sequences that are being received
// from the connected client.
func (c *Client) Controls() *control.Stream { return c.controlStream }

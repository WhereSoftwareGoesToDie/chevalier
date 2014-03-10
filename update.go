package chevalier

import (
	zmq "github.com/pebbe/zmq4"
)

// GetFullContents requests a contents list for origin from a Vaultaire
// readerd listening on endpoint, returning it as a DataSourceBurst.
func GetContents(endpoint, origin string) (*DataSourceBurst, error) {
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}
	err = sock.Connect(endpoint)
	if err != nil {
		return nil, err
	}
	_, err = sock.Send(origin, 0)
	if err != nil {
		return nil, err
	}
	b, err := sock.RecvBytes(0)
	if err != nil {
		return nil, err
	}
	burst, err := UnmarshalSourceBurst(b)
	if err != nil {
		return nil, err
	}
	return burst, nil
}

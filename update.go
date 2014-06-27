package chevalier

import (
	zmq "github.com/pebbe/zmq4"
)

// GetContents list for origin from a Vaultaire
// readerd listening on endpoint, returning it as a DataSourceBurst.
func GetContents(endpoint, origin string) ([]*ElasticsearchSource, error) {
	sources := make([]*ElasticsearchSource, 0)
	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return nil, err
	}
	err = sock.Connect(endpoint)
	if err != nil {
		return nil, err
	}
	request := make([][]byte, 2)
	request[0] = make([]byte, len(origin))
	for idx, chr := range origin {
		request[0][idx] = byte(chr)
	}
	request[1] = make([]byte, 1)
	request[1][0] = byte(ContentsListRequest)
	_, err = sock.SendMessage(request[0], request[1])
	if err != nil {
		return nil, err
	}
	res := new(ContentsResponse)
	for res, err = recvContentsMessage(sock); notStopResponse(res); res, err = recvContentsMessage(sock)  {
		if err != nil {
			return nil, err
		}
		res_ := unpackSourceResponse(origin, res)
		sources = append(sources, res_)
	}
	return sources, nil
}

func notStopResponse(res *ContentsResponse) bool {
	if res == nil {
		return false // first iteration or error
	}
	if res.opCode == ContentsListEntry {
		return false // data response, continue
	}
	return true
}

func recvContentsMessage(sock *zmq.Socket) (*ContentsResponse, error) {
	bs, err := sock.RecvBytes(0)
	if err != nil {
		return nil, err
	}
	return unpackContentsResponse(bs)
}

func unpackSourceResponse(origin string, res *ContentsResponse) *ElasticsearchSource {
	source := new(ElasticsearchSource)
	source.Source = make(map[string]string, 0)
	source.Address = res.entry.address
	for k, v := range res.entry.tags {
		source.Source[k] = v
	}
	source.Origin = origin
	return source
}

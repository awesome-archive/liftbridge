package server

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"

	client "github.com/tylertreat/go-jetbridge/proto"

	"github.com/tylertreat/jetbridge/server/commitlog"
	"github.com/tylertreat/jetbridge/server/proto"
)

var envelopeCookie = []byte("jetb")

type stream struct {
	*proto.Stream
	mu          sync.RWMutex
	sub         *nats.Subscription // Subscription to stream NATS subject
	replSub     *nats.Subscription // Subscription for replication requests
	log         CommitLog
	srv         *Server
	subjectHash string
	replicating bool
	replicas    map[string]struct{}
	isr         map[string]struct{}
	replicators map[string]*replicator
}

func (s *Server) newStream(protoStream *proto.Stream) (*stream, error) {
	log, err := commitlog.New(commitlog.Options{
		Path:            filepath.Join(s.config.Clustering.RaftPath, protoStream.Subject, protoStream.Name),
		MaxSegmentBytes: s.config.Log.MaxSegmentBytes,
		MaxLogBytes:     s.config.Log.RetentionBytes,
		Compact:         s.config.Log.Compact,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create commit log")
	}

	h := sha1.New()
	h.Write([]byte(protoStream.Subject))
	subjectHash := fmt.Sprintf("%x", h.Sum(nil))

	replicas := make(map[string]struct{}, len(protoStream.Replicas))
	for _, replica := range protoStream.Replicas {
		replicas[replica] = struct{}{}
	}

	isr := make(map[string]struct{}, len(protoStream.Isr))
	for _, replica := range protoStream.Isr {
		isr[replica] = struct{}{}
	}

	st := &stream{
		Stream:      protoStream,
		log:         log,
		srv:         s,
		subjectHash: subjectHash,
		replicas:    replicas,
		isr:         isr,
	}

	replicators := make(map[string]*replicator, len(protoStream.Replicas))
	for _, replica := range protoStream.Replicas {
		if replica == s.config.Clustering.NodeID {
			// Don't replicate to ourselves.
			continue
		}
		replicators[replica] = &replicator{
			stream:   st,
			requests: make(chan *proto.ReplicationRequest),
			hw:       -1,
		}
	}

	return st, nil
}

func (s *stream) String() string {
	return fmt.Sprintf("[subject=%s, name=%s]", s.Subject, s.Name)
}

func (s *stream) close() error {
	if err := s.log.Close(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			return err
		}
	}
	if s.replSub != nil {
		return s.replSub.Unsubscribe()
	}
	return nil
}

func (s *stream) handleMsg(msg *nats.Msg) {
	envelope := getEnvelope(msg.Data)
	m := &proto.Message{
		MagicByte: 2,
		Timestamp: time.Now(),
		Headers:   make(map[string][]byte),
	}
	if envelope != nil {
		m.Key = envelope.Key
		m.Value = envelope.Value
		for key, value := range envelope.Headers {
			m.Headers[key] = value
		}
	} else {
		m.Value = msg.Data
	}
	m.Headers["subject"] = []byte(msg.Subject)
	if msg.Reply != "" {
		m.Headers["reply"] = []byte(msg.Reply)
	}

	ms := &proto.MessageSet{Messages: []*proto.Message{m}}
	data, err := proto.Encode(ms)
	if err != nil {
		panic(err)
	}
	offset, err := s.log.Append(data)
	if err != nil {
		s.srv.logger.Errorf("Failed to append to log %s: %v", s, err)
		return
	}

	// Publish ack.
	if envelope != nil && envelope.AckInbox != "" {
		ack := &client.Ack{
			StreamSubject: s.Subject,
			StreamName:    s.Name,
			MsgSubject:    msg.Subject,
			Offset:        offset,
		}
		data, err := ack.Marshal()
		if err != nil {
			panic(err)
		}
		s.srv.nc.Publish(envelope.AckInbox, data)
	}
}

func (s *stream) handleReplicationRequest(msg *nats.Msg) {
	req := &proto.ReplicationRequest{}
	if err := req.Unmarshal(msg.Data); err != nil {
		s.srv.logger.Errorf("Invalid replication request for stream %s: %v", s, err)
		return
	}
	s.mu.Lock()
	if _, ok := s.replicas[req.ReplicaID]; !ok {
		s.srv.logger.Warnf("Received replication request for stream %s from non-replica %s",
			s, req.ReplicaID)
		s.mu.Unlock()
		return
	}
	replicator, ok := s.replicators[req.ReplicaID]
	if !ok {
		panic(fmt.Sprintf("No replicator for stream %s and replica %s", s, req.ReplicaID))
	}
	s.mu.Unlock()
	replicator.replicate(req)
}

func (s *stream) getReplicationInbox() string {
	return fmt.Sprintf("%s.%s.%s.replicate",
		s.srv.config.Clustering.Namespace, s.subjectHash, s.Name)
}

func (s *stream) startReplicating() {
	s.mu.Lock()
	defer s.mu.Unlock()
	replicating := s.replicating
	if replicating {
		return
	}
	s.replicating = true
	for _, replicator := range s.replicators {
		go replicator.start()
	}
}

func (s *stream) stopReplicating() {
	s.mu.Lock()
	defer s.mu.Unlock()
	replicating := s.replicating
	if !replicating {
		return
	}
	s.replicating = false
	for _, replicator := range s.replicators {
		replicator.stop()
	}
}

func getEnvelope(data []byte) *client.Message {
	if len(data) <= 4 {
		return nil
	}
	if !bytes.Equal(data[0:4], envelopeCookie) {
		return nil
	}
	msg := &client.Message{}
	if err := msg.Unmarshal(data[4:]); err != nil {
		return nil
	}
	return msg
}

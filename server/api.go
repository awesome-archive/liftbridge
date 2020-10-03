package server

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	client "github.com/liftbridge-io/liftbridge-api/go"
	"github.com/liftbridge-io/liftbridge/server/commitlog"
	proto "github.com/liftbridge-io/liftbridge/server/protocol"
)

var hasher = crc32.ChecksumIEEE

// apiServer implements the gRPC server interface clients interact with.
type apiServer struct {
	*Server
}

// CreateStream creates a new stream attached to a NATS subject. It returns an
// AlreadyExists status code if a stream with the given subject and name
// already exists.
func (a *apiServer) CreateStream(ctx context.Context, req *client.CreateStreamRequest) (
	*client.CreateStreamResponse, error) {

	resp := &client.CreateStreamResponse{}
	if req.ReplicationFactor == 0 {
		req.ReplicationFactor = 1
	}
	if req.Partitions == 0 {
		req.Partitions = 1
	}
	a.logger.Debugf("api: CreateStream [name=%s, subject=%s, partitions=%d, replicationFactor=%d]",
		req.Name, req.Subject, req.Partitions, req.ReplicationFactor)

	if req.Name == "" {
		a.logger.Errorf("api: Failed to create stream: name cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "Name cannot be empty")
	}
	// TODO: Check if valid NATS subject?
	if req.Subject == "" {
		a.logger.Errorf("api: Failed to create stream: subject cannot be empty")
		return nil, status.Error(codes.InvalidArgument, "Subject cannot be empty")
	}

	partitions := make([]*proto.Partition, req.Partitions)
	for i := int32(0); i < req.Partitions; i++ {
		partitions[i] = &proto.Partition{
			Subject:           req.Subject,
			Stream:            req.Name,
			Group:             req.Group,
			ReplicationFactor: req.ReplicationFactor,
			Id:                i,
		}
	}

	stream := &proto.Stream{
		Name:       req.Name,
		Subject:    req.Subject,
		Partitions: partitions,
		Config:     getStreamConfig(req),
	}

	if e := a.metadata.CreateStream(ctx, &proto.CreateStreamOp{Stream: stream}); e != nil {
		if e.Code() != codes.AlreadyExists {
			a.logger.Errorf("api: Failed to create stream %s: %v", req.Name, e.Err())
		}
		return nil, e.Err()
	}

	return resp, nil
}

// DeleteStream deletes a stream attached to a NATS subject.
func (a *apiServer) DeleteStream(ctx context.Context, req *client.DeleteStreamRequest) (
	*client.DeleteStreamResponse, error) {

	resp := &client.DeleteStreamResponse{}
	a.logger.Debugf("api: DeleteStream [name=%s]",
		req.Name)

	if e := a.metadata.DeleteStream(ctx, &proto.DeleteStreamOp{
		Stream: req.Name,
	}); e != nil {
		a.logger.Errorf("api: Failed to delete stream %v: %v", req.Name, e.Err())
		return nil, e.Err()
	}

	return resp, nil
}

// PauseStream pauses a stream's partitions. If no partitions are specified,
// all of the stream's partitions will be paused. Partitions are resumed when
// they are published to via the Liftbridge Publish API.
func (a *apiServer) PauseStream(ctx context.Context, req *client.PauseStreamRequest) (
	*client.PauseStreamResponse, error) {

	resp := &client.PauseStreamResponse{}
	a.logger.Debugf("api: PauseStream [name=%s, partitions=%v, resumeAll=%v]",
		req.Name, req.Partitions, req.ResumeAll)

	if len(req.Partitions) == 0 {
		stream := a.metadata.GetStream(req.Name)
		if stream == nil {
			return nil, status.Error(codes.NotFound, "stream not found")
		}
		for _, partition := range stream.GetPartitions() {
			req.Partitions = append(req.Partitions, partition.Id)
		}
	}

	if e := a.metadata.PauseStream(ctx, &proto.PauseStreamOp{
		Stream:     req.Name,
		Partitions: req.Partitions,
		ResumeAll:  req.ResumeAll,
	}); e != nil {
		a.logger.Errorf("api: Failed to pause stream %v: %v", req.Name, e.Err())
		return nil, e.Err()
	}

	return resp, nil
}

// SetStreamReadonly sets the readonly status on a stream's partitions. If no
// partitions are specified, all of the stream's partitions will have their
// readonly status set.
func (a *apiServer) SetStreamReadonly(ctx context.Context, req *client.SetStreamReadonlyRequest) (
	*client.SetStreamReadonlyResponse, error) {

	resp := &client.SetStreamReadonlyResponse{}
	a.logger.Debugf("api: SetStreamReadonly [name=%s, partitions=%v, readonly=%v]",
		req.Name, req.Partitions, req.Readonly)

	if len(req.Partitions) == 0 {
		stream := a.metadata.GetStream(req.Name)
		if stream == nil {
			return nil, status.Error(codes.NotFound, "stream not found")
		}
		for _, partition := range stream.GetPartitions() {
			req.Partitions = append(req.Partitions, partition.Id)
		}
	}

	if e := a.metadata.SetStreamReadonly(ctx, &proto.SetStreamReadonlyOp{
		Stream:     req.Name,
		Partitions: req.Partitions,
		Readonly:   req.Readonly,
	}); e != nil {
		a.logger.Errorf("api: Failed to set stream readonly flag %v: %v", req.Name, e.Err())
		return nil, e.Err()
	}

	return resp, nil
}

// Subscribe creates an ephemeral subscription for the given stream partition.
// It begins to receive messages starting at the given offset and waits for new
// messages when it reaches the end of the partition. Use the request context
// to close the subscription.
func (a *apiServer) Subscribe(req *client.SubscribeRequest, out client.API_SubscribeServer) error {
	a.logger.Debugf("api: Subscribe [stream=%s, partition=%d, start=%s, offset=%d, timestamp=%d]",
		req.Stream, req.Partition, req.StartPosition, req.StartOffset, req.StartTimestamp)

	partition := a.metadata.GetPartition(req.Stream, req.Partition)
	if partition == nil {
		a.logger.Errorf("api: Failed to subscribe to partition "+
			"[stream=%s, partition=%d]: no such partition",
			req.Stream, req.Partition)
		return status.Error(codes.NotFound, "No such partition")
	}

	leader, _ := partition.GetLeader()
	if leader != a.config.Clustering.ServerID {
		if req.ReadISRReplica {
			a.logger.Info("api: Accepting subscription to partition %s: server not stream leader", partition)
		} else {
			a.logger.Errorf("api: Failed to subscribe to partition %s: server not stream leader", partition)
			return status.Error(codes.FailedPrecondition, "Server not partition leader")
		}
	}

	cancel := make(chan struct{})
	defer close(cancel)
	ch, errCh, err := a.subscribe(out.Context(), partition, req, cancel)
	if err != nil {
		a.logger.Errorf("api: Failed to subscribe to partition %s: %v", partition, err.Err())
		return err.Err()
	}

	// Send an empty message which signals the subscription was successfully
	// created.
	if err := out.Send(&client.Message{}); err != nil {
		return err
	}

	// Update the active subscriber count.
	partition.IncreaseSubscriberCount()
	defer partition.DecreaseSubscriberCount()

	for {
		select {
		case <-out.Context().Done():
			return nil
		case m := <-ch:
			if err := out.Send(m); err != nil {
				return err
			}
		case err := <-errCh:
			return err.Err()
		}
	}
}

// FetchMetadata retrieves the latest cluster metadata, including stream broker
// information.
func (a *apiServer) FetchMetadata(ctx context.Context, req *client.FetchMetadataRequest) (
	*client.FetchMetadataResponse, error) {
	a.logger.Debugf("api: FetchMetadata %s", req.Streams)

	resp, err := a.metadata.FetchMetadata(ctx, req)
	if err != nil {
		a.logger.Errorf("api: Failed to fetch metadata: %v", err.Err())
		return nil, err.Err()
	}

	return resp, nil
}

// FetchPartitionMetadata retrieves metatadata from the partition leader. This
// is mainly useful when client would like to know the high watermark and
// newest offset for a partition.
func (a *apiServer) FetchPartitionMetadata(ctx context.Context, req *client.FetchPartitionMetadataRequest) (
	*client.FetchPartitionMetadataResponse, error) {
	a.logger.Debug("api: FetchPartitionMetadata [stream=%s, partition=%s]", req.Stream, req.Partition)

	resp, err := a.metadata.FetchPartitionMetadata(ctx, req)
	if err != nil {
		a.logger.Errorf("api: Failed to fetch partition metadata: %v", err.Err())
		return nil, err.Err()
	}
	return resp, nil
}

// Publish a new message to a stream. If the AckPolicy is not NONE and a
// deadline is provided, this will synchronously block until the ack is
// received. If the ack is not received in time, a DeadlineExceeded status code
// is returned. A FailedPrecondition status code is returned if the partition is
// readonly.
func (a *apiServer) Publish(ctx context.Context, req *client.PublishRequest) (
	*client.PublishResponse, error) {

	// TODO: Deprecate in favor of PublishAsync and log a warning.
	a.logger.Debugf("api: Publish [stream=%s, partition=%d]", req.Stream, req.Partition)

	subject, e := a.getPublishSubject(req)
	if e != nil {
		a.logger.Errorf("api: Failed to publish message: %v", e.Message)
		return nil, convertPublishAsyncError(e)
	}

	if e := a.ensureStreamNotReadonly(req.Stream, req.Partition); e != nil {
		return nil, convertPublishAsyncError(e)
	}

	if err := a.resumeStream(ctx, req.Stream, req.Partition); err != nil {
		a.logger.Errorf("api: Failed to resume stream: %v", err)
		return nil, err
	}

	if req.AckInbox == "" {
		req.AckInbox = a.getAckInbox()
	}

	var (
		msg = &client.Message{
			Key:           req.Key,
			Value:         req.Value,
			Stream:        req.Stream,
			Subject:       subject,
			Headers:       req.Headers,
			AckInbox:      req.AckInbox,
			CorrelationId: req.CorrelationId,
			AckPolicy:     req.AckPolicy,
		}
		resp = new(client.PublishResponse)
	)

	ack, err := a.publish(ctx, subject, req.AckInbox, req.AckPolicy, msg)
	if err != nil {
		a.logger.Errorf("api: Failed to publish message: %v", err)
		return nil, err
	}

	resp.Ack = ack
	return resp, nil
}

// Asynchronously publish messages to a stream. This returns a stream which
// will yield PublishResponses for messages whose AckPolicy is not NONE.
func (a *apiServer) PublishAsync(stream client.API_PublishAsyncServer) error {
	ackInbox := a.getAckInbox()
	sub, err := a.ncPublishes.Subscribe(ackInbox, func(m *nats.Msg) {
		ack, err := proto.UnmarshalAck(m.Data)
		if err != nil {
			a.logger.Errorf("api: Invalid ack received on ack inbox: %v", err)
			return
		}
		if err := stream.Send(&client.PublishResponse{CorrelationId: ack.CorrelationId, Ack: ack}); err != nil {
			a.logger.Errorf("api: Failed to send PublishAsync response: %v", err)
		}
	})
	if err != nil {
		return err
	}
	sub.SetPendingLimits(-1, -1)
	defer sub.Unsubscribe()

	if err := a.publishAsyncLoop(stream, ackInbox); err != nil {
		a.logger.Errorf("api: Failed to publish async message: %v", err)
		return err
	}
	return nil
}

// Publish a Liftbridge message to a NATS subject. If the AckPolicy is not NONE
// and a deadline is provided, this will synchronously block until the first
// ack is received. If an ack is not received in time, a DeadlineExceeded
// status code is returned.
func (a *apiServer) PublishToSubject(ctx context.Context, req *client.PublishToSubjectRequest) (
	*client.PublishToSubjectResponse, error) {
	a.logger.Debugf("api: PublishToSubject [subject=%s]", req.Subject)

	if req.AckInbox == "" {
		req.AckInbox = a.getAckInbox()
	}

	var (
		msg = &client.Message{
			Key:           req.Key,
			Value:         req.Value,
			Subject:       req.Subject,
			Headers:       req.Headers,
			AckInbox:      req.AckInbox,
			CorrelationId: req.CorrelationId,
			AckPolicy:     req.AckPolicy,
		}
		resp = new(client.PublishToSubjectResponse)
	)

	ack, err := a.publish(ctx, req.Subject, req.AckInbox, req.AckPolicy, msg)
	if err != nil {
		a.logger.Errorf("api: Failed to publish message: %v", err)
		return nil, err
	}

	resp.Ack = ack
	return resp, nil
}

// SetCursor stores a cursor position for a particular stream partition which
// is uniquely identified by an opaque string.
//
// NOTE: This is a beta endpoint and is subject to change. It is not included
// as part of Liftbridge's semantic versioning scheme.
func (a *apiServer) SetCursor(ctx context.Context, req *client.SetCursorRequest) (
	*client.SetCursorResponse, error) {
	a.logger.Debugf("api: SetCursor [stream=%s, partition=%d, cursorId=%s, offset=%d]",
		req.Stream, req.Partition, req.CursorId, req.Offset)

	if req.Stream == "" {
		return nil, status.Error(codes.InvalidArgument, "No stream provided")
	}
	if req.CursorId == "" {
		return nil, status.Error(codes.InvalidArgument, "No cursorId provided")
	}

	if status := a.cursors.SetCursor(ctx, req.Stream, req.CursorId, req.Partition, req.Offset); status != nil {
		return nil, status.Err()
	}
	return new(client.SetCursorResponse), nil
}

// FetchCursor retrieves a partition cursor position.
//
// NOTE: This is a beta endpoint and is subject to change. It is not included
// as part of Liftbridge's semantic versioning scheme.
func (a *apiServer) FetchCursor(ctx context.Context, req *client.FetchCursorRequest) (
	*client.FetchCursorResponse, error) {
	a.logger.Debugf("api: FetchCursor [stream=%s, partition=%d, cursorId=%s]",
		req.Stream, req.Partition, req.CursorId)

	if req.Stream == "" {
		return nil, status.Error(codes.InvalidArgument, "No stream provided")
	}
	if req.CursorId == "" {
		return nil, status.Error(codes.InvalidArgument, "No cursorId provided")
	}

	offset, status := a.cursors.GetCursor(ctx, req.Stream, req.CursorId, req.Partition)
	if status != nil {
		return nil, status.Err()
	}
	return &client.FetchCursorResponse{Offset: offset}, nil
}

func (a *apiServer) ensureStreamNotReadonly(name string, partitionID int32) *client.PublishAsyncError {
	stream := a.metadata.GetStream(name)
	if stream == nil {
		return &client.PublishAsyncError{
			Code:    client.PublishAsyncError_NOT_FOUND,
			Message: fmt.Sprintf("no such stream: %s", name),
		}
	}
	partition := stream.GetPartition(partitionID)
	if partition == nil {
		return &client.PublishAsyncError{
			Code:    client.PublishAsyncError_NOT_FOUND,
			Message: fmt.Sprintf("no such partition: %d", partitionID),
		}
	}
	if partition.IsReadonly() {
		return &client.PublishAsyncError{
			Code:    client.PublishAsyncError_READONLY,
			Message: fmt.Sprintf("readonly partition: %d", partitionID),
		}
	}

	return nil
}

func (a *apiServer) resumeStream(ctx context.Context, streamName string, partitionID int32) error {
	stream := a.metadata.GetStream(streamName)
	if stream == nil {
		return status.Error(codes.NotFound, fmt.Sprintf("No such stream: %s", streamName))
	}
	var toResume []int32
	if stream.GetResumeAll() {
		// If ResumeAll is enabled, resume any paused partitions in the stream.
		partitions := stream.GetPartitions()
		toResume = make([]int32, 0, len(partitions))
		for _, partition := range partitions {
			if !partition.IsPaused() {
				continue
			}
			toResume = append(toResume, partition.Id)
		}
	} else {
		// Otherwise just resume the partition being published to if it's
		// paused.
		partition := stream.GetPartition(partitionID)
		if partition == nil {
			return status.Error(codes.NotFound, fmt.Sprintf("No such partition: %d", partitionID))
		}
		if partition.IsPaused() {
			toResume = []int32{partition.Id}
		}
	}

	if len(toResume) == 0 {
		return nil
	}

	req := &proto.ResumeStreamOp{
		Stream:     stream.GetName(),
		Partitions: toResume,
	}
	if e := a.metadata.ResumeStream(ctx, req); e != nil {
		return e.Err()
	}

	// Reset the ResumeAll flag on the stream.
	stream.SetResumeAll(false)
	return nil
}

func (a *apiServer) getPublishSubject(req *client.PublishRequest) (string, *client.PublishAsyncError) {
	if req.Stream == "" {
		return "", &client.PublishAsyncError{
			Code:    client.PublishAsyncError_BAD_REQUEST,
			Message: "no stream provided",
		}
	}
	stream := a.metadata.GetStream(req.Stream)
	if stream == nil {
		return "", &client.PublishAsyncError{
			Code:    client.PublishAsyncError_NOT_FOUND,
			Message: fmt.Sprintf("no such stream: %s", req.Stream),
		}
	}
	subject := stream.GetSubject()
	if req.Partition > 0 {
		subject = fmt.Sprintf("%s.%d", subject, req.Partition)
	}
	return subject, nil
}

func (a *apiServer) publishAsyncLoop(stream client.API_PublishAsyncServer, ackInbox string) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF || status.Code(err) == codes.Canceled {
				return nil
			}
			return err
		}

		if e := a.ensureStreamNotReadonly(req.Stream, req.Partition); e != nil {
			a.logger.Errorf("api: Failed to publish async message: %v", e.Message)
			a.sendPublishAsyncError(stream, req.CorrelationId, e)
			continue
		}

		req.AckInbox = ackInbox

		a.logger.Debugf("api: PublishAsync [stream=%s, partition=%d]", req.Stream, req.Partition)

		subject, e := a.getPublishSubject(req)
		if e != nil {
			a.logger.Errorf("api: Failed to publish async message: %v", e.Message)
			a.sendPublishAsyncError(stream, req.CorrelationId, e)
			continue
		}
		if err := a.resumeStream(stream.Context(), req.Stream, req.Partition); err != nil {
			err = errors.Wrap(err, "failed to resume stream")
			a.logger.Errorf("api: Failed to publish async message: %v", err)
			a.sendPublishAsyncError(stream, req.CorrelationId, &client.PublishAsyncError{
				Code:    client.PublishAsyncError_INTERNAL,
				Message: err.Error(),
			})
			continue
		}
		msg, err := proto.MarshalPublish(&client.Message{
			Key:           req.Key,
			Value:         req.Value,
			Stream:        req.Stream,
			Subject:       subject,
			Headers:       req.Headers,
			AckInbox:      req.AckInbox,
			CorrelationId: req.CorrelationId,
			AckPolicy:     req.AckPolicy,
		})
		if err != nil {
			err = errors.Wrap(err, "failed to marshal message")
			a.logger.Errorf("api: Failed to publish async message: %v", err)
			a.sendPublishAsyncError(stream, req.CorrelationId, &client.PublishAsyncError{
				Code:    client.PublishAsyncError_INTERNAL,
				Message: err.Error(),
			})
			continue
		}
		if err := a.ncPublishes.Publish(subject, msg); err != nil {
			err = errors.Wrap(err, "failed to publish to NATS")
			a.logger.Errorf("api: Failed to publish async message: %v", err)
			a.sendPublishAsyncError(stream, req.CorrelationId, &client.PublishAsyncError{
				Code:    client.PublishAsyncError_INTERNAL,
				Message: err.Error(),
			})
		}
	}
}

func (a *apiServer) sendPublishAsyncError(stream client.API_PublishAsyncServer,
	correlationID string, err *client.PublishAsyncError) {

	resp := &client.PublishResponse{
		CorrelationId: correlationID,
		// Set an Ack with an empty correlation id so we don't break older
		// clients that are unaware of AsyncError. TODO (2.0.0): Remove when
		// clients are expected to check for AsyncError.
		Ack:        &client.Ack{CorrelationId: ""},
		AsyncError: err,
	}
	if err := stream.Send(resp); err != nil {
		a.logger.Errorf("api: Failed to send PublishAsync error response: %v", err)
	}
}

func (a *apiServer) publish(ctx context.Context, subject, ackInbox string,
	ackPolicy client.AckPolicy, msg *client.Message) (*client.Ack, error) {

	buf, err := proto.MarshalPublish(msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}

	// If AckPolicy is NONE or a timeout isn't specified, then we will fire and
	// forget.
	_, hasDeadline := ctx.Deadline()
	if ackPolicy == client.AckPolicy_NONE || !hasDeadline {
		if err := a.ncPublishes.Publish(subject, buf); err != nil {
			return nil, errors.Wrap(err, "failed to publish to NATS")
		}
		return nil, nil
	}

	// Otherwise we need to publish and wait for the ack.
	return a.publishSync(ctx, subject, ackInbox, buf)
}

func (a *apiServer) publishSync(ctx context.Context, subject,
	ackInbox string, msg []byte) (*client.Ack, error) {

	sub, err := a.ncPublishes.SubscribeSync(ackInbox)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to ack inbox")
	}
	if err := sub.AutoUnsubscribe(1); err != nil {
		return nil, errors.Wrap(err, "failed to auto unsubscribe from ack inbox")
	}

	if err := a.ncPublishes.Publish(subject, msg); err != nil {
		return nil, errors.Wrap(err, "failed to publish to NATS")
	}

	ackMsg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		if err == nats.ErrTimeout {
			err = status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, err
	}

	ack, err := proto.UnmarshalAck(ackMsg.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid ack for publish")
	}
	return ack, nil
}

// subscribe sets up a subscription on the given partition and begins sending
// messages on the returned channel. The subscription will run until the cancel
// channel is closed, the context is canceled, or an error is returned
// asynchronously on the status channel.
func (a *apiServer) subscribe(ctx context.Context, partition *partition,
	req *client.SubscribeRequest, cancel chan struct{}) (
	<-chan *client.Message, <-chan *status.Status, *status.Status) {

	if req.Resume {
		if err := a.resumeStream(ctx, req.Stream, req.Partition); err != nil {
			return nil, nil, status.New(
				codes.Internal, fmt.Sprintf("Failed to resume stream: %v", err))
		}

		// Resuming a partition creates a new one, so we have to get a pointer
		// to it.
		partition = a.metadata.GetPartition(req.Stream, req.Partition)
		if partition == nil {
			a.logger.Errorf("api: Failed to subscribe to partition "+
				"[stream=%s, partition=%d]: no such partition",
				req.Stream, req.Partition)
			return nil, nil, status.New(codes.NotFound, "No such partition")
		}
	}

	startOffset, st := getStartOffset(req, partition.log)
	if st != nil {
		return nil, nil, st
	}

	endOffset := partition.log.NewestOffset()

	var (
		ch          = make(chan *client.Message)
		errCh       = make(chan *status.Status)
		reader, err = partition.log.NewReader(startOffset, false)
	)
	if err != nil {
		return nil, nil, status.New(
			codes.Internal, fmt.Sprintf("Failed to create stream reader: %v", err))
	}

	a.startGoroutine(func() {
		headersBuf := make([]byte, 28)
		for {
			// TODO: this could be more efficient.
			m, offset, timestamp, _, err := reader.ReadMessage(ctx, headersBuf)
			if err != nil {
				var s *status.Status
				if err == commitlog.ErrCommitLogDeleted {
					// Partition was deleted while subscribed.
					s = status.New(codes.NotFound, err.Error())
				} else if err == commitlog.ErrCommitLogClosed {
					// Partition was closed while subscribed (likely paused).
					code := codes.Internal
					if partition.IsPaused() {
						code = codes.FailedPrecondition
					}
					s = status.New(code, err.Error())
				} else if err == commitlog.ErrCommitLogReadonly {
					// Partition was set to readonly while subscribed.
					s = status.New(codes.ResourceExhausted, fmt.Sprintf("End of readonly partition"))
				} else {
					s = status.Convert(err)
				}

				select {
				case errCh <- s:
				case <-cancel:
				}
				return
			}
			headers := m.Headers()
			var (
				msg = &client.Message{
					Stream:       partition.Stream,
					Partition:    partition.Id,
					Offset:       offset,
					Key:          m.Key(),
					Value:        m.Value(),
					Timestamp:    timestamp,
					Headers:      headers,
					Subject:      string(headers["subject"]),
					ReplySubject: string(headers["reply"]),
				}
			)
			select {
			case ch <- msg:
			case <-cancel:
				return
			}
			if partition.IsReadonly() && offset == endOffset {
				s := status.New(codes.ResourceExhausted, fmt.Sprintf("End of readonly partition"))

				select {
				case errCh <- s:
				case <-cancel:
				}
				return
			}
		}
	})

	return ch, errCh, nil
}

func getStartOffset(req *client.SubscribeRequest, log commitlog.CommitLog) (int64, *status.Status) {
	var startOffset int64
	switch req.StartPosition {
	case client.StartPosition_OFFSET:
		startOffset = req.StartOffset
	case client.StartPosition_TIMESTAMP:
		offset, err := log.OffsetForTimestamp(req.StartTimestamp)
		if err != nil {
			return startOffset, status.New(
				codes.Internal, fmt.Sprintf("Failed to lookup offset for timestamp: %v", err))
		}
		startOffset = offset
	case client.StartPosition_EARLIEST:
		startOffset = log.OldestOffset()
	case client.StartPosition_LATEST:
		startOffset = log.NewestOffset()
	case client.StartPosition_NEW_ONLY:
		startOffset = log.NewestOffset() + 1
	default:
		return startOffset, status.New(
			codes.InvalidArgument,
			fmt.Sprintf("Unknown StartPosition %s", req.StartPosition))
	}

	// If log is empty, next offset will be 0.
	if startOffset < 0 {
		startOffset = 0
	}

	return startOffset, nil
}

func getStreamConfig(req *client.CreateStreamRequest) *proto.StreamConfig {
	config := new(proto.StreamConfig)
	if req.RetentionMaxAge != nil {
		config.RetentionMaxAge = &proto.NullableInt64{Value: req.RetentionMaxAge.Value}
	}
	if req.CleanerInterval != nil {
		config.CleanerInterval = &proto.NullableInt64{Value: req.CleanerInterval.Value}
	}
	if req.SegmentMaxBytes != nil {
		config.SegmentMaxBytes = &proto.NullableInt64{Value: req.SegmentMaxBytes.Value}
	}
	if req.SegmentMaxAge != nil {
		config.SegmentMaxAge = &proto.NullableInt64{Value: req.SegmentMaxAge.Value}
	}
	if req.CompactMaxGoroutines != nil {
		config.CompactMaxGoroutines = &proto.NullableInt32{Value: req.CompactMaxGoroutines.Value}
	}
	if req.RetentionMaxBytes != nil {
		config.RetentionMaxBytes = &proto.NullableInt64{Value: req.RetentionMaxBytes.Value}
	}
	if req.RetentionMaxMessages != nil {
		config.RetentionMaxMessages = &proto.NullableInt64{Value: req.RetentionMaxMessages.Value}
	}
	if req.CompactEnabled != nil {
		config.CompactEnabled = &proto.NullableBool{Value: req.CompactEnabled.Value}
	}
	if req.AutoPauseTime != nil {
		config.AutoPauseTime = &proto.NullableInt64{Value: req.AutoPauseTime.Value}
	}
	if req.AutoPauseDisableIfSubscribers != nil {
		config.AutoPauseDisableIfSubscribers = &proto.NullableBool{Value: req.AutoPauseDisableIfSubscribers.Value}
	}
	if req.MinIsr != nil {
		config.MinIsr = &proto.NullableInt32{Value: req.MinIsr.Value}
	}
	return config
}

func convertPublishAsyncError(err *client.PublishAsyncError) error {
	if err == nil {
		return nil
	}
	var code codes.Code
	switch err.Code {
	case client.PublishAsyncError_NOT_FOUND:
		code = codes.NotFound
	case client.PublishAsyncError_BAD_REQUEST:
		code = codes.InvalidArgument
	case client.PublishAsyncError_READONLY:
		code = codes.FailedPrecondition
	case client.PublishAsyncError_UNKNOWN:
		fallthrough
	default:
		code = codes.Unknown
	}

	return status.Error(code, err.Message)
}

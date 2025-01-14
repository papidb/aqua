// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: notification.proto

package notification

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NotificationService_AddNotification_FullMethodName       = "/notification.NotificationService/AddNotification"
	NotificationService_GetAllNotifications_FullMethodName   = "/notification.NotificationService/GetAllNotifications"
	NotificationService_ClearNotification_FullMethodName     = "/notification.NotificationService/ClearNotification"
	NotificationService_ClearAllNotifications_FullMethodName = "/notification.NotificationService/ClearAllNotifications"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The NotificationService defines the RPC methods.
type NotificationServiceClient interface {
	// Adds a new notification for a user.
	AddNotification(ctx context.Context, in *AddNotificationRequest, opts ...grpc.CallOption) (*Empty, error)
	// Retrieves all notifications for a user.
	GetAllNotifications(ctx context.Context, in *GetAllNotificationsRequest, opts ...grpc.CallOption) (*NotificationList, error)
	// Clears a single notification by its ID.
	ClearNotification(ctx context.Context, in *ClearNotificationRequest, opts ...grpc.CallOption) (*Empty, error)
	// Clears all notifications for a user.
	ClearAllNotifications(ctx context.Context, in *ClearAllNotificationsRequest, opts ...grpc.CallOption) (*Empty, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) AddNotification(ctx context.Context, in *AddNotificationRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, NotificationService_AddNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) GetAllNotifications(ctx context.Context, in *GetAllNotificationsRequest, opts ...grpc.CallOption) (*NotificationList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NotificationList)
	err := c.cc.Invoke(ctx, NotificationService_GetAllNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) ClearNotification(ctx context.Context, in *ClearNotificationRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, NotificationService_ClearNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) ClearAllNotifications(ctx context.Context, in *ClearAllNotificationsRequest, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, NotificationService_ClearAllNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
//
// The NotificationService defines the RPC methods.
type NotificationServiceServer interface {
	// Adds a new notification for a user.
	AddNotification(context.Context, *AddNotificationRequest) (*Empty, error)
	// Retrieves all notifications for a user.
	GetAllNotifications(context.Context, *GetAllNotificationsRequest) (*NotificationList, error)
	// Clears a single notification by its ID.
	ClearNotification(context.Context, *ClearNotificationRequest) (*Empty, error)
	// Clears all notifications for a user.
	ClearAllNotifications(context.Context, *ClearAllNotificationsRequest) (*Empty, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) AddNotification(context.Context, *AddNotificationRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNotification not implemented")
}
func (UnimplementedNotificationServiceServer) GetAllNotifications(context.Context, *GetAllNotificationsRequest) (*NotificationList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) ClearNotification(context.Context, *ClearNotificationRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearNotification not implemented")
}
func (UnimplementedNotificationServiceServer) ClearAllNotifications(context.Context, *ClearAllNotificationsRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAllNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_AddNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).AddNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_AddNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).AddNotification(ctx, req.(*AddNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_GetAllNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetAllNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_GetAllNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetAllNotifications(ctx, req.(*GetAllNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_ClearNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).ClearNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_ClearNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).ClearNotification(ctx, req.(*ClearNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_ClearAllNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearAllNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).ClearAllNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_ClearAllNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).ClearAllNotifications(ctx, req.(*ClearAllNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddNotification",
			Handler:    _NotificationService_AddNotification_Handler,
		},
		{
			MethodName: "GetAllNotifications",
			Handler:    _NotificationService_GetAllNotifications_Handler,
		},
		{
			MethodName: "ClearNotification",
			Handler:    _NotificationService_ClearNotification_Handler,
		},
		{
			MethodName: "ClearAllNotifications",
			Handler:    _NotificationService_ClearAllNotifications_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}

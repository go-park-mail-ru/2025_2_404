package handler

import (
	"2025_2_404/internal/delivery/grpc/interceptor"
	modeluser "2025_2_404/internal/service/profile/domain"
	"2025_2_404/protos/profile"
	"bytes"
	"context"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileUsecaseI interface{
	Update(ctx context.Context, client modeluser.User, file io.Reader, ext string) error
	Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, []byte, error)
	Delete(ctx context.Context, clientID modeluser.ID) error
}

type ProfileServer struct {
	profile.UnimplementedProfileServer 
	profileUsecase ProfileUsecaseI
}

func NewProfileServer(profileUsecase ProfileUsecaseI) *ProfileServer{
	return &ProfileServer{
		profileUsecase: profileUsecase,
	}	
}

func (h *ProfileServer) Update(ctx context.Context, req *profile.UpdateRequest) (*profile.UpdateResponse, error){
	clientID, err := interceptor.GetUserID(ctx) 
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	client := modeluser.User{
		ID:            clientID,
		Email:         req.GetEmail(),
		UserName:      req.GetUserName(),
		UserFirstName: req.GetFirstName(),
		UserLastName:  req.GetLastName(),
		Company:       req.GetCompany(),
		Phone:         req.GetPhone(),
	}
	
	var fileReader io.Reader
	var ext string

	avatarBytes := req.GetAvatar()
	if len(avatarBytes) > 0 {
		fileReader = bytes.NewReader(avatarBytes)
		mimeType := http.DetectContentType(avatarBytes)

		switch mimeType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unsupported image format: %s", mimeType)
		}
	} else {
		fileReader = nil
		ext = ""
	}
	
	err = h.profileUsecase.Update(ctx, client, fileReader, ext)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &profile.UpdateResponse{
		UserName:   client.UserName,
		Email:      client.Email,
		FirstName:  client.UserFirstName,
		LastName: client.UserLastName,
		Company:    client.Company,
		Phone:      client.Phone,
	}, nil
}

func (h *ProfileServer) Show(ctx context.Context, req *profile.ShowRequest) (*profile.ShowResponse, error){
	clientID, err := interceptor.GetUserID(ctx) 
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	user, imgBytes, err := h.profileUsecase.Show(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to show profile: %v", err)
	}

	return &profile.ShowResponse{
		UserName:      user.UserName,
		Email:         user.Email,
		FirstName:     user.UserFirstName,
		LastName:    user.UserLastName,
		Company:       user.Company,
		Phone:         user.Phone,
		Avatar: imgBytes, 
	}, nil
}

func (h *ProfileServer) Delete(ctx context.Context, req *profile.DeleteRequest) (*profile.DeleteResponse, error) {
	clientID, err := interceptor.GetUserID(ctx) 
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	err = h.profileUsecase.Delete(ctx, clientID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete profile: %v", err)
	}

	return &profile.DeleteResponse{}, nil
}



package boards_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	mock_accomodations "github.com/PegasusMKD/travel-dream-board/internal/accomodations/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	mock_activities "github.com/PegasusMKD/travel-dream-board/internal/activities/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	mock_boards "github.com/PegasusMKD/travel-dream-board/internal/boards/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
	mock_transport "github.com/PegasusMKD/travel-dream-board/internal/transport/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBoardsService_CreateBoard(t *testing.T) {
	mockRepo := new(mock_boards.Repository)
	mockAccSvc := new(mock_accomodations.Service)
	mockActSvc := new(mock_activities.Service)
	mockTransSvc := new(mock_transport.Service)

	svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

	ctx := context.Background()
	boardToCreate := &boards.Board{Name: "Test Board"}
	expectedBoard := &boards.Board{Name: "Test Board", Uuid: "uuid-123"}

	mockRepo.On("CreateBoard", ctx, boardToCreate).Return(expectedBoard, nil)

	result, err := svc.CreateBoard(ctx, boardToCreate)
	assert.NoError(t, err)
	assert.Equal(t, expectedBoard, result)

	mockRepo.AssertExpectations(t)
}

func TestBoardsService_GetBoardById(t *testing.T) {
	ctx := context.Background()
	boardUuid := "uuid-123"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_boards.Repository)
		mockAccSvc := new(mock_accomodations.Service)
		mockActSvc := new(mock_activities.Service)
		mockTransSvc := new(mock_transport.Service)
		svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

		mockBoard := &boards.Board{Uuid: boardUuid, Name: "Test"}
		mockAccs := []*accomodations.Accomodation{{Uuid: "acc-1"}}
		mockActs := []*activities.Activity{{Uuid: "act-1"}}
		mockTrans := []*transport.Transport{{Uuid: "trans-1"}}

		mockRepo.On("GetBoardById", ctx, boardUuid).Return(mockBoard, nil)
		mockAccSvc.On("GetAccomodationsByBoardId", ctx, boardUuid).Return(mockAccs, nil)
		mockActSvc.On("GetActivitiesByBoardId", ctx, boardUuid).Return(mockActs, nil)
		mockTransSvc.On("GetTransportsByBoardId", ctx, boardUuid).Return(mockTrans, nil)

		result, err := svc.GetBoardById(ctx, boardUuid)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockBoard.Uuid, result.Uuid)
		assert.Len(t, result.Accomodations, 1)
		assert.Len(t, result.Activities, 1)
		assert.Len(t, result.Transport, 1)

		mockRepo.AssertExpectations(t)
		mockAccSvc.AssertExpectations(t)
		mockActSvc.AssertExpectations(t)
		mockTransSvc.AssertExpectations(t)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(mock_boards.Repository)
		mockAccSvc := new(mock_accomodations.Service)
		mockActSvc := new(mock_activities.Service)
		mockTransSvc := new(mock_transport.Service)
		svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

		mockRepo.On("GetBoardById", ctx, boardUuid).Return(nil, errors.New("repo error"))

		result, err := svc.GetBoardById(ctx, boardUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "repo error", err.Error())
	})

	t.Run("Acc Error", func(t *testing.T) {
		mockRepo := new(mock_boards.Repository)
		mockAccSvc := new(mock_accomodations.Service)
		mockActSvc := new(mock_activities.Service)
		mockTransSvc := new(mock_transport.Service)
		svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

		mockBoard := &boards.Board{Uuid: boardUuid, Name: "Test"}

		mockRepo.On("GetBoardById", ctx, boardUuid).Return(mockBoard, nil)
		mockAccSvc.On("GetAccomodationsByBoardId", ctx, boardUuid).Return(nil, errors.New("acc err"))

		result, err := svc.GetBoardById(ctx, boardUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Act Error", func(t *testing.T) {
		mockRepo := new(mock_boards.Repository)
		mockAccSvc := new(mock_accomodations.Service)
		mockActSvc := new(mock_activities.Service)
		mockTransSvc := new(mock_transport.Service)
		svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

		mockBoard := &boards.Board{Uuid: boardUuid, Name: "Test"}

		mockRepo.On("GetBoardById", ctx, boardUuid).Return(mockBoard, nil)
		mockAccSvc.On("GetAccomodationsByBoardId", ctx, boardUuid).Return([]*accomodations.Accomodation{}, nil)
		mockActSvc.On("GetActivitiesByBoardId", ctx, boardUuid).Return(nil, errors.New("act err"))

		result, err := svc.GetBoardById(ctx, boardUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Trans Error", func(t *testing.T) {
		mockRepo := new(mock_boards.Repository)
		mockAccSvc := new(mock_accomodations.Service)
		mockActSvc := new(mock_activities.Service)
		mockTransSvc := new(mock_transport.Service)
		svc := boards.NewService(mockRepo, mockAccSvc, mockActSvc, mockTransSvc)

		mockBoard := &boards.Board{Uuid: boardUuid, Name: "Test"}

		mockRepo.On("GetBoardById", ctx, boardUuid).Return(mockBoard, nil)
		mockAccSvc.On("GetAccomodationsByBoardId", ctx, boardUuid).Return([]*accomodations.Accomodation{}, nil)
		mockActSvc.On("GetActivitiesByBoardId", ctx, boardUuid).Return([]*activities.Activity{}, nil)
		mockTransSvc.On("GetTransportsByBoardId", ctx, boardUuid).Return(nil, errors.New("trans err"))

		result, err := svc.GetBoardById(ctx, boardUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestBoardsService_GetAllBoards(t *testing.T) {
	mockRepo := new(mock_boards.Repository)
	svc := boards.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	userUuid := "user-1"
	expectedBoards := []*boards.Board{{Uuid: "b1"}, {Uuid: "b2"}}

	mockRepo.On("GetAllBoards", ctx, userUuid).Return(expectedBoards, nil)

	result, err := svc.GetAllBoards(ctx, userUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedBoards, result)
}

func TestBoardsService_UpdateBoardById(t *testing.T) {
	mockRepo := new(mock_boards.Repository)
	svc := boards.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	boardUuid := "b1"
	boardData := &boards.Board{Name: "Updated"}

	mockRepo.On("UpdateBoardById", ctx, boardUuid, boardData).Return(nil)

	err := svc.UpdateBoardById(ctx, boardUuid, boardData)
	assert.NoError(t, err)
}

func TestBoardsService_DeleteBoardById(t *testing.T) {
	mockRepo := new(mock_boards.Repository)
	svc := boards.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	boardUuid := "b1"

	mockRepo.On("DeleteBoardById", ctx, boardUuid).Return(nil)

	err := svc.DeleteBoardById(ctx, boardUuid)
	assert.NoError(t, err)
}

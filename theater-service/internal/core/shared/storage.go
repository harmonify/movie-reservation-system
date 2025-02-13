package shared

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
)

type (
	TheaterStorage interface {
		WithTx(tx *database.Transaction) TheaterStorage
		SaveTheater(ctx context.Context, createModel *entity.SaveTheater) error
		UpdateTheater(ctx context.Context, findModel *entity.FindOneTheater, updateModel *entity.UpdateTheater) error
		SoftDeleteTheater(ctx context.Context, findModel *entity.FindOneTheater) error
		FindOneTheater(ctx context.Context, findModel *entity.FindOneTheater) (*entity.Theater, error)
		FindManyTheaters(ctx context.Context, findModel *entity.FindManyTheaters) ([]*entity.Theater, error)
	}

	RoomStorage interface {
		WithTx(tx *database.Transaction) RoomStorage
		SaveRoom(ctx context.Context, createModel *entity.SaveRoom) error
		UpdateRoom(ctx context.Context, findModel *entity.FindOneRoom, updateModel *entity.UpdateRoom) error
		SoftDeleteRoom(ctx context.Context, findModel *entity.FindOneRoom) error
		FindOneRoom(ctx context.Context, findModel *entity.FindOneRoom) (*entity.Room, error)
		FindManyRooms(ctx context.Context, findModel *entity.FindManyRooms) ([]*entity.Room, error)
	}

	SeatStorage interface {
		WithTx(tx *database.Transaction) SeatStorage
		SaveSeat(ctx context.Context, createModel *entity.SaveSeat) error
		UpdateSeat(ctx context.Context, findModel *entity.FindSeat, updateModel *entity.UpdateSeat) error
		SoftDeleteSeat(ctx context.Context, findModel *entity.FindSeat) error
		CountRoomSeats(ctx context.Context, roomIds []string) ([]*entity.CountRoomSeats, error)
		FindShowtimeAvailableSeats(ctx context.Context, findModel *entity.FindShowtimeAvailableSeats) ([]*entity.Seat, error)
	}

	ShowtimeStorage interface {
		WithTx(tx *database.Transaction) ShowtimeStorage
		SaveShowtime(ctx context.Context, createModel *entity.SaveShowtime) error
		UpdateShowtime(ctx context.Context, findModel *entity.FindOneShowtime, updateModel *entity.UpdateShowtime) error
		SoftDeleteShowtime(ctx context.Context, findModel *entity.FindOneShowtime) error
		FindOneShowtime(ctx context.Context, findModel *entity.FindOneShowtime) (*entity.Showtime, error)
		FindManyShowtimes(ctx context.Context, findModel *entity.FindManyShowtimes) ([]*entity.Showtime, error)
	}

	TicketStorage interface {
		WithTx(tx *database.Transaction) TicketStorage
		SaveTicket(ctx context.Context, createModel *entity.SaveTicket) error
		UpdateTicket(ctx context.Context, findModel *entity.FindOneTicket, updateModel *entity.UpdateTicket) error
		SoftDeleteTicket(ctx context.Context, findModel *entity.FindOneTicket) error
		FindOneTicket(ctx context.Context, findModel *entity.FindOneTicket) (*entity.Ticket, error)
		FindManyTickets(ctx context.Context, findModel *entity.FindManyTickets) ([]*entity.Ticket, error)
		CountShowtimeTickets(ctx context.Context, showtimeIds []string) ([]*entity.CountShowtimeTicket, error)
	}
)

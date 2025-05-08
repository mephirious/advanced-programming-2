package dto

type UserOrderStatisticsRequestDTO struct {
	UserID string
}

type UserOrderStatisticsResponseDTO struct {
	UserID               string
	TotalOrders          int
	TotalCompletedOrders int
	TotalCancelledOrders int
	OrdersPerHour        map[int]int
}

type UserStatisticsRequestDTO struct {
	UserID string
}

type UserStatisticsResponseDTO struct {
	UserID         string
	TotalUsers     int
	UserOrderCount int
	MostActiveHour int
}
